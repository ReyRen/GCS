/*
package main

import (

	"container/list"
	"fmt"
	"net/http"
	"sync"

)
*/
package main

import (
	"github.com/sevlyar/go-daemon"
	"log/slog"
)

/*func main() {
	flowControl := NewFlowControl()
	myHandler := MyHandler{
		flowControl: flowControl,
	}

	http.Handle("/", &myHandler)

	http.ListenAndServe("172.18.127.64:8080", nil)
}*/

/*type MyHandler struct {
	flowControl *FlowControl
}

//func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("Hello Go"))
//}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("recieve http request")
	job := &Job{
		DoneChan: make(chan struct{}, 1),
		handleJob: func(job *Job) error {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Hello World"))
			return nil
		},
	}

	h.flowControl.CommitJob(job)
	fmt.Println("commit job to job queue success")
	job.WaitDone()
}

type FlowControl struct {
	jobQueue *JobQueue
	wm       *WorkerManager
}

func NewFlowControl() *FlowControl {
	jobQueue := NewJobQueue(10)
	fmt.Println("init job queue success")

	m := NewWorkerManager(jobQueue)
	m.createWorker()
	fmt.Println("init worker success")

	control := &FlowControl{
		jobQueue: jobQueue,
		wm:       m,
	}
	fmt.Println("init flowcontrol success")
	return control
}

func (c *FlowControl) CommitJob(job *Job) {
	c.jobQueue.PushJob(job)
	fmt.Println("commit job success")
}

type Job struct {
	DoneChan  chan struct{}
	handleJob func(j *Job) error
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

type WorkerManager struct {
	jobQueue *JobQueue
}

func NewWorkerManager(jobQueue *JobQueue) *WorkerManager {
	return &WorkerManager{
		jobQueue: jobQueue,
	}
}

func (m *WorkerManager) createWorker() error {

	go func() {
		fmt.Println("start the worker success")
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
}*/

func main() {

	//Setup log system
	log_sys_init()
	//log system ready

	//Setup daemon system
	cntxt := &daemon.Context{
		PidFileName: "gcs.pid",
		PidFilePerm: 0644,
		LogFileName: "./log/gcs.log",
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[gcs]"},
	}
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			slog.Error("cntxt.Search error", "ERR_MSG", err.Error())
		}
		daemon.SendCommands(d)
		return
	}
	d, err := cntxt.Reborn()
	if err != nil {
		slog.Error("cntxt.Reborn error", "ERR_MSG", err.Error())
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	slog.Info("- - - - - - - - - - - - - - -")
	slog.Info("[GCS] started")
	defer func() {
		slog.Info("[GCS] exited")
	}()
	//Daemon system ready

	listen_handler()

	//nvme_sys_init()
}

/*func main() {
ctx := context.Background()

cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
	panic(err)
}

defer cli.Close()

reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
if err != nil {
	panic(err)
}

defer reader.Close()
io.Copy(os.Stdout, reader)

resp, err := cli.ContainerCreate(ctx, &container.Config{
	Image: "alpine",
	Cmd:   []string{"echo", "hello world"},
	Tty:   false,
}, nil, nil, nil, "")
if err != nil {
	panic(err)
}

if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	panic(err)
}

statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
select {
case err := <-errCh:
	if err != nil {
		panic(err)
	}
case <-statusCh:
}

out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
if err != nil {
	panic(err)
}

stdcopy.StdCopy(os.Stdout, os.Stderr, out)
*/
