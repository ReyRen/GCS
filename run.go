package main

import (
	"bytes"
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
func (h *MyHandler) recvMsgHandler(conn *websocket.Conn) {

	slog.Debug("Receive HTTP request, and upgrade to websocket", "SOURCE_ADDR", conn.RemoteAddr().String())

	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Error("recvMsgHandler conn.Close err", "ERR_MSG", err.Error())
		}
		/* TODO pop队列 */
	}()
	for {
		receiveMsg := newReceiveMsg(conn)
		slog.Debug("newReceiveMsg done")

		_, message, err := conn.ReadMessage() // This is a block func, once ws closed, this would be get err
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket.IsUnexpectedCloseError error", "ERR_MSG", err.Error())
			}
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		jsonHandler(message, receiveMsg)

		slog.Debug("receive message display", "RAW_MSG", string(message))
		slog.Debug("receive msg to struct display", "UID", receiveMsg.Content.IDs.Uid)
		/*接收连接进来的消息*/

		sendMsg := newSendMsg()
		//一开始的 job 默认是发送资源信息，不提交到队列中执行
		job := &Job{
			receiveMsg: receiveMsg,
			sendMsg:    sendMsg,
			conn:       conn,
			DoneChan:   make(chan struct{}, 1),
			handleJob: func(j *Job) error {
				//TODO 提交任务执行的程序

				return nil
			},
		}

		go func() {
			switch job.receiveMsg.Type {
			case 1:
			//TODO 获取资源信息
			case 2:
			//TODO 开始任务
			default:
			}
		}()
	}
}
