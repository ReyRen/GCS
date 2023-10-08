package main

import "github.com/gorilla/websocket"

const (
	GCS_ADDR_WITH_PORT = "172.18.127.64:8066" // gcs self address and port
)

// http to websocket upgrade variables
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}