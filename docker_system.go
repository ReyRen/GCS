package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log/slog"
)

func docker_test() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer cli.Close()
	/*swarm, err := cli.SwarmInspect(ctx)
	//slog.Debug("swarm information", "SWARM_JoinTokens", swarm.JoinTokens.Manager)
	slog.Debug("swarm information", "SWARM_JoinTokens", swarm.ClusterInfo.)
	*/
	swarmNodes, err := cli.NodeList(ctx, types.NodeListOptions{})
	for _, v := range swarmNodes {
		slog.Debug("swarm information", "HOSTNAME", v.Description.Hostname)
		slog.Debug("swarm information", "STATUS", v.Status)
	}
	/*reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)*/
}
