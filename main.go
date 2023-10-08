package main

import (
	"github.com/sevlyar/go-daemon"
	"log/slog"
)

func main() {

	//Setup log system
	logSysInit()
	//log system ready

	//Setup daemon system
	cntxt := &daemon.Context{
		PidFileName: "gcs.pid",
		PidFilePerm: 0644,
		LogFileName: "./log/gcs.log",
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[gcs]"},
	}
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			slog.Error("cntxt.Search error", "ERR_MSG", err.Error())
		}
		daemon.SendCommands(d)
		return
	}
	d, err := cntxt.Reborn()
	if err != nil {
		slog.Error("cntxt.Reborn error", "ERR_MSG", err.Error())
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	slog.Info("- - - - - - -[GCS] started - - - - - - -")
	defer func() {
		slog.Info("- - - - - - -[GCS] exited- - - - - - -")
	}()
	//Daemon system ready

	slog.Debug("listenHandler start")
	listenHandler()
	slog.Debug("listenHandler done")

	//nvme_sys_init()
}

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
