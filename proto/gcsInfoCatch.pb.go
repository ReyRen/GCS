// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: gcsInfoCatch.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DockerLogStorReqMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogFilePath   string `protobuf:"bytes,1,opt,name=logFilePath,proto3" json:"logFilePath,omitempty"`
	ContainerName string `protobuf:"bytes,2,opt,name=containerName,proto3" json:"containerName,omitempty"`
}

func (x *DockerLogStorReqMsg) Reset() {
	*x = DockerLogStorReqMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DockerLogStorReqMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DockerLogStorReqMsg) ProtoMessage() {}

func (x *DockerLogStorReqMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DockerLogStorReqMsg.ProtoReflect.Descriptor instead.
func (*DockerLogStorReqMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{0}
}

func (x *DockerLogStorReqMsg) GetLogFilePath() string {
	if x != nil {
		return x.LogFilePath
	}
	return ""
}

func (x *DockerLogStorReqMsg) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

type DockerLogStorRespMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogStorResp string `protobuf:"bytes,1,opt,name=logStorResp,proto3" json:"logStorResp,omitempty"`
}

func (x *DockerLogStorRespMsg) Reset() {
	*x = DockerLogStorRespMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DockerLogStorRespMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DockerLogStorRespMsg) ProtoMessage() {}

func (x *DockerLogStorRespMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DockerLogStorRespMsg.ProtoReflect.Descriptor instead.
func (*DockerLogStorRespMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{1}
}

func (x *DockerLogStorRespMsg) GetLogStorResp() string {
	if x != nil {
		return x.LogStorResp
	}
	return ""
}

type ContainerRunRequestMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageName     string `protobuf:"bytes,1,opt,name=imageName,proto3" json:"imageName,omitempty"`
	ContainerName string `protobuf:"bytes,2,opt,name=containerName,proto3" json:"containerName,omitempty"`
	GpuIdx        string `protobuf:"bytes,3,opt,name=gpuIdx,proto3" json:"gpuIdx,omitempty"`
	Master        bool   `protobuf:"varint,4,opt,name=master,proto3" json:"master,omitempty"`
	Paramaters    string `protobuf:"bytes,5,opt,name=paramaters,proto3" json:"paramaters,omitempty"`
}

func (x *ContainerRunRequestMsg) Reset() {
	*x = ContainerRunRequestMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerRunRequestMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerRunRequestMsg) ProtoMessage() {}

func (x *ContainerRunRequestMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerRunRequestMsg.ProtoReflect.Descriptor instead.
func (*ContainerRunRequestMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{2}
}

func (x *ContainerRunRequestMsg) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

func (x *ContainerRunRequestMsg) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

func (x *ContainerRunRequestMsg) GetGpuIdx() string {
	if x != nil {
		return x.GpuIdx
	}
	return ""
}

func (x *ContainerRunRequestMsg) GetMaster() bool {
	if x != nil {
		return x.Master
	}
	return false
}

func (x *ContainerRunRequestMsg) GetParamaters() string {
	if x != nil {
		return x.Paramaters
	}
	return ""
}

type ContainerRunRespondMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RunResp     string `protobuf:"bytes,1,opt,name=runResp,proto3" json:"runResp,omitempty"`
	ContainerIp string `protobuf:"bytes,2,opt,name=containerIp,proto3" json:"containerIp,omitempty"`
}

func (x *ContainerRunRespondMsg) Reset() {
	*x = ContainerRunRespondMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerRunRespondMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerRunRespondMsg) ProtoMessage() {}

func (x *ContainerRunRespondMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerRunRespondMsg.ProtoReflect.Descriptor instead.
func (*ContainerRunRespondMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{3}
}

func (x *ContainerRunRespondMsg) GetRunResp() string {
	if x != nil {
		return x.RunResp
	}
	return ""
}

func (x *ContainerRunRespondMsg) GetContainerIp() string {
	if x != nil {
		return x.ContainerIp
	}
	return ""
}

type DeleteRequestMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContainName string `protobuf:"bytes,1,opt,name=containName,proto3" json:"containName,omitempty"`
}

func (x *DeleteRequestMsg) Reset() {
	*x = DeleteRequestMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequestMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequestMsg) ProtoMessage() {}

func (x *DeleteRequestMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequestMsg.ProtoReflect.Descriptor instead.
func (*DeleteRequestMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRequestMsg) GetContainName() string {
	if x != nil {
		return x.ContainName
	}
	return ""
}

type DeleteRespondMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeleteResp string `protobuf:"bytes,1,opt,name=deleteResp,proto3" json:"deleteResp,omitempty"`
}

