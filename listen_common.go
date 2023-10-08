package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
)

const (
	GCS_ADDR_WITH_PORT = "172.18.127.64:8066" // gcs self address and port

)

// http to websocket upgrade variables
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func jsonHandler(data []byte, v interface{}) {
	errJson := json.Unmarshal(data, v)
	if errJson != nil {
		slog.Error("jsonHandler error", "ERR_MSG", errJson.Error())
	}
}
