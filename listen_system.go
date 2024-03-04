package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func listenHandler() {
	flowControl := NewFlowControl()
	myHTTPHandler := MyHandler{
		flowControl: flowControl,
	}
	//ServeHTTP是MyHandler自己实现的 http 接口方法
	http.Handle("/", &myHTTPHandler)
	err := http.ListenAndServe(GCS_ADDR_WITH_PORT, nil)
	if err != nil {
		slog.Error("ListenAndServe error", "ERR_MSG", err)
	}
	/*自动调用run.go中，实现的接口函数func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)*/
}

func listenResourceHandler() {
	http.HandleFunc("/resource", func(writer http.ResponseWriter, request *http.Request) {
		resourcehandler(writer, request)
	})
	err := http.ListenAndServe(GCS_RESOURCE_WITH_PORT, nil)
	if err != nil {
		slog.Error("ListenAndServe error", "ERR_MSG", err)
	}
}

/*
支持程序“热重启”
*/
func (job *Job) recordToUpdate(statusID int) error {
	//这里是记录 job
	/*
		handle map:
			{
				"4-129":                //job.receiveMsg.Content.IDs.Uid-job.receiveMsg.Content.IDs.Tid
					[
						"statusID"
						"contaierName", //job.sendMsg.Content.ContainerName
						"LogPathName",  //job.receiveMsg.LogPathName
						"LogAddress",   //job.receiveMsg.Content.LogAddress
						"nodeName1+nodeIP1+gpuIndex,nodeName2+nodeIP2+gpuIndex,.." //job.receiveMsg.Content.SelectedNodes
					]
			}
	*/
	mapKey := strconv.Itoa(job.receiveMsg.Content.IDs.Uid) + "-" + strconv.Itoa(job.receiveMsg.Content.IDs.Tid)
	file, err := os.OpenFile("update.file", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		slog.Error("recordToUpdate OpenFile error", "ERR_MSG", err)
		return err
	}
	defer file.Close()
	if _, ok := UPDATEMAP[mapKey]; ok {
		delete(UPDATEMAP, mapKey)
	}

	UPDATEMAP[mapKey] = append(UPDATEMAP[mapKey], strconv.Itoa(statusID))            //0
	UPDATEMAP[mapKey] = append(UPDATEMAP[mapKey], job.sendMsg.Content.ContainerName) //1
	UPDATEMAP[mapKey] = append(UPDATEMAP[mapKey], job.receiveMsg.LogPathName)        //2
	UPDATEMAP[mapKey] = append(UPDATEMAP[mapKey], job.receiveMsg.Content.LogAddress) //3
	var tmpSlice []string
	for _, v := range *job.receiveMsg.Content.SelectedNodes {
		tmpSlice = append(tmpSlice, v.NodeName+"+"+v.NodeAddress+"+"+v.GPUIndex)
	}
	UPDATEMAP[mapKey] = append(UPDATEMAP[mapKey], strings.Join(tmpSlice, "|")) //4
	dataReady, err := json.Marshal(UPDATEMAP)
	if err != nil {
		slog.Error("recordToUpdate Marshal error", "ERR_MSG", err)
		return err
	}
	file.Write(dataReady)
	return nil
}

func (job *Job) removeToUpdate() error {
	mapKey := strconv.Itoa(job.receiveMsg.Content.IDs.Uid) + "-" + strconv.Itoa(job.receiveMsg.Content.IDs.Tid)
	delete(UPDATEMAP, mapKey)

	file, err := os.OpenFile("update.file", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	defer file.Close()
	if err != nil {
		slog.Error("removeToUpdate OpenFile error", "ERR_MSG", err)
		return err
	}
	dataReady, err := json.Marshal(UPDATEMAP)
	if err != nil {
		slog.Error("removeToUpdate Marshal error", "ERR_MSG", err)
		return err
	}
	file.Write(dataReady)
	return nil
}

func reloadUpdateInfo() error {
	slog.Info("reload status file......")

	file, err := os.OpenFile("update.file", os.O_RDONLY, 0766)
	if err != nil {
		slog.Error("reloadUpdateInfo OpenFile error", "ERR_MSG", err)
		return err
	}
	defer file.Close()

	//这里应该是重新激活一个 job。因为后续的方法都是以 job 为单位的（包括 socket 和 log）
	tmpbyte := make([]byte, 4096)
	total, err := file.Read(tmpbyte)
	if err != nil {
		slog.Error("reloadUpdateInfo Read error", "ERR_MSG", err)
		return err
	}

	err = json.Unmarshal(tmpbyte[:total], &UPDATEMAP)
	if err != nil {
		slog.Error("reloadUpdateInfo Unmarshal error", "ERR_MSG", err)
		return err
	}

	if len(UPDATEMAP) == 0 {
		//如果发现 map 是空，那么说明没有需要 reload 的，或者是有问题了
		slog.Info("update file is empty, nothing to reload")
		return nil
	}
	//遍历 map
	for k, v := range UPDATEMAP {
		receiveMsg := newReceiveMsg()
		slog.Debug("reload:newReceiveMsg initialed ok")
		sendMsg := newSendMsg()
		slog.Debug("reload:newSendMsg initialed ok")

		job := &Job{
			receiveMsg: receiveMsg,
			sendMsg:    sendMsg,
		}
		slog.Debug("reload:Job initialed ok")

		var ids Ids
		ids.Uid, _ = strconv.Atoi(strings.Split(k, "-")[0])
		ids.Tid, _ = strconv.Atoi(strings.Split(k, "-")[1])

		job.receiveMsg.Content.IDs = &ids

		statusID, _ := strconv.Atoi(v[0])
		job.sendMsg.Content.ContainerName = v[1]
		job.receiveMsg.LogPathName = v[2]
		job.receiveMsg.Content.LogAddress = v[3]

		//nodeName1+nodeIP1+gpuIndex｜nodeName2+nodeIP2+gpuIndex,..
		nodeNameIPGPU := strings.Split(v[4], "|")
		//nodeName1+nodeIP1+gpuIndex
		var selectNodesList []SelectNodes
		var selectNodesInfo SelectNodes
		for _, vnodeNameIPGPU := range nodeNameIPGPU {
			selectNodesInfo.NodeName = strings.Split(vnodeNameIPGPU, "+")[0]
			selectNodesInfo.NodeAddress = strings.Split(vnodeNameIPGPU, "+")[1]
			selectNodesInfo.GPUIndex = strings.Split(vnodeNameIPGPU, "+")[2]
			selectNodesList = append(selectNodesList, selectNodesInfo)
		}
		job.receiveMsg.Content.SelectedNodes = &selectNodesList

		err := socketClientCreate(job, statusID)
		if err != nil {
			slog.Error("reloadUpdateInfo socketClientCreate error", "ERR_MSG", err)
			return err
		}
		slog.Debug("reload JOB done", "UID", job.receiveMsg.Content.IDs.Uid,
			"TID", job.receiveMsg.Content.IDs.Tid)
		go func() {
			err := logStoreHandler(job)
			if err != nil {
				return
			}
		}()
		slog.Debug("reload LOG done", "UID", job.receiveMsg.Content.IDs.Uid,
			"TID", job.receiveMsg.Content.IDs.Tid)
	}
	slog.Info("reload JOB and LOG all done")
	return nil
}
