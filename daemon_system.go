package main

import (
	"github.com/sevlyar/go-daemon"
	"log/slog"
)

func daemon_sys_init() {
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
			slog.Error("cntxt.Search error", "others", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		slog.Error("cntxt.Reborn error", "others", err.Error())
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	slog.Info("- - - - - - - - - - - - - - -")
	slog.Info("[GCS] started")
}
