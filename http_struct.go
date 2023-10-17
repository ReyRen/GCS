package main

import (
	"container/list"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
)

/*************QUEUE JOB HANDLE STRUCT SET*************/
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
	receiveMsg *RecvMsg
	sendMsg    *SendMsg
	conn       *websocket.Conn
	DoneChan   chan struct{}
	//handleJob只执行create容器的操作，一旦容器创建成功，即可退出
	handleJob         func(j *Job, addrWithPort string) error
	sendMsgSignalChan chan struct{}
}

/*************QUEUE JOB HANDLE STRUCT SET*************/

/*************RECEIVE MESSAGE STRUCT SET*************/
type Ids struct {
	Uid int `json:"uid"`
	Tid int `json:"tid"`
}
type RecvMsgContent struct {
	IDs                *Ids           `json:"ids"`
	OriginalModelUrl   string         `json:"originalModelUrl"`
	ContinuousModelUrl string         `json:"continuousModelUrl"`
	ModelName          string         `json:"modelName"`
	ResourceType       string         `json:"resourceType"`
	SelectedNodes      *[]SelectNodes `json:"selectedNodes"`
	ModelType          int            `json:"modelType"`
	Command            string         `json:"command"`
	FrameworkType      int            `json:"frameworkType"`
	ToolBoxName        string         `json:"toolBoxName"`
	Params             string         `json:"params"`
	SelectedDataset    string         `json:"selectedDataset"`
	ImageName          string         `json:"imageName"`
	DistributingMethod int            `json:"distributingMethod"`
	CommandBox         string         `json:"cmd"`
}
type RecvMsg struct {
	Type        int             `json:"type"`
	Admin       bool            `json:"admin"`
	Content     *RecvMsgContent `json:"content"`
	FtpFileName string
}

func newReceiveMsg() *RecvMsg {
	receiveMsgContent := newReceiveMsgContent()
	slog.Debug("newReceiveMsgContent done")
	return &RecvMsg{
		Type:        -1,
		Admin:       false,
		Content:     receiveMsgContent,
		FtpFileName: "",
	}
}

func newSelectNodes() *[]SelectNodes {
	slog.Debug("newSelectNodes done")
	var selectNodes []SelectNodes
	return &selectNodes
}

func newReceiveMsgContent() *RecvMsgContent {
	var ids *Ids
	selectNodes := newSelectNodes()
	slog.Debug("newReceiveMsgContent done")
	return &RecvMsgContent{
		IDs:           ids,
		SelectedNodes: selectNodes,
	}
}

type SelectNodes struct {
	NodeNames string `json:"nodeName"`
	NodeIp    string `json:"nodeIp"`
	GPUIndex  string `json:"gpuIndex"`
}

/*************RECEIVE MESSAGE STRUCT SET*************/

/*************SEND MESSAGE STRUCT SET*************/
type ResourceInfo struct {
	NodesListerName   string `json:"nodesListerName"`
	NodesListerAddr   string `json:"nodesListerAddr"`
	NodesListerLabel  string `json:"nodesListerLabel"`
	NodesListerStatus string `json:"nodesListerStatus"`
}
type SendMsgContent struct {
	Log           string `json:"log"`
	ContainerName string
	ResourceInfo  *ResourceInfo `json:"resourceInfo"`
}
type SendMsg struct {
	Type    int             `json:"type"`
	Content *SendMsgContent `json:"content"`
}

func newSendMsgContent() *SendMsgContent {
	return &SendMsgContent{
		Log: "",
		ResourceInfo: &ResourceInfo{
			NodesListerName:   "",
			NodesListerAddr:   "",
			NodesListerLabel:  "",
			NodesListerStatus: "",
		},
	}
}

func newSendMsg() *SendMsg {
	content := newSendMsgContent()
	return &SendMsg{
		Type:    -1,
		Content: content,
	}
}

/*************SEND MESSAGE STRUCT SET*************/

/*************NODE RESOURCE STRUCT SET*************/
type recvResourceMsg struct {
	Type         int              `json:"type"`
	NodeName     string           `json:"nodeName"`
	NodeAddress  string           `json:"nodeAddress"`
	OccupiedList *[]OccupiedLists `json:"occupiedList"`
}
type sendResourceMsg struct {
	Utilize     string `json:"utilize"`
	MemUtilize  string `json:"memUtilize"`
	Temperature string `json:"temp"`
	Occupied    string `json:"occupied"`
	GPUIndex    string `json:"gpuIndex"`
}
type OccupiedLists struct {
	NodeAddress string `json:"nodeAddress"`
	GPUIndex    string `json:"gpuIndex"`
}
type ResourceClient struct {
	conn *websocket.Conn
	rm   *recvResourceMsg
	sm   *sendResourceMsg
}

/*************NODE RESOURCE STRUCT SET*************/
