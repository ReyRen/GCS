package main

func main() {

	log_sys_init()
	daemon_sys_init()

	/*for i := 0; i < 100000; i++ {
		slog.Info("greeting", "say", s)
	}*/
}

var (
	SUCCESS = 1
	FAILED  = 0
)

/*func main() {
ctx := context.Background()

cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
	panic(err)
}

defer cli.Close()

reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
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

stdcopy.StdCopy(os.Stdout, os.Stderr, out)
*/
