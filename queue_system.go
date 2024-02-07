package main

import (
	"log/slog"
	"strings"
)

func NewFlowControl() *FlowControl {
	jobQueue := NewJobQueue(1024)
	slog.Debug("init job queue success")

	nm := NewWorkerManager(jobQueue)
	_ = nm.createWorker()
	slog.Debug("init worker success")

	control := &FlowControl{
		jobQueue: jobQueue,
		wm:       nm,
	}
	slog.Debug("init flowcontrol success")
	return control
}

func (m *WorkerManager) createWorker() error {
	go func() {
		slog.Debug("start the worker success")
		var job *Job

		for {
			select {
			case <-m.jobQueue.waitJob(): // 执行完下面的语句才会回来继续监听waitJob()的 notice channel
				slog.Debug("get a job notification from job queue")
				job = m.jobQueue.PopJob()
				slog.Debug("pop the job",
					"UID", job.receiveMsg.Content.IDs.Uid,
					"TID", job.receiveMsg.Content.IDs.Tid)

				//判断资源是否满足
				slog.Debug("Check resource free or not")
				free := true
				for _, v := range *job.receiveMsg.Content.SelectedNodes {
					//先获取一下 每个容器GPU 的个数
					job.receiveMsg.GpuCountPerContainer = len(strings.Split(v.GPUIndex, ","))
					// true说明有被占用
					if checkGPUOccupiedOrNot(v.NodeAddress, v.GPUIndex) {
						free = false
						job.sendMsg.Type = WS_STATUS_BACK_RESOURCE_INSUFFICIENT
						//job.sendMsgSignalChan <- struct{}{}
						break
					} else {
						//job.sendMsgSignalChan <- struct{}{}
						slog.Debug("selected gpu free", "NODE_ADDR", v.NodeAddress)
					}
					job.sendMsgSignalChan <- struct{}{}
				}
				if !free {
					//说明资源不满足，那么就必须停止任务创建
					slog.Debug("resource cannot use, task over")
					job.Done()
				} else {
					err := socketClientCreate(job, SOCKET_STATUS_BACK_CREATE_START)
					if err != nil {
						slog.Debug("socketClientCreate error in  container creating")
						return
					}
					slog.Debug("socket send:container creating")
					//给 websocket 发送创建容器开始
					job.sendMsg.Type = WS_STATUS_BACK_CREATING
					/*job.sendMsgSignalChan <- struct{}{}*/
					err = job.Execute()
					if err != nil {
						//说明创建任务的过程中，有问题
						slog.Debug("execute job error, execute delete containers",
							"UID", job.receiveMsg.Content.IDs.Uid,
							"TID", job.receiveMsg.Content.IDs.Tid)
						err := socketClientCreate(job, SOCKET_STATUS_BACK_CREATE_FAILED)
						if err != nil {
							slog.Debug("socketClientCreate error in  container create failed")
							return
						}
						slog.Debug("socket send:container create failed")
						job.sendMsg.Type = WS_STATUS_BACK_CREATE_FAILED
						/*job.sendMsgSignalChan <- struct{}{}*/
						//将该任务的容器不管有没有创建成功都删除一遍
						for _, v := range *job.receiveMsg.Content.SelectedNodes {
							err := dockerDeleteHandler(v.NodeAddress, job.sendMsg.Content.ContainerName)
							if err != nil {
								slog.Error("dockerDeleteHandler get error")
							}
						}
						job.Done()
					} else {
						err := socketClientCreate(job, SOCKET_STATUS_BACK_TRAINNING)
						if err != nil {
							slog.Debug("socketClientCreate error in  container running")
							return
						}
						slog.Debug("socket send:container running")
						job.sendMsg.Type = WS_STATUS_BACK_TRAINNING
						/*job.sendMsg.Content.ContainerName = job.sendMsg.Content.ContainerName*/
						err = resourceInfo(job)
						if err != nil {
							slog.Error("get resourceInfo err", "ERR_MSG", err)
						}
						/*job.sendMsgSignalChan <- struct{}{}*/
						slog.Debug("execute job done",
							"UID", job.receiveMsg.Content.IDs.Uid,
							"TID", job.receiveMsg.Content.IDs.Tid)
						job.Done()
					}
				}
			}
		}
	}()
	return nil
}

