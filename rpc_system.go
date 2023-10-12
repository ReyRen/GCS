package main

import (
	pb "GCS/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"log/slog"
	"net/rpc"
	"time"
)

func rpc_test() {
	client, err := rpc.Dial("tcp", "172.18.127.62:40062")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply string
	err = client.Call("gcs-info-catch-service.DockerContainerRun", "172.18.127.68:80/base-images/deepspeed091_python382_pytorch112_cuda116_pjx:v1.4", &reply)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("get msg from remote", "MSG", reply)
	//fmt.Println(reply)
	/*err = client.Call("gcs-info-catch-service.GoodLuck", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)*/
}

func rpc_handler(rpcServerAddrPort string,

	nvmlRequestInfo NvmlRequestInfo,
	dockerequestInfo DockerRequestInfo) {
	// 连接grpc服务器
	conn, err := grpc.Dial(rpcServerAddrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("grpc.Dial get error", "ERR_MSG", err.Error())
	}
	// 延迟关闭连接
	defer conn.Close()

	// 初始化Greeter服务客户端
	client := pb.NewGcsInfoCatchServiceClient(conn)
	// 初始化上下文，设置请求超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 延迟关闭请求会话
	defer cancel()

	// handler
	switch {

	}
}
