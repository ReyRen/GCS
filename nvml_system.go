package main

import (
	pb "GCS/proto"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log/slog"
	"time"
)

func (c *ResourceClient) nvme_sys_handler() error {

	for _, v := range *c.rm.OccupiedList {
		err := c.grpcHandler(v.NodeAddress, v.GPUIndex)
		if err != nil {
			slog.Error("resource grpcHandler error",
				"ERR_MSG", err.Error())
			return err
		}
	}
	return nil
}

func (c *ResourceClient) grpcHandler(addr string, gpuIdx string) error {
	// 连接grpc服务器
	conn, err := grpc.Dial(addr+GCS_INFO_CATCH_GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("resource grpc.Dial get error",
			"ERR_MSG", err.Error())
		return err
	}
	// 延迟关闭连接
	defer conn.Close()

	// 初始化客户端
	rpcClient := pb.NewGcsInfoCatchServiceDockerClient(conn)
	// 初始化上下文，设置请求超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 延迟关闭请求会话
	defer cancel()

	stream, err := rpcClient.NvmlUtilizationRate(ctx, &pb.NvmlInfoReuqestMsg{IndexID: gpuIdx})
	if err != nil {
		slog.Error(" client.NvmlUtilizationRate stream get err",
			"ERR_MSG", err.Error())
		return err
	}
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("resource rpc stream read over")
			break
		}
		if err != nil {
			slog.Error("resource rpc stream server receive error",
				"ERR_MSG", err.Error())
			continue
		}
		//组装到 sendmsg 中
		slog.Debug("GetIndexID", "VALUE", resp.GetIndexID())
		c.sm.GPUIndex = AssembleToRespondString(resp.GetIndexID())
		slog.Debug("GetOccupied", "VALUE", resp.GetOccupied())
		c.sm.Occupied = AssembleToRespondString(resp.GetOccupied())
		slog.Debug("GetTemperature", "VALUE", resp.GetTemperature())
		c.sm.Temperature = AssembleToRespondString(resp.GetTemperature())
		slog.Debug("GetMemRate", "VALUE", resp.GetMemRate())
		c.sm.MemUtilize = AssembleToRespondString(resp.GetMemRate())
		slog.Debug("GetUtilizationRate", "VALUE", resp.GetUtilizationRate())
		c.sm.Utilize = AssembleToRespondString(resp.GetUtilizationRate())

		c.sm.NodeAddress = addr
	}

	//send
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		slog.Error("c.conn.NextWriter error",
			"ERR_MSG", err.Error())
		return err
	}
	sdmsg, _ := json.Marshal(c.sm)
	_, err = w.Write(sdmsg)
	if err != nil {
		slog.Error("w.Write error",
			"ERR_MSG", err.Error())
	}
	if err := w.Close(); err != nil {
		slog.Error("w.Close error",
			"ERR_MSG", err.Error())
		return err
	}
	return nil
}
