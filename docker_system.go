package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log/slog"
	"strings"
)

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
	job.sendMsg.Type = 1 //这个 type 没有被前端处理
	job.sendMsg.Content.ResourceInfo.NodesListerAddr = strings.Join(nodesListerAddr, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerName = strings.Join(nodesListerName, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerStatus = strings.Join(nodesListerStatus, ",")
	job.sendMsg.Content.ResourceInfo.NodesListerLabel = strings.Join(nodesListerLabel, ",")

	return nil
}
