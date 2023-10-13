package main

import (
	pb "GCS/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log/slog"
	"strconv"
	"time"
)

func nvme_sys_handler(c *ResourceClient) error {
	// 连接grpc服务器
	conn, err := grpc.Dial(c.rm.NodeAddress+GCS_INFO_CATCH_GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("resource grpc.Dial get error",
			"ERR_MSG", err.Error(),
			"RPC_ADDR", c.rm.NodeAddress)
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

	stream, err := rpcClient.NvmlUtilizationRate(ctx, &pb.NvmlInfoReuqestMsg{Type: strconv.Itoa(1)})
	if err != nil {
		slog.Error(" client.NvmlUtilizationRate stream get err",
			"ERR_MSG", err.Error(),
			"RPC_ADDR", c.rm.NodeAddress)
		return err
	}
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("resource rpc stream server closed")
			break
		}
		if err != nil {
			slog.Error("resource rpc stream server receive error",
				"ERR_MSG", err.Error(),
				"RPC_ADDR", c.rm.NodeAddress)
			continue
		}
		slog.Info("GetIndexID", "VALUE", resp.GetIndexID())
		slog.Info("GetOccupied", "VALUE", resp.GetOccupied())
		slog.Info("GetTemperature", "VALUE", resp.GetTemperature())
		slog.Info("GetMemRate", "VALUE", resp.GetMemRate())
		slog.Info("GetUtilizationRate", "VALUE", resp.GetUtilizationRate())
	}

	return nil
}
