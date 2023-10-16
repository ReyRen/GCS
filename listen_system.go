package main

import (
	"log/slog"
	"net/http"
)

func listenHandler() {
	flowControl := NewFlowControl()
	myHTTPHandler := MyHandler{
		flowControl: flowControl,
	}
	//ServeHTTP是MyHandler自己实现的 http 接口方法
	http.Handle("/", &myHTTPHandler)
	err := http.ListenAndServe(GCS_ADDR_WITH_PORT, nil)
	if err != nil {
		slog.Error("ListenAndServe error", "ERR_MSG", err)
	}
	/*自动调用run.go中，实现的接口函数func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)*/
}

func listenResourceHandler() {
	http.HandleFunc("/resource", func(writer http.ResponseWriter, request *http.Request) {
		resourcehandler(writer, request)
	})
	err := http.ListenAndServe(GCS_RESOURCE_WITH_PORT, nil)
	if err != nil {
		slog.Error("ListenAndServe error", "ERR_MSG", err)
	}
}
