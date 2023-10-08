package main

/*
	Log system created by yuanren 2023/10/7
*/
import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
)

func log_sys_init() {
	opts := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}
	r := &lumberjack.Logger{
		Filename:   "./log/gcs.log",
		LocalTime:  true,
		MaxSize:    512,  // max log size is 512M
		MaxAge:     3650, // max age is 10 years
		MaxBackups: 99,   // max backup count is 99
		Compress:   true,
	}
	logger := slog.New(slog.NewJSONHandler(r, &opts))
	slog.SetDefault(logger)
}
