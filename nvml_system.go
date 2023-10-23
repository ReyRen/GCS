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
	"strings"
)

func (c *ResourceClient) nvml_sys_handler() error {

	switch c.rm.Type {
	case RESOUECE_GET_TYPE_ALL:
		slog.Debug("get all GPU resource")
		err := c.grpcHandler(c.rm.NodeAddress, GPU_ALL_INDEX_STRING)
		if err != nil {
			slog.Error("resource grpcHandler error",
				"ERR_MSG", err.Error())
			return err
		}
	case RESOUECE_GET_TYPE_PARTIAL:
		slog.Debug("get partial GPU resource")
		for _, v := range *c.rm.OccupiedList {
			if c.rm.NodeAddress == v.NodeAddress {
				err := c.grpcHandler(c.rm.NodeAddress, v.GPUIndex)
				if err != nil {
					slog.Error("resource grpcHandler error",
						"ERR_MSG", err.Error())
					return err
				}
				break
			}
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
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 延迟关闭请求会话
	//defer cancel()

	stream, err := rpcClient.NvmlUtilizationRate(context.Background(), &pb.NvmlInfoReuqestMsg{IndexID: gpuIdx})
	if err != nil {
		slog.Error(" client.NvmlUtilizationRate stream get err",
			"ERR_MSG", err.Error())
		return err
	}
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF { //说明资源读取正常完成了，所以 break 出循环
			slog.Debug("resource rpc stream read over")
			break
		}
		if err != nil && err != io.EOF {
			slog.Error("resource rpc stream server receive error",
				"ERR_MSG", err.Error())
			return err
		}
		//组装到 sendmsg 中
		c.sm.GPUIndex = AssembleToRespondString(resp.GetIndexID())
		c.sm.Occupied = AssembleToRespondString(resp.GetOccupied())
		c.sm.Temperature = AssembleToRespondString(resp.GetTemperature())
		c.sm.MemUtilize = AssembleToRespondString(resp.GetMemRate())
		c.sm.Utilize = AssembleToRespondString(resp.GetUtilizationRate())
	}
	slog.Debug("GetIndexID", "VALUE", c.sm.GPUIndex, "RPC_NODE", c.rm.NodeName)
	slog.Debug("GetOccupied", "VALUE", c.sm.Occupied, "RPC_NODE", c.rm.NodeName)
	slog.Debug("GetTemperature", "VALUE", c.sm.Temperature, "RPC_NODE", c.rm.NodeName)
	slog.Debug("GetMemRate", "VALUE", c.sm.MemUtilize, "RPC_NODE", c.rm.NodeName)
	slog.Debug("GetUtilizationRate", "VALUE", c.sm.Utilize, "RPC_NODE", c.rm.NodeName)

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

func checkGPUOccupiedOrNot(addr string, gpuIndex string) bool {
	conn, err := grpc.Dial(addr+GCS_INFO_CATCH_GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("resource grpc.Dial get error",
			"ERR_MSG", err.Error())
		return false
	}
	// 延迟关闭连接
	defer conn.Close()

	// 初始化客户端
	rpcClient := pb.NewGcsInfoCatchServiceDockerClient(conn)

	stream, err := rpcClient.NvmlUtilizationRate(context.Background(), &pb.NvmlInfoReuqestMsg{IndexID: gpuIndex})
	if err != nil {
		slog.Error(" client.NvmlUtilizationRate stream get err",
			"ERR_MSG", err.Error())
		return false
	}

	// 通过 Recv() 不断获取服务端send()推送的消息
	resp, err := stream.Recv()
	// err==io.EOF则表示服务端关闭stream了 退出
	if err == io.EOF {
		slog.Debug("resource rpc stream read over")
		return true //因为只有一次读取，如果第一次读取就是 EOF，说明有问题，所以不能被使用
	}
	if err != nil && err != io.EOF {
		slog.Error("resource rpc stream server receive error",
			"ERR_MSG", err.Error())
		return true //读取 stream 出错了，说明有问题，所以不能被使用
	}
	return strings.Contains(AssembleToRespondString(resp.GetOccupied()), "1")
}