func (x *DeleteRespondMsg) Reset() {
	*x = DeleteRespondMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRespondMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRespondMsg) ProtoMessage() {}

func (x *DeleteRespondMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRespondMsg.ProtoReflect.Descriptor instead.
func (*DeleteRespondMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteRespondMsg) GetDeleteResp() string {
	if x != nil {
		return x.DeleteResp
	}
	return ""
}

type StatusRequestMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContainerName string `protobuf:"bytes,1,opt,name=containerName,proto3" json:"containerName,omitempty"`
}

func (x *StatusRequestMsg) Reset() {
	*x = StatusRequestMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusRequestMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusRequestMsg) ProtoMessage() {}

func (x *StatusRequestMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusRequestMsg.ProtoReflect.Descriptor instead.
func (*StatusRequestMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{6}
}

func (x *StatusRequestMsg) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

type StatusRespondMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusResp string `protobuf:"bytes,1,opt,name=statusResp,proto3" json:"statusResp,omitempty"`
}

func (x *StatusRespondMsg) Reset() {
	*x = StatusRespondMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusRespondMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusRespondMsg) ProtoMessage() {}

func (x *StatusRespondMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusRespondMsg.ProtoReflect.Descriptor instead.
func (*StatusRespondMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{7}
}

func (x *StatusRespondMsg) GetStatusResp() string {
	if x != nil {
		return x.StatusResp
	}
	return ""
}

type LogsRequestMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContainerName string `protobuf:"bytes,1,opt,name=containerName,proto3" json:"containerName,omitempty"`
}

func (x *LogsRequestMsg) Reset() {
	*x = LogsRequestMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogsRequestMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsRequestMsg) ProtoMessage() {}

func (x *LogsRequestMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsRequestMsg.ProtoReflect.Descriptor instead.
func (*LogsRequestMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{8}
}

func (x *LogsRequestMsg) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

type LogsRespondMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogsResp string `protobuf:"bytes,1,opt,name=logsResp,proto3" json:"logsResp,omitempty"`
}

func (x *LogsRespondMsg) Reset() {
	*x = LogsRespondMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogsRespondMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsRespondMsg) ProtoMessage() {}

func (x *LogsRespondMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsRespondMsg.ProtoReflect.Descriptor instead.
func (*LogsRespondMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{9}
}

func (x *LogsRespondMsg) GetLogsResp() string {
	if x != nil {
		return x.LogsResp
	}
	return ""
}

// **************************NVML Message*****************************
type NvmlInfoReuqestMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndexID string `protobuf:"bytes,1,opt,name=indexID,proto3" json:"indexID,omitempty"`
}

func (x *NvmlInfoReuqestMsg) Reset() {
	*x = NvmlInfoReuqestMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NvmlInfoReuqestMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NvmlInfoReuqestMsg) ProtoMessage() {}

func (x *NvmlInfoReuqestMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NvmlInfoReuqestMsg.ProtoReflect.Descriptor instead.
func (*NvmlInfoReuqestMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{10}
}

func (x *NvmlInfoReuqestMsg) GetIndexID() string {
	if x != nil {
		return x.IndexID
	}
	return ""
}

type NvmlInfoRespondMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndexID         []int32  `protobuf:"varint,1,rep,packed,name=indexID,proto3" json:"indexID,omitempty"`
	UtilizationRate []uint32 `protobuf:"varint,2,rep,packed,name=utilizationRate,proto3" json:"utilizationRate,omitempty"`
	MemRate         []uint32 `protobuf:"varint,3,rep,packed,name=memRate,proto3" json:"memRate,omitempty"`
	Temperature     []uint32 `protobuf:"varint,4,rep,packed,name=temperature,proto3" json:"temperature,omitempty"`
	Occupied        []uint32 `protobuf:"varint,5,rep,packed,name=occupied,proto3" json:"occupied,omitempty"` //通过 runningprocess 有无来进行判断
}

func (x *NvmlInfoRespondMsg) Reset() {
	*x = NvmlInfoRespondMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcsInfoCatch_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NvmlInfoRespondMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NvmlInfoRespondMsg) ProtoMessage() {}

