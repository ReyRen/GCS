package main

import (
	pb "GCS/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log/slog"
	"time"
)

func rpcHandlerDockercontainerimagepull(rpcServerAddrPort string,
	nvmlRequestInfo *NvmlRequestInfo,
	dockerequestInfo *DockerRequestInfo,
	job *Job) error {
	// 连接grpc服务器
	conn, err := grpc.Dial(rpcServerAddrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error", "ERR_MSG", err.Error())
		return err
	}
	// 延迟关闭连接
	defer conn.Close()

	// 初始化客户端
	client := pb.NewGcsInfoCatchServiceClient(conn)
	// 初始化上下文，设置请求超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 延迟关闭请求会话
	defer cancel()

	// 调用获取stream
	stream, err := client.DockerContainerImagePull(ctx, &pb.ImagePullRequestMsg{ImageName: dockerequestInfo.imageName})
	if err != nil {
		slog.Error(" client.DockerContainerImagePull stream get err", "ERR_MSG", err.Error())
		return err
	}

	// 循环获取服务端推送的消息
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// 4. err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("rpc stream server closed")
			break
		}
		if err != nil {
			slog.Error("rpc stream server receive error", "ERR_MSG", err.Error())
			continue
		}
		resp.GetImageResp()
		job.sendMsg.Content.Log = resp.GetImageResp()
		slog.Debug("receive rpc image pull log", "OTHER_MSG", job.sendMsg.Content.Log)
		job.sendMsg.Type = 3
		job.sendMsgSignalChan <- struct{}{}
	}
	return nil
}