func (c *FlowControl) CommitJob(job *Job) {
	c.jobQueue.PushJob(job)
	slog.Debug("commit job success",
		"UID", job.receiveMsg.Content.IDs.Uid,
		"TID", job.receiveMsg.Content.IDs.Tid)
}

func (job *Job) Done() {
	slog.Debug("job done",
		"UID", job.receiveMsg.Content.IDs.Uid,
		"TID", job.receiveMsg.Content.IDs.Tid)
	job.DoneChan <- struct{}{}
}

func (job *Job) WaitDone() {
	select {
	case <-job.DoneChan:
		return
	}
}

func (job *Job) Execute() error {
	slog.Debug("start execute job",
		"UID", job.receiveMsg.Content.IDs.Uid,
		"TID", job.receiveMsg.Content.IDs.Tid)
	master := false
	for k, v := range *job.receiveMsg.Content.SelectedNodes {
		if k+1 == len(*job.receiveMsg.Content.SelectedNodes) {
			//说明是最后一个容器的创建
			master = true
			job.receiveMsg.Content.LogAddress = v.NodeAddress
		}
		slog.Debug("grpc execute to server",
			"GRPC_SERVER", v.NodeAddress+GCS_INFO_CATCH_GRPC_PORT,
			"UID", job.receiveMsg.Content.IDs.Uid,
			"TID", job.receiveMsg.Content.IDs.Tid)
		err := job.handleJob(job, v.NodeAddress+GCS_INFO_CATCH_GRPC_PORT, v.GPUIndex, master)
		if err != nil {
			// 说明单个任务创建失败
			slog.Error("grpc execute to server: execute error!!",
				"GRPC_SERVER", v.NodeAddress+GCS_INFO_CATCH_GRPC_PORT,
				"UID", job.receiveMsg.Content.IDs.Uid,
				"TID", job.receiveMsg.Content.IDs.Tid)
			return err
		}
	}
	return nil
}

func (q *JobQueue) PushJob(job *Job) {
	slog.Debug("11111111111111111111111111111111 before q.mu.Lock()")
	//q.mu.Lock()
	slog.Debug("11111111111111111111111111111111 after q.mu.Lock()")
	//defer q.mu.Unlock()
	q.size++
	slog.Debug("queue size=", q.size, "queue capacity=", q.capacity)
	if q.size > q.capacity {
		slog.Debug("queue is overwhelmed, remove least job")
		q.RemoveLeastJob()
	}

	q.queue.PushBack(job)
	slog.Debug("push the job to the jobqueue",
		"UID", job.receiveMsg.Content.IDs.Uid,
		"TID", job.receiveMsg.Content.IDs.Tid)

	/*
		q.noticeChan是个带有 1 个缓冲区的 channel
		当任务 1 执行任务中
		  任务 2 来了，执行完 push 后，可以写入q.noticeChan中，无须阻塞，但是此时q.noticeChan还没轮到监听，所以阻塞在 pop执行任务上
			任务 3 来了，因为noticeChan缓冲区的还没读，所以任务 3 会阻塞在写入q.noticeChan的时候
	*/
	q.noticeChan <- struct{}{} //struct{}表示类型，struct{}{}表示实例化
}

func (q *JobQueue) PopJob() *Job {
	//q.mu.Lock()
	//defer q.mu.Unlock()

	if q.size == 0 {
		return nil
	}

	q.size--
	return q.queue.Remove(q.queue.Front()).(*Job)
}

func (q *JobQueue) RemoveLeastJob() {
	if q.queue.Len() != 0 {
		back := q.queue.Back()
		abandonJob := back.Value.(*Job)
		abandonJob.Done() //释放waitdone 的阻塞
		q.queue.Remove(back)
		slog.Debug("remove the lease job",
			"UID", abandonJob.receiveMsg.Content.IDs.Uid,
			"TID", abandonJob.receiveMsg.Content.IDs.Tid)
	}
}

func (q *JobQueue) waitJob() <-chan struct{} {
	return q.noticeChan
}
