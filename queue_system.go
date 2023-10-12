package main

import (
	"log/slog"
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
				job.Execute()
				slog.Debug("execute job done",
					"UID", job.receiveMsg.Content.IDs.Uid,
					"TID", job.receiveMsg.Content.IDs.Tid)
				job.Done()
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
	job.DoneChan <- struct{}{}
	slog.Debug("job done",
		"UID", job.receiveMsg.Content.IDs.Uid,
		"TID", job.receiveMsg.Content.IDs.Tid)
	if _, ok := <-job.DoneChan; ok {
		close(job.DoneChan)
	}
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
	return job.handleJob(job, "172.18.127.62:40062")
}

func (q *JobQueue) PushJob(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.size++
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
	q.mu.Lock()
	defer q.mu.Unlock()

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
