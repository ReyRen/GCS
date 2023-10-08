package main

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"log/slog"
)

func nvme_sys_init() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		slog.Error("Unable to initialize NVML", "ERR_MSG", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			slog.Error("Unable to shutdown NVML", "ERR_MSG", nvml.ErrorString(ret))
		}
	}()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		slog.Error("Unable to get device count", "ERR_MSG", nvml.ErrorString(ret))
	}

	for i := 0; i < count; i++ {
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			slog.Error("Unable to get device at index", "INDEX", i, "ERR_MSG", nvml.ErrorString(ret))
		}
		uuid, ret := device.GetUUID()
		if ret != nvml.SUCCESS {
			slog.Error("Unable to get uuid of device at index", "INDEX", i, "ERR_MSG", nvml.ErrorString(ret))
		}

		pInfo, ret := device.GetComputeRunningProcesses()
		if ret != nvml.SUCCESS {
			slog.Info("GetComputeRunningProcesses empty")
		} else {
			for _, value := range pInfo {
				slog.Info("Output GPU info",
					"UUID", uuid,
					"COMPUTE_INSTANCE_ID", value.ComputeInstanceId,
					"GPU_INSTANCE_ID", value.GpuInstanceId,
					"PID", value.Pid,
					"USED_MEM", value.UsedGpuMemory/1024/1024)
			}
		}
	}
}