func (x *NvmlInfoRespondMsg) ProtoReflect() protoreflect.Message {
	mi := &file_gcsInfoCatch_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NvmlInfoRespondMsg.ProtoReflect.Descriptor instead.
func (*NvmlInfoRespondMsg) Descriptor() ([]byte, []int) {
	return file_gcsInfoCatch_proto_rawDescGZIP(), []int{11}
}

func (x *NvmlInfoRespondMsg) GetIndexID() []int32 {
	if x != nil {
		return x.IndexID
	}
	return nil
}

func (x *NvmlInfoRespondMsg) GetUtilizationRate() []uint32 {
	if x != nil {
		return x.UtilizationRate
	}
	return nil
}

func (x *NvmlInfoRespondMsg) GetMemRate() []uint32 {
	if x != nil {
		return x.MemRate
	}
	return nil
}

func (x *NvmlInfoRespondMsg) GetTemperature() []uint32 {
	if x != nil {
		return x.Temperature
	}
	return nil
}

func (x *NvmlInfoRespondMsg) GetOccupied() []uint32 {
	if x != nil {
		return x.Occupied
	}
	return nil
}

var File_gcsInfoCatch_proto protoreflect.FileDescriptor

var file_gcsInfoCatch_proto_rawDesc = []byte{
	0x0a, 0x12, 0x67, 0x63, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x43, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5d, 0x0a, 0x13, 0x44,
	0x6f, 0x63, 0x6b, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x4d,
	0x73, 0x67, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6c, 0x6f, 0x67, 0x46, 0x69, 0x6c, 0x65,
	0x50, 0x61, 0x74, 0x68, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x38, 0x0a, 0x14, 0x44, 0x6f,
	0x63, 0x6b, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x4d,
	0x73, 0x67, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x53, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6c, 0x6f, 0x67, 0x53, 0x74, 0x6f, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x22, 0xac, 0x01, 0x0a, 0x16, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x73, 0x67, 0x12,
	0x1c, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a,
	0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x70, 0x75, 0x49, 0x64, 0x78, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x70, 0x75, 0x49, 0x64, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x6d,
	0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x6d, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x61, 0x74, 0x65, 0x72,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x61, 0x74,
	0x65, 0x72, 0x73, 0x22, 0x54, 0x0a, 0x16, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x52, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x18, 0x0a,
	0x07, 0x72, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x72, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x61,
	0x69, 0x6e, 0x65, 0x72, 0x49, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f,
	0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x70, 0x22, 0x34, 0x0a, 0x10, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x20, 0x0a,
	0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x32, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64,
	0x4d, 0x73, 0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x22, 0x38, 0x0a, 0x10, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61,
	0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x32, 0x0a,
	0x10, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73,
	0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x22, 0x36, 0x0a, 0x0e, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x4d, 0x73, 0x67, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2c, 0x0a, 0x0e, 0x4c, 0x6f, 0x67,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x6c,
	0x6f, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c,
	0x6f, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x22, 0x2e, 0x0a, 0x12, 0x4e, 0x76, 0x6d, 0x6c, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x75, 0x71, 0x65, 0x73, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x18, 0x0a,
	0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x44, 0x22, 0xb0, 0x01, 0x0a, 0x12, 0x4e, 0x76, 0x6d, 0x6c,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x18,
	0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x44, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52,
	0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x0f, 0x75, 0x74, 0x69, 0x6c,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0d, 0x52, 0x0f, 0x75, 0x74, 0x69, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61,
	0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x52, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0d, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x52, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0d, 0x52, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x69, 0x65, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0d,
	0x52, 0x08, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x69, 0x65, 0x64, 0x32, 0xf7, 0x03, 0x0a, 0x19, 0x47,
	0x63, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x43, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x12, 0x44, 0x6f, 0x63, 0x6b,
	0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x75, 0x6e, 0x12, 0x1d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x1d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52,
	0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x30, 0x01,
	0x12, 0x4d, 0x0a, 0x15, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d,
	0x73, 0x67, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x30, 0x01, 0x12,
	0x4d, 0x0a, 0x15, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x73,
	0x67, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x30, 0x01, 0x12, 0x47,
	0x0a, 0x13, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x4c, 0x6f, 0x67, 0x73, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f,
	0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x15, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64,
	0x4d, 0x73, 0x67, 0x22, 0x00, 0x30, 0x01, 0x12, 0x4a, 0x0a, 0x0d, 0x44, 0x6f, 0x63, 0x6b, 0x65,
	0x72, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x6f, 0x72, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x6f, 0x72, 0x52, 0x65,
	0x71, 0x4d, 0x73, 0x67, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x4d, 0x73,
	0x67, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x13, 0x4e, 0x76, 0x6d, 0x6c, 0x55, 0x74, 0x69, 0x6c, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x76, 0x6d, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x75, 0x71, 0x65,
	0x73, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x76,
	0x6d, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x4d, 0x73, 0x67,
	0x22, 0x00, 0x30, 0x01, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gcsInfoCatch_proto_rawDescOnce sync.Once
	file_gcsInfoCatch_proto_rawDescData = file_gcsInfoCatch_proto_rawDesc
)

func file_gcsInfoCatch_proto_rawDescGZIP() []byte {
	file_gcsInfoCatch_proto_rawDescOnce.Do(func() {
		file_gcsInfoCatch_proto_rawDescData = protoimpl.X.CompressGZIP(file_gcsInfoCatch_proto_rawDescData)
	})
	return file_gcsInfoCatch_proto_rawDescData
}

var file_gcsInfoCatch_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_gcsInfoCatch_proto_goTypes = []interface{}{
	(*DockerLogStorReqMsg)(nil),    // 0: proto.DockerLogStorReqMsg
	(*DockerLogStorRespMsg)(nil),   // 1: proto.DockerLogStorRespMsg
	(*ContainerRunRequestMsg)(nil), // 2: proto.ContainerRunRequestMsg
	(*ContainerRunRespondMsg)(nil), // 3: proto.ContainerRunRespondMsg
	(*DeleteRequestMsg)(nil),       // 4: proto.DeleteRequestMsg
	(*DeleteRespondMsg)(nil),       // 5: proto.DeleteRespondMsg
	(*StatusRequestMsg)(nil),       // 6: proto.StatusRequestMsg
	(*StatusRespondMsg)(nil),       // 7: proto.StatusRespondMsg
	(*LogsRequestMsg)(nil),         // 8: proto.LogsRequestMsg
	(*LogsRespondMsg)(nil),         // 9: proto.LogsRespondMsg
	(*NvmlInfoReuqestMsg)(nil),     // 10: proto.NvmlInfoReuqestMsg
	(*NvmlInfoRespondMsg)(nil),     // 11: proto.NvmlInfoRespondMsg
}
var file_gcsInfoCatch_proto_depIdxs = []int32{
	2,  // 0: proto.GcsInfoCatchServiceDocker.DockerContainerRun:input_type -> proto.ContainerRunRequestMsg
	4,  // 1: proto.GcsInfoCatchServiceDocker.DockerContainerDelete:input_type -> proto.DeleteRequestMsg
	6,  // 2: proto.GcsInfoCatchServiceDocker.DockerContainerStatus:input_type -> proto.StatusRequestMsg
	8,  // 3: proto.GcsInfoCatchServiceDocker.DockerContainerLogs:input_type -> proto.LogsRequestMsg
	0,  // 4: proto.GcsInfoCatchServiceDocker.DockerLogStor:input_type -> proto.DockerLogStorReqMsg
	10, // 5: proto.GcsInfoCatchServiceDocker.NvmlUtilizationRate:input_type -> proto.NvmlInfoReuqestMsg
	3,  // 6: proto.GcsInfoCatchServiceDocker.DockerContainerRun:output_type -> proto.ContainerRunRespondMsg
	5,  // 7: proto.GcsInfoCatchServiceDocker.DockerContainerDelete:output_type -> proto.DeleteRespondMsg
	7,  // 8: proto.GcsInfoCatchServiceDocker.DockerContainerStatus:output_type -> proto.StatusRespondMsg
	9,  // 9: proto.GcsInfoCatchServiceDocker.DockerContainerLogs:output_type -> proto.LogsRespondMsg
	1,  // 10: proto.GcsInfoCatchServiceDocker.DockerLogStor:output_type -> proto.DockerLogStorRespMsg
	11, // 11: proto.GcsInfoCatchServiceDocker.NvmlUtilizationRate:output_type -> proto.NvmlInfoRespondMsg
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_gcsInfoCatch_proto_init() }
func file_gcsInfoCatch_proto_init() {
	if File_gcsInfoCatch_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gcsInfoCatch_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DockerLogStorReqMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DockerLogStorRespMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerRunRequestMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerRunRespondMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequestMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRespondMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusRequestMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusRespondMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogsRequestMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogsRespondMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NvmlInfoReuqestMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gcsInfoCatch_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NvmlInfoRespondMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gcsInfoCatch_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gcsInfoCatch_proto_goTypes,
		DependencyIndexes: file_gcsInfoCatch_proto_depIdxs,
		MessageInfos:      file_gcsInfoCatch_proto_msgTypes,
	}.Build()
	File_gcsInfoCatch_proto = out.File
	file_gcsInfoCatch_proto_rawDesc = nil
	file_gcsInfoCatch_proto_goTypes = nil
	file_gcsInfoCatch_proto_depIdxs = nil
}
