package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// mute : websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
		return
	}

	h.recvMsgHandler(conn)

	// 一旦点击提交任务，那么就将 job 进行 commit，并且进入队列中

	/*job := &Job{
		DoneChan: make(chan struct{}, 1),
		handleJob: func(job *Job) error {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Hello World"))
			return nil
		},
	}

	h.flowControl.CommitJob(job)
	fmt.Println("commit job to job queue success")
	job.WaitDone()*/
}
func (h *MyHandler) recvMsgHandler(conn *websocket.Conn) {

	slog.Debug("Receive HTTP request, and upgrade to websocket", "SOURCE_ADDR", conn.RemoteAddr().String())

	//初始化
	/*
		1. 接收数据结构体
		2. 发送数据结构体
		3. Job 结构体
			一开始的 job 默认是发送资源信息，不提交到队列中执行
			一个新的 conn 连接进来就产生一个新的 job
	*/
	receiveMsg := newReceiveMsg()
	slog.Debug("newReceiveMsg initialed ok")
	sendMsg := newSendMsg()
	slog.Debug("newSendMsg initialed ok")
	job := &Job{
		receiveMsg: receiveMsg,
		sendMsg:    sendMsg,
		conn:       conn,
		DoneChan:   make(chan struct{}),
		handleJob: func(j *Job) error {
			//TODO 提交任务执行的程序

			return nil
		},
		sendMsgSignalChan: make(chan struct{}),
	}
	slog.Debug("Job initialed ok")

	defer func() {
		err := job.conn.Close()
		if err != nil {
			slog.Error("recvMsgHandler conn.Close err", "ERR_MSG", err.Error())
		}
		/* TODO pop队列 */
	}()

	go h.sendMsgHandler(job)

	for {
		_, message, err := conn.ReadMessage() // This is a block func, once ws closed, this would be get err
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket.IsUnexpectedCloseError error", "ERR_MSG", err.Error())
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		jsonHandler(message, receiveMsg)

		slog.Debug("receive message display", "RAW_MSG", string(message))
		slog.Debug("receive msg to struct display", "UID", receiveMsg.Content.IDs.Uid)
		/*接收连接进来的消息*/

		//将新收到的 receiveMsg 更新给当前 job
		job.receiveMsg = receiveMsg

		go func() {
			switch job.receiveMsg.Type {
			case MESSAGE_TYPE_RESOURCE_INFO:
				//TODO 获取物理节点状态信息（docker swarm）
				//TODO 如果物理节点状态正常，获取 GPU 信息（nvml_system）使用 grpc
			case MESSAGE_TYPE_CREATE:
				//TODO 创建任务（docker_system） 使用 grpc
				//TODO 新建任务需要加入队列中
			case MESSAGE_TYPE_LOG:
				//TODO 获取容器日志
			case MESSAGE_TYPE_STOP:
				//TODO 任务停止（docker_system） 使用 grpc
			default:
				slog.Warn("receive message type not implemented", "OTHER_MSG", job.receiveMsg.Type)
			}
		}()

	}
}

func (h *MyHandler) sendMsgHandler(job *Job) {
	for {
		select {
		case _, ok := <-job.sendMsgSignalChan:
			if !ok {
				// chan closed
				slog.Debug("sendMsgSignalChan is closed",
					"UID", job.receiveMsg.Content.IDs.Uid,
					"TID", job.receiveMsg.Content.IDs.Tid)
				_ = job.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := job.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error("sendMsgSignalChan nextWriter error",
					"UID", job.receiveMsg.Content.IDs.Uid,
					"TID", job.receiveMsg.Content.IDs.Tid,
					"ERR_MSG", err.Error())
				return
			}
			sdmsg, _ := json.Marshal(job.sendMsg)
			_, err = w.Write(sdmsg)
			if err != nil {
				slog.Error("resourceInfoLogChan write error",
					"UID", job.receiveMsg.Content.IDs.Uid,
					"TID", job.receiveMsg.Content.IDs.Tid,
					"ERR_MSG", err.Error())
				return
			}
			if err := w.Close(); err != nil {
				slog.Error("resourceInfoLogChan close error",
					"UID", job.receiveMsg.Content.IDs.Uid,
					"TID", job.receiveMsg.Content.IDs.Tid,
					"ERR_MSG", err.Error())
				return
			}
		}
	}
}
