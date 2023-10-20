# GCS
graceful container schedule for AI trainning, support Infiniband, GDRDMA

## 简介
这个程序需要和 grpc 客户端程序进行交互 https://github.com/ReyRen/GCS-Info-Catch。 如图为两者之间的关系。

GCS（Graceful Container Schedule）主要是摒弃了在分布式训练中使用kubernetes/slurm等训练集群，完全轻量化的部署，
轻量化的使用。通过这样，可以除去冗余的部分，精准做到以下特点：
* 训练性能损耗最小
* 训练部署方式最简单
* “集群”极简部署（其实就是运行环境部署）
* 任务隔离（看到即可用）
* （我们环境）全通路（1:1:1）IB 链路训练，实现 GDRDMA功能
* 任务等待队列，保证创建任务不冲突
* GRPC 过程调度，加入节点很方便

其中用到了nvml-go和docker的 API。NVML 获取精准的物理机GPU 信息，docker 用来对创建的任务进行全流程追踪。


## 大体使用方式
1. 安装 https://golang.google.cn/doc/install. 这里使用的是v1.21.2
2. 环境变量设置为(.bashrc下)
```shell
  export PATH=$PATH:/usr/local/go/bin
  export GOPROXY="https://goproxy.cn"
  export GO111MODULE=on
  export PATH="$PATH:$(go env GOPATH)/bin"
```
3. 每个节点上安装好 ofed 驱动以及 nvidia GPU 驱动
4. 相应的镜像仓库/http、ftp /大容量告诉存储挂载，这些都准备好
5. Docker 默认安装好（最新版本），nvidia-docker 安装好并且设置为默认 runtime
6. 执行 make run即可