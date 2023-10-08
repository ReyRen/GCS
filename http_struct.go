package main

import (
	"container/list"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
)

/*********QUEUE JOB HANDLE STRUCT SET*********/
type MyHandler struct {
	flowControl *FlowControl
}
type FlowControl struct {
	jobQueue *JobQueue
	wm       *WorkerManager
}
type JobQueue struct {
	mu         sync.Mutex
	noticeChan chan struct{}
	queue      *list.List
	size       int
	capacity   int
}

func NewJobQueue(cap int) *JobQueue {
	return &JobQueue{
		capacity:   cap,
		queue:      list.New(),
		noticeChan: make(chan struct{}, 1),
	}
}

type WorkerManager struct {
	jobQueue *JobQueue
}

func NewWorkerManager(jobQueue *JobQueue) *WorkerManager {
	return &WorkerManager{
		jobQueue: jobQueue,
	}
}

type Job struct {
	receiveMsg *recvMsg
	sendMsg    *sendMsg
	conn       *websocket.Conn
	DoneChan   chan struct{}
	handleJob  func(j *Job) error
}

/*********QUEUE JOB HANDLE STRUCT SET*********/

/*********xxxxxx xxxxx struct*********/
type Ids struct {
	Uid int `json:"uid"`
	Tid int `json:"tid"`
}
type recvMsgContent struct {
	IDs                *Ids   `json:"ids"`
	OriginalModelUrl   string `json:"originalModelUrl"`
	ContinuousModelUrl string `json:"continuousModelUrl"`
	ModelName          string `json:"modelName"`
	ResourceType       string `json:"resourceType"`
	//SelectedNodes      *[]selectNodes `json:"selectedNodes"`
	ModelType          int    `json:"modelType"`
	Command            string `json:"command"`
	FrameworkType      int    `json:"frameworkType"`
	ToolBoxName        string `json:"toolBoxName"`
	Params             string `json:"params"`
	SelectedDataset    string `json:"selectedDataset"`
	ImageName          string `json:"imageName"`
	DistributingMethod int    `json:"distributingMethod"`
	CommandBox         string `json:"cmd"`
}
type recvMsg struct {
	Type        int             `json:"type"`
	Admin       bool            `json:"admin"`
	Content     *recvMsgContent `json:"content"`
	FtpFileName string
}

func newReceiveMsg(conn *websocket.Conn) *recvMsg {
	receiveMsgContent := newReceiveMsgContent(conn)
	slog.Debug("newReceiveMsgContent done")
	return &recvMsg{
		Type:        -1,
		Admin:       false,
		Content:     receiveMsgContent,
		FtpFileName: "",
	}
}

func newReceiveMsgContent(conn *websocket.Conn) *recvMsgContent {
	var ids *Ids
	slog.Debug("newReceiveMsgContent done")
	return &recvMsgContent{
		IDs: ids,
	}
}

type resourceInfo struct {
	NodesListerName   string `json:"nodesListerName"`
	NodesListerLabel  string `json:"nodesListerLabel"`
	NodesListerStatus string `json:"nodesListerStatus"`
}
type sendMsgContent struct {
	Log          string        `json:"log"`
	ResourceInfo *resourceInfo `json:"resourceInfo"`
}
type sendMsg struct {
	Type    int             `json:"type"`
	Content *sendMsgContent `json:"content"`
}

func newSendMsgContent() *sendMsgContent {
	return &sendMsgContent{
		Log: "",
		ResourceInfo: &resourceInfo{
			NodesListerName:   "",
			NodesListerLabel:  "",
			NodesListerStatus: "",
		},
	}
}

func newSendMsg() *sendMsg {
	content := newSendMsgContent()
	return &sendMsg{
		Type:    -1,
		Content: content,
	}
}
