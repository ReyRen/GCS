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

	slog.Debug("swarm nodeListerName", "nodeListerName", nodesListerName)
	slog.Debug("swarm nodeListerAddr", "nodeListerAddr", nodesListerAddr)
	slog.Debug("swarm nodesListerLabel", "nodesListerLabel", nodesListerLabel)
	slog.Debug("swarm nodesListerStatus", "nodesListerStatus", nodesListerStatus)

	job.sendMsg.Content.ResourceInfo.NodesListerAddr = strings.Join(nodesListerAddr, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerName = strings.Join(nodesListerName, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerStatus = strings.Join(nodesListerStatus, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerLabel = strings.Join(nodesListerLabel, ",")

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
