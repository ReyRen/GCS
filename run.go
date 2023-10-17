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
	"strconv"
	"time"
)

/*******************单独的获取资源状况的websocket handler*******************/
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
		slog.Debug("receive resource message display", "RAW_MSG", string(message))
		jsonHandler(message, c.rm)
		_ = c.nvme_sys_handler()
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
	//初始化结构体准备接收数据
	var rMsg recvResourceMsg
	var sMsg sendResourceMsg
	client := ResourceClient{
		conn: conn,
		rm:   &rMsg,
		sm:   &sMsg,
	}
	go client.resourceHandler()
}

/*******************单独的获取资源状况的websocket handler*******************/

/*******************任务执行的的websocket handler*******************/
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

			j.sendMsg.Type = TRAINNING_CREATION_SEND
			j.sendMsgSignalChan <- struct{}{}

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
			//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			// 延迟关闭请求会话
			//defer cancel()

			// create container
			for _, v := range *j.receiveMsg.Content.SelectedNodes {
				if v.NodeIp+GCS_INFO_CATCH_GRPC_PORT == addrWithPort {
					stream, err := client.DockerContainerRun(context.Background(), &pb.ContainerRunRequestMsg{
						ImageName:     j.sendMsg.Content.ContainerName,
						ContainerName: v.GPUIndex,
						GpuIdx:        j.receiveMsg.Content.ImageName,
					})

					if err != nil {
						slog.Error("rpc DockerContainerStart receive error",
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
							break
						}
						//j.sendMsg.Content.Log = resp.GetImageResp()
						slog.Debug("receive rpc image pull log",
							"UID", j.receiveMsg.Content.IDs.Uid,
							"TID", j.receiveMsg.Content.IDs.Tid,
							"OTHER_MSG", resp.GetRunResp())
					}
					break
				}
			}
			//查看状态

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
		/*接收连接进来的消息*/

		//将新收到的 receiveMsg 更新给当前 job
		job.receiveMsg = receiveMsg

		go func() {
			switch job.receiveMsg.Type {
			//收到信息type是 1，表示获取物理节点状态信息
			case MESSAGE_TYPE_NODE_INFO:
				slog.Debug("get resource info")
				err = resourceInfo(job)
				if err != nil {
					slog.Error("get resourceInfo err", "ERR_MSG", err)
					break
				}
				job.sendMsgSignalChan <- struct{}{}
			case MESSAGE_TYPE_START_CREATION:
				//新建任务需要加入队列中
				slog.Info("create task, job commit to queue")
				//一个任务统一一个 containerName，所以先创建
				job.sendMsg.Content.ContainerName = GetContainerName(strconv.Itoa(job.receiveMsg.Content.IDs.Uid),
					strconv.Itoa(job.receiveMsg.Content.IDs.Tid))
				//提交任务
				h.flowControl.CommitJob(job)
				slog.Info("commit job to job queue success")
				/*
					阻塞函数，直到有向job.DoneChan中写入
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

/*******************任务执行的的websocket handler*******************/

func docker_test() {
	conn, err := grpc.Dial("172.18.127.62:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error",
			"ERR_MSG", err.Error())
	}
	// 延迟关闭连接
	defer conn.Close()

	// 初始化客户端
	client := pb.NewGcsInfoCatchServiceDockerClient(conn)
	// 初始化上下文，设置请求超时时间为1秒
	//ctx, cancel := context.WithTimeout(context.Background(), 999*time.Hour)
	// 延迟关闭请求会话
	//defer cancel()

	// create container
	stream, err := client.DockerContainerRun(context.Background(), &pb.ContainerRunRequestMsg{
		ImageName: "172.18.127.68:80/base-images/ubuntu22_cuda118_python311:v1.1",
		//ImageName:     "172.18.127.68:80/web-images/229end:v1",
		ContainerName: "ssss-1111-222-333",
		GpuIdx:        "2,3,4",
	})

	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("rpc stream server closed")
			break
		}
		if err != nil {
			slog.Error("receive rpc container create log",
				"ERR_MSG", err.Error())
			break
		}
		//j.sendMsg.Content.Log = resp.GetImageResp()
		slog.Debug("receive rpc container create log",
			"OTHER_MSG", resp.GetRunResp())
		if resp.GetRunResp() == "CONTAINER_RUNNING" {
			slog.Debug("receive rpc container create log",
				"CONTAINER_IPS", resp.GetContainerIp())
		}
	}
	//查看状态
	go func() {
		for {
			time.Sleep(3 * time.Second)
			docker_status_test("ssss-1111-222-333")
		}
	}()

	for {

	}
}

// 查看状态
func docker_status_test(containerName string) {
	conn, err := grpc.Dial("172.18.127.62:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error",
			"ERR_MSG", err.Error())
	}
	// 延迟关闭连接
	defer conn.Close()
	client := pb.NewGcsInfoCatchServiceDockerClient(conn)

	stream, err := client.DockerContainerStatus(context.Background(), &pb.StatusRequestMsg{ContainerName: containerName})
	if err != nil {
		slog.Error("DockerContainerStatus error",
			"ERR_MSG", err.Error())
		return
	}

	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("rpc stream server closed")
			break
		}
		if err != nil {
			if resp.GetStatusResp() == "CONTAINER_REMOVE" {
				slog.Debug("receive rpc container status",
					"OTHER_MSG", resp.GetStatusResp())
			} else {
				slog.Error("receive rpc container create log",
					"ERR_MSG", err.Error())
			}
			break
		}
		//j.sendMsg.Content.Log = resp.GetImageResp()
		slog.Debug("receive rpc container status",
			"OTHER_MSG", resp.GetStatusResp())
	}
}
