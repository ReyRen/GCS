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
	"strings"
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
		handleJob: func(j *Job, addrWithPort string, gpuIndex string, master bool) error {
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

			if master {
				//参数补充完整
				originalModuleURL := "--originalModelUrl=" + j.receiveMsg.Content.OriginalModelUrl
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, originalModuleURL)
				ips := "--ip=" + strings.Join(j.receiveMsg.ContainerIps, ",")
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, ips)
				nodes := "--nodes=" + strconv.Itoa(len(*j.receiveMsg.Content.SelectedNodes))
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, nodes)
				modelName := "--modelName=" + j.receiveMsg.Content.ModelName
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, modelName)
				modelType := "--modelType=" + strconv.Itoa(j.receiveMsg.Content.ModelType)
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, modelType)
				frameworkType := "--frameworkType=" + strconv.Itoa(j.receiveMsg.Content.FrameworkType)
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, frameworkType)
				toolBoxName := "--toolBoxName=" + j.receiveMsg.Content.ToolBoxName
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, toolBoxName)
				params := "--params=" + j.receiveMsg.Content.Params
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, params)
				selectedDataset := "--selectedDataset=" + j.receiveMsg.Content.SelectedDataset
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, selectedDataset)
				distributingMethod := "--distributingMethod=" + strconv.Itoa(j.receiveMsg.Content.DistributingMethod)
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, distributingMethod)
				cmd := "--cmd=" + j.receiveMsg.Content.CommandBox
				j.receiveMsg.Paramaters = append(j.receiveMsg.Paramaters, cmd)
			}

			// 初始化客户端
			client := pb.NewGcsInfoCatchServiceDockerClient(conn)

			// create container
			stream, err := client.DockerContainerRun(context.Background(), &pb.ContainerRunRequestMsg{
				ImageName:     j.receiveMsg.Content.ImageName,
				ContainerName: j.sendMsg.Content.ContainerName,
				GpuIdx:        gpuIndex,
				Master:        master,
				Paramaters:    strings.Join(j.receiveMsg.Paramaters, " "),
				//paramaters 组装好
			})
			// 循环获取服务端推送的消息
			for {
				// 通过 Recv() 不断获取服务端send()推送的消息
				resp, err := stream.Recv()
				// err==io.EOF则表示服务端关闭stream了 退出
				if err != nil && err != io.EOF {
					//这是真正的函数执行错误
					//stream服务端返回一种是真正错误，一种是read得到了 EOF 然后return nil
					slog.Error("rpc stream server receive error",
						"ERR_MSG", err.Error(),
						"UID", j.receiveMsg.Content.IDs.Uid,
						"TID", j.receiveMsg.Content.IDs.Tid)
					return err
				}
				if err == io.EOF {
					break
				}
				//如果 error 是 EOF有两种情况，一种是全部都去完成，一种是 PULL 结束了
				//j.sendMsg.Content.Log = resp.GetImageResp()
				slog.Debug("receive container run",
					"UID", j.receiveMsg.Content.IDs.Uid,
					"TID", j.receiveMsg.Content.IDs.Tid,
					"OTHER_MSG", resp.GetRunResp())
				if resp.GetRunResp() == "CONTAINER_RUNNING" {
					j.receiveMsg.ContainerIps = append(j.receiveMsg.ContainerIps, resp.GetContainerIp())
					return nil
				}
			}
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
			case MESSAGE_TYPE_START_CREATION:
				job.sendMsg.Type = 10 //表示训练指令已发送
				job.sendMsgSignalChan <- struct{}{}
				//新建任务需要加入队列中
				slog.Info("create container signal")
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

			//收到信息type是 1，表示获取物理节点状态信息
			case MESSAGE_TYPE_NODE_INFO:
				slog.Debug("get resource info")
				err = resourceInfo(job)
				if err != nil {
					slog.Error("get resourceInfo err", "ERR_MSG", err)
					break
				}
				job.sendMsgSignalChan <- struct{}{}
			case MESSAGE_TYPE_LOG:
				//TODO 获取容器日志
				//TODO 如果发现获取日志出现 EOF，说明日志结束，name 就执行delete操作
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
		ImageName:     "172.18.127.68:80/base-images/ubuntu22_cuda118_python311:v1.1",
		ContainerName: "ssss-1111-222-333",
		GpuIdx:        "2,3,4",
		Master:        true,
		Paramaters:    "--A sss --B sdd --C ss",
	})

	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("111rpc stream server closed")
			break
		}
		if err != nil {
			slog.Error("receive rpc container create log",
				"ERR_MSG", err.Error())
			break
		}
		slog.Debug("receive rpc container create",
			"OTHER_MSG", resp.GetRunResp())
		//j.sendMsg.Content.Log = resp.GetImageResp()

		if resp.GetRunResp() == "CONTAINER_RUNNING" {
			slog.Debug("receive rpc container create",
				"CONTAINER_IPS", resp.GetContainerIp())
		}
		//查看状态
		go func() {
			for {
				time.Sleep(3 * time.Second)
				err := docker_status_test("ssss-1111-222-333")
				if err != nil {
					break
				}
				/*err = docker_exec_test()
				if err != nil {
					break
				}*/
				err = docker_log_test("ssss-1111-222-333")
				if err != nil {
					break
				}
			}
		}()
	}
	for {

	}
}

// 查看日志
func docker_log_test(containerName string) error {
	conn, err := grpc.Dial("172.18.127.62:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error",
			"ERR_MSG", err.Error())
	}
	// 延迟关闭连接
	defer conn.Close()
	client := pb.NewGcsInfoCatchServiceDockerClient(conn)
	stream, err := client.DockerContainerLogs(context.Background(), &pb.LogsRequestMsg{ContainerName: containerName})
	if err != nil {
		slog.Error("DockerContainerLogs error",
			"ERR_MSG", err.Error())
		return err
	}
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("rpc stream server closed")
			return nil
		}
		if err != nil {
			slog.Error("receive rpc container log",
				"ERR_MSG", err.Error())
			return err
		}
		/*if resp.GetStatusResp() == "CONTAINER_REMOVE" {
			slog.Debug("receive rpc container status",
				"OTHER_MSG", resp.GetStatusResp())
		}*/
		//j.sendMsg.Content.Log = resp.GetImageResp()
		slog.Debug("receive rpc container log",
			"OTHER_MSG", resp.GetLogsResp())
		//return nil
	}
}

// 查看状态
func docker_status_test(containerName string) error {
	conn, err := grpc.Dial("172.18.127.62:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error",
			"ERR_MSG", err.Error())
	}
	// 延迟关闭连接
	defer conn.Close()
	client := pb.NewGcsInfoCatchServiceDockerClient(conn)

	stream, err := client.DockerContainerStatus(context.Background(), &pb.StatusRequestMsg{ContainerName: containerName})

	//这错个没用啊！！
	if err != nil {
		slog.Error("DockerContainerStatus error",
			"ERR_MSG", err.Error())
		return err
	}
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			slog.Debug("rpc stream server closed")
			return nil
		}
		if err != nil {
			slog.Error("receive rpc container status",
				"ERR_MSG", err.Error())
			return err
		}
		slog.Debug("receive rpc container status",
			"OTHER_MSG", resp.GetStatusResp())
	}
}
