package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	slog.Info("recieve http sssrequest")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// mute : websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
		return
	}

	slog.Info("recieve http request")

	/*这里处理没有开始任务的逻辑*/
	//read the first connection
	receiveMsg := newReceiveMsg(conn)
	_, message, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			slog.Error("websocket.IsUnexpectedCloseError error", "ERR_MSG", err.Error())
		}
	}
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	jsonHandler(message, receiveMsg)

	slog.Debug("receive message raw", "MSG", message)
	slog.Debug("receive message receiveMsg", "UID", receiveMsg.Content.IDs.Uid)
	/*这里处理没有开始任务的逻辑*/

	sendMsg := newSendMsg()
	//一开始的 job 默认是发送资源信息，不提交到队列中
	_ = &Job{
		receiveMsg: receiveMsg,
		sendMsg:    sendMsg,
		conn:       conn,
		DoneChan:   make(chan struct{}, 1),
		handleJob:  nil,
	}

	//job.resourceInfoResponse()

	// 一单点击提交任务，那么就将 job 进行 commit，并且进入队列中

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

/*
func (job *Job) resourceInfoResponse() {
	if job.receiveMsg.Type {

	}
}*/
