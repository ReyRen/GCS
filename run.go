package main

import (
	pb "GCS/proto"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func (c *ResourceClient) resourceHandler() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			slog.Error("resourceHandler conn close err", "ERR_MSG", err.Error())
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage() // This is a block func, once ws closed, this would be get err
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket.IsUnexpectedCloseError error", "ERR_MSG", err.Error())
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		jsonHandler(message, c.rm)
		_ = nvme_sys_handler(c)
	}
}

func resourcehandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// mute : websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
		return
	}
	slog.Debug("Receive HTTP Resource request, and upgrade to websocket", "SOURCE_ADDR", conn.RemoteAddr().String())
	//TODO 初始化结构体准备接收数据
	var rMsg recvResourceMsg
	var sMsg sendResourceMsg
	client := ResourceClient{
		conn: conn,
		rm:   &rMsg,
		sm:   &sMsg,
	}
	go client.resourceHandler()
}

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
		handleJob: func(j *Job, addrWithPort string) error {
			// 连接grpc服务器
			conn, err := grpc.Dial(addrWithPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				slog.Error("grpc.Dial get error",
					"ERR_MSG", err.Error(),
					"UID", j.receiveMsg.Content.IDs.Uid,
					"TID", j.receiveMsg.Content.IDs.Tid)
				return err
			}
			// 延迟关闭连接
			defer conn.Close()

			// 初始化客户端
			client := pb.NewGcsInfoCatchServiceDockerClient(conn)
			// 初始化上下文，设置请求超时时间为1秒
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			// 延迟关闭请求会话
			defer cancel()

			// 调用获取stream
			stream, err := client.DockerContainerImagePull(ctx, &pb.ImagePullRequestMsg{ImageName: j.receiveMsg.Content.ImageName})
			if err != nil {
				slog.Error(" client.DockerContainerImagePull stream get err",
					"ERR_MSG", err.Error(),
					"UID", j.receiveMsg.Content.IDs.Uid,
					"TID", j.receiveMsg.Content.IDs.Tid)
				return err
			}

			// 循环获取服务端推送的消息
			for {
				// 通过 Recv() 不断获取服务端send()推送的消息
				resp, err := stream.Recv()
				// err==io.EOF则表示服务端关闭stream了 退出
				if err == io.EOF {
					slog.Debug("rpc stream server closed")
					break
				}
				if err != nil {
					slog.Error("rpc stream server receive error",
						"ERR_MSG", err.Error(),
						"UID", j.receiveMsg.Content.IDs.Uid,
						"TID", j.receiveMsg.Content.IDs.Tid)
					continue
				}
				j.sendMsg.Content.Log = resp.GetImageResp()
				slog.Debug("receive rpc image pull log",
					"UID", j.receiveMsg.Content.IDs.Uid,
					"TID", j.receiveMsg.Content.IDs.Tid,
					"OTHER_MSG", j.sendMsg.Content.Log)
				j.sendMsg.Type = 3
				j.sendMsgSignalChan <- struct{}{}
			}

			// TODO create container
			//client.DockerContainerStart()

			return nil
		},
		sendMsgSignalChan: make(chan struct{}),
	}
	slog.Debug("Job initialed ok")

	//一个 job 的结束就是这个方法的结束
	defer func() {
		err := job.conn.Close()
		if err != nil {
			slog.Error("job conn.Close err", "ERR_MSG", err.Error())
		}

		if _, ok := <-job.DoneChan; ok {
			close(job.DoneChan)
		}
		if _, ok := <-job.sendMsgSignalChan; ok {
			close(job.sendMsgSignalChan)
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
				//新建任务需要加入队列中
				slog.Info("create task, job commit to queue")
				h.flowControl.CommitJob(job)
				slog.Info("commit job to job queue success")
				/*
					这是一个会阻塞的函数，直到有向job.DoneChan中写入
				*/
				job.WaitDone()
				//执行完 done 就可以释放队列任务，并且此处不阻塞了

				//TODO 创建任务（docker_system） 使用 grpc
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
