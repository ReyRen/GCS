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

	file, err := os.OpenFile("update.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		slog.Error("recordToUpdate OpenFile error", "ERR_MSG", err)
		return err
	}
	defer file.Close()

	var cInfoMaps []containerInfoMap
	var cInfoMap containerInfoMap
	for _, v := range *job.receiveMsg.Content.SelectedNodes {
		cInfoMap.NodeName = v.NodeName
		cInfoMap.NodeAddress = v.NodeAddress
		cInfoMap.GPUIndex = v.GPUIndex
		cInfoMaps = append(cInfoMaps, cInfoMap)
	}
	mapValue := MapValue{
		StatusID:         strconv.Itoa(statusID),
		ContainerName:    job.sendMsg.Content.ContainerName,
		LogPathName:      job.receiveMsg.LogPathName,
		LogAddress:       job.receiveMsg.Content.LogAddress,
		ContainerInfoMap: &cInfoMaps,
	}
	UPDATEMAP.Store(mapKey, mapValue)
	//将sync.map转换为普通的map然后才能写入到文件
	normalMap := make(map[string]interface{})
	UPDATEMAP.Range(func(key, value interface{}) bool {
		normalMap[key.(string)] = value
		return true
	})
	dataReady, err := json.Marshal(normalMap)
	if err != nil {
		slog.Error("recordToUpdate Marshal error", "ERR_MSG", err)
		return err
	}

	file.Write(dataReady)
	return nil
}

func (job *Job) removeToUpdate() error {
	mapKey := strconv.Itoa(job.receiveMsg.Content.IDs.Uid) + "-" + strconv.Itoa(job.receiveMsg.Content.IDs.Tid)

	UPDATEMAP.Delete(mapKey)

	file, err := os.OpenFile("update.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		slog.Error("removeToUpdate OpenFile error", "ERR_MSG", err)
		return err
	}
	defer file.Close()

	//将sync.map转换为普通的map然后才能写入到文件
	normalMap := make(map[string]interface{})
	UPDATEMAP.Range(func(key, value interface{}) bool {
		normalMap[key.(string)] = value
		return true
	})
	dataReady, err := json.Marshal(normalMap)

	if err != nil {
		slog.Error("removeToUpdate Marshal error", "ERR_MSG", err)
		return err
	}
	file.Write(dataReady)
	return nil
}

func reloadUpdateInfo() error {
	slog.Info("reload status file......")

	file, err := os.OpenFile("update.json", os.O_RDONLY, 0766)
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

	normalMap := make(map[string]MapValue)
	err = json.Unmarshal(tmpbyte[:total], &normalMap)
	if err != nil {
		slog.Error("reloadUpdateInfo Unmarshal error", "ERR_MSG", err)
		return err
	}

	for k, v := range normalMap {
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

		statusID, _ := strconv.Atoi(v.StatusID)
		job.sendMsg.Content.ContainerName = v.ContainerName
		job.receiveMsg.LogPathName = v.LogPathName
		job.receiveMsg.Content.LogAddress = v.LogAddress
		var selectNodesList []SelectNodes
		var selectNodesInfo SelectNodes
		for _, vnodeNameIPGPU := range *(v.ContainerInfoMap) {
			selectNodesInfo.NodeName = vnodeNameIPGPU.NodeName
			selectNodesInfo.NodeAddress = vnodeNameIPGPU.NodeAddress
			selectNodesInfo.GPUIndex = vnodeNameIPGPU.GPUIndex
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

		UPDATEMAP.Store(k, v)
	}

	slog.Info("reload JOB and LOG all done")
	return nil
}
