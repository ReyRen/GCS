package main

import (
	pb "GCS/proto"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log/slog"
	"strings"
)

/*
这个函数是通过docker swarm获取到各个物理 node 的信息
*/
func resourceInfo(job *Job) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("client.NewClientWithOpts err", "ERR_MSG", err.Error())
		return err
	}

	defer cli.Close()

	swarmNodes, err := cli.NodeList(ctx, types.NodeListOptions{})

	var nodesListerName []string
	var nodesListerAddr []string
	var nodesListerLabel []string
	var nodesListerStatus []string

	for _, v := range swarmNodes {
		if v.Spec.Role == "manager" {
			continue
		}
		nodesListerName = append(nodesListerName, v.Description.Hostname)
		nodesListerAddr = append(nodesListerAddr, v.Status.Addr)
		nodesListerLabel = append(nodesListerLabel, GPU_TYPE)
		nodesListerStatus = append(nodesListerStatus, string(v.Status.State))
	}

	slog.Debug("swarm info:", "nodeListerName", nodesListerName,
		"nodeListerAddr", nodesListerAddr,
		"nodesListerLabel", nodesListerLabel,
		"nodesListerStatus", nodesListerStatus)

	job.sendMsg.Content.ResourceInfo.NodesListerAddr = strings.Join(nodesListerAddr, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerName = strings.Join(nodesListerName, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerStatus = strings.Join(nodesListerStatus, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerLabel = strings.Join(nodesListerLabel, ",")

	return nil
}

func logStoreHandler(job *Job) error {
	conn, err := grpc.Dial(job.receiveMsg.Content.LogAddress+GCS_INFO_CATCH_GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error",
			"ERR_MSG", err.Error())
		return err
	}
	// 延迟关闭连接
	defer conn.Close()
	client := pb.NewGcsInfoCatchServiceDockerClient(conn)
	_, err = client.DockerLogStor(context.Background(), &pb.DockerLogStorReqMsg{
		LogFilePath:   job.receiveMsg.LogPathName,
		ContainerName: job.sendMsg.Content.ContainerName,
	})
	if err != nil && err != io.EOF {
		slog.Error("DockerLogStor get error", "ERR_MSG", err.Error())

		/*job.sendMsg.Type = WS_STATUS_BACK_CREATE_FAILED
		job.sendMsgSignalChan <- struct{}{}*/
		_ = socketClientCreate(job, SOCKET_STATUS_BACK_CREATE_FAILED)
		return err
	} else {
		// EOF了 说明日志没了
		for _, v := range *job.receiveMsg.Content.SelectedNodes {
			err := dockerDeleteHandler(v.NodeAddress, job.sendMsg.Content.ContainerName)
			if err != nil {
				slog.Error("dockerDeleteHandler get error")
			}
		}
		/*job.sendMsg.Type = WS_STATUS_BACK_STOP_NORMAL
		job.sendMsgSignalChan <- struct{}{}*/
		err := socketClientCreate(job, SOCKET_STATUS_BACK_STOP_NORMAL)
		if err != nil {
			slog.Debug("socketClientCreate error in logstor")
		}
	}
	return nil
}

func dockerLogHandler(job *Job) error {
	conn, err := grpc.Dial(job.receiveMsg.Content.LogAddress+GCS_INFO_CATCH_GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error",
			"ERR_MSG", err.Error())
		return err
	}
	// 延迟关闭连接
	defer conn.Close()
	client := pb.NewGcsInfoCatchServiceDockerClient(conn)
	stream, err := client.DockerContainerLogs(context.Background(), &pb.LogsRequestMsg{ContainerName: job.receiveMsg.Content.ContainerName})
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("log EOF")
			break
		}
		if err != nil && err != io.EOF {
			slog.Error("receive rpc container delete",
				"ERR_MSG", err.Error())
			return err
		}
		job.sendMsg.Type = WS_STATUS_BACK_SEND_LOG
		job.sendMsg.Content.Log = resp.GetLogsResp()
		job.sendMsgSignalChan <- struct{}{}
		/*slog.Debug("receive rpc container log",
		"OTHER_MSG", resp.GetLogsResp())*/
	}
	return nil
}

func dockerDeleteHandler(addr string, containerName string) error {
	conn, err := grpc.Dial(addr+GCS_INFO_CATCH_GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error",
			"ERR_MSG", err.Error())
	}
	// 延迟关闭连接
	defer conn.Close()
	client := pb.NewGcsInfoCatchServiceDockerClient(conn)
	stream, err := client.DockerContainerDelete(context.Background(), &pb.DeleteRequestMsg{ContainName: containerName})
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("delete EOF")
			break
		}
		if err != nil && err != io.EOF {
			slog.Error("receive rpc container delete",
				"ERR_MSG", err.Error())
			break
			//这里如果没有对应的容器名字，会出现 delete 的错误，但是其实是正常的，所以break
		}
		slog.Debug("receive rpc container delete",
			"OTHER_MSG", resp.GetDeleteResp())
	}
	return nil
}
