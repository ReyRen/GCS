package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"math/rand"
	"time"
)

const (
	GCS_ADDR_WITH_PORT = "172.18.127.64:8066" // gcs self address and port

	MESSAGE_TYPE_RESOURCE_INFO = 1
	MESSAGE_TYPE_CREATE        = 2
	MESSAGE_TYPE_LOG           = 3
	MESSAGE_TYPE_STOP          = 4
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

func GetRandomString(l int) string {
	str := "0123456789abcefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetContainerName(uid string, tid string) string {
	randomString := GetRandomString(8)
	return randomString + "-" + uid + "-" + tid
}
