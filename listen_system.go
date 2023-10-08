package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func listenHandler() {
	flowControl := NewFlowControl()
	myHTTPHandler := MyHandler{
		flowControl: flowControl,
	} //ServeHTTP是MyHandler自己实现的 http 接口方法
	http.Handle("/", &myHTTPHandler)
	err := http.ListenAndServe(GCS_ADDR_WITH_PORT, nil)
	if err != nil {
		slog.Error("ListenAndServe error", "ERR_MSG", err)
	}
}

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
			case <-m.jobQueue.waitJob():
				fmt.Println("get a job from job queue")
				job = m.jobQueue.PopJob()
				fmt.Println("start to execute job")
				job.Execute()
				fmt.Print("execute job done")
				job.Done()
			}
		}
	}()

	return nil
}

func (c *FlowControl) CommitJob(job *Job) {
	c.jobQueue.PushJob(job)
	fmt.Println("commit job success")
}

func (job *Job) Done() {
	job.DoneChan <- struct{}{}
	close(job.DoneChan)
}

func (job *Job) WaitDone() {
	select {
	case <-job.DoneChan:
		return
	}
}

func (job *Job) Execute() error {
	fmt.Println("job start to execute ")
	return job.handleJob(job)
}

func (q *JobQueue) PushJob(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.size++
	if q.size > q.capacity {
		q.RemoveLeastJob()
	}

	q.queue.PushBack(job)

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
		abandonJob.Done()
		q.queue.Remove(back)
	}
}

func (q *JobQueue) waitJob() <-chan struct{} {
	return q.noticeChan
}
