package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	/******************GCS log 文件名以及路径******************/
	GCS_LOG_PATH = "./log/gcs.log"

	/******************GCS 和 GCS resource 监听地址端口******************/
	GCS_ADDR_WITH_PORT       = "172.18.127.66:8066" // gcs self address and port
	GCS_RESOURCE_WITH_PORT   = "172.18.127.66:8067" // gcs resource self address and port
	GCS_INFO_CATCH_GRPC_PORT = ":50001"

	DOCKER_IMAGES_PREFIX = "172.18.127.68:80/"

	/******************GPU 类型，用于传输给 ws 前段显示******************/
	GPU_TYPE = "SXM-A800-80G"

	/******************训练日志存储的前置路径******************/
	LOG_STOR_PRE_PATH = "/storage-ftp-data/user/"
	//LOG_STOR_PRE_PATH = "/home/ftper/ftp/user/"

	/******************从 ws 前端收到的状态码******************/
	MESSAGE_TYPE_NODE_INFO      = 1 //表示获取物理节点状态信息
	MESSAGE_TYPE_START_CREATION = 2 //表示容器资源创建
	MESSAGE_TYPE_LOG            = 3 //表示获取容器日志
	MESSAGE_TYPE_STOP           = 4 //表示强制停止

	/******************从 ws resource前端收到的状态码******************/
	RESOUECE_GET_TYPE_ALL     = 1 //获取所有资源
	RESOUECE_GET_TYPE_PARTIAL = 2 //获取gpuindex 资源
	GPU_ALL_INDEX_STRING      = "0,1,2,3,4,5,6,7"

	/******************从 ws resource前端收到的状态码******************/
	WS_STATUS_BACK_SEND_LOG              = 3  //表示发送日志
	WS_STATUS_BACK_CREATE_RECEIVE        = 10 //表示训练指令已发送
	WS_STATUS_BACK_RESOURCE_INSUFFICIENT = 11 //表示训练资源不满足
	WS_STATUS_BACK_CREATING              = 12 //表示容器创建中
	WS_STATUS_BACK_TRAINNING             = 13 //表示任务创建成功，并且开始训练了
	WS_STATUS_BACK_CREATE_FAILED         = 14 //表示容器创建失败
	WS_STATUS_BACK_STOP_INNORMAL         = 16 //表示点击停止强制结束的任务，是非正常的任务
	WS_STATUS_BACK_STOP_NORMAL           = 15 //表示任务正常结束，日志读取到 EOF

	/******************与后端 socket 建立连接的 ip 地址和端口******************/
	socketServer                           = "172.18.127.68:8020"
	SOCKET_STATUS_BACK_CREATE_START        = 4   //表示创建容器开始（资源已判定为满足）
	SOCKET_STATUS_BACK_TRAINNING           = 6   //表示容器创建全部成功，并且训练中
	SOCKET_STATUS_BACK_STOP_NORMAL         = 7   //表示任务正常结束，日志读取到 EOF
	SOCKET_STATUS_BACK_STOP_INNOMAL        = 8   //表示任务手动结束
	SOCKET_STATUS_BACK_CREATE_FAILED       = 401 //表示容器创建失败
	SOCKET_STATUS_BACK_RESOURCE_INSUFFICIE = 402
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

var (
	//UPDATEMAP map[string][]string

	UPDATEMAP sync.Map
)

/*
[

	"statusID"
	"contaierName", //job.sendMsg.Content.ContainerName
	"LogPathName",  //job.receiveMsg.LogPathName
	"LogAddress",   //job.receiveMsg.Content.LogAddress
	"nodeName1+nodeIP1+gpuIndex,nodeName2+nodeIP2+gpuIndex,.." //job.receiveMsg.Content.SelectedNodes

]
*/

type MapValue struct {
	StatusID         string              `json:"status_id"`
	ContainerName    string              `json:"container_name"`
	LogPathName      string              `json:"log_path_name"`
	LogAddress       string              `json:"log_address"`
	ContainerInfoMap *[]containerInfoMap `json:"containerInfo"`
}
type containerInfoMap struct {
	NodeAddress string `json:"nodeAddress"`
	NodeName    string `json:"nodeName"`
	GPUIndex    string `json:"gpuIndex"`
}

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

func AssembleToRespondString(raw interface{}) string {

	var tmpString []string
	switch raw.(type) {
	case []int32:
		for _, value := range raw.([]int32) {
			if value == 99999 {
				// error get
				tmpString = append(tmpString, " ")
				continue
			}
			tmpString = append(tmpString, strconv.Itoa(int(value)))
		}
	case []uint32:
		for _, value := range raw.([]uint32) {
			if value == 99999 {
				// error get
				tmpString = append(tmpString, " ")
				continue
			}
			tmpString = append(tmpString, strconv.Itoa(int(value)))
		}
	}
	return strings.Join(tmpString, ",")
}

func socketClientCreate(job *Job, statusCode int) error {
	slog.Debug("socket client creating")
	// create socket client
	conn, err := net.Dial("tcp", socketServer)
	if err != nil {
		slog.Error("socketClientCreate err", "UID", job.receiveMsg.Content.IDs.Uid, "TID", job.receiveMsg.Content.IDs.Tid)
		return err
	}
	defer conn.Close()
	var containerInfo containerInfoList
	var tmpSlice []containerInfoList
	var socketSendMsg socketSendMsg
	socketSendMsg.Uid = job.receiveMsg.Content.IDs.Uid
	socketSendMsg.Tid = job.receiveMsg.Content.IDs.Tid
	socketSendMsg.ContainerName = job.sendMsg.Content.ContainerName
	socketSendMsg.StatusId = statusCode

	for _, v := range *job.receiveMsg.Content.SelectedNodes {
		containerInfo.GPUIndex = v.GPUIndex
		containerInfo.NodeAddress = v.NodeAddress
		containerInfo.NodeName = v.NodeName
		tmpSlice = append(tmpSlice, containerInfo)
	}
	socketSendMsg.ContainerInfoList = &tmpSlice
	socketmsg, _ := json.Marshal(socketSendMsg)
	_, err = conn.Write(socketmsg)
	if err != nil {
		slog.Error("socketClientCreate write err", "UID", job.receiveMsg.Content.IDs.Uid, "TID", job.receiveMsg.Content.IDs.Tid)
		return err
	}
	return nil
}
