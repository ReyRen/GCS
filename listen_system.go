package main

import (
	"log/slog"
	"net/http"
)

func listen_handler() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		run_handler(writer, request)
	})
	err := http.ListenAndServe(GCS_ADDR_WITH_PORT, nil)
	if err != nil {
		slog.Error("ListenAndServe error", "ERR_MSG", err)
	}
}
