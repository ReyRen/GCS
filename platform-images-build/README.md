# 特定任务镜像构造

该使用说明的目的是，结合自己的需求构建出符合分布式并行训练平台使用的镜像

## 01-base-cuda
该目录下的 Dockerfile 将创建出一个拥有 cuda 、 cudnn 、一些平台使用 IB的相关工具以及做一些免密等相关工作的基础镜像

## 02-miniconda
基于 01-base-cuda 中生成的基础镜像，增加了 minconda 和对应python版本的安装

## 03-pytorch
基于 02-miniconda 中生成的基础镜像，增加了对应版本的 pytorch

......

**说明**：使用者可根据该结果，进行回溯到相对应的 dockerfile 中进行更新和定制。但其中**01-base-cuda**中除了换 From 源外，其余安装的工具为平台
必须。