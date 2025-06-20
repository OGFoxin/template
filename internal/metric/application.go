package metric

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/process"
)

func (m *metrics) GetCpuInfo() []float64 {
	cpuData, err := cpu.Percent(0, true)

	if err != nil {
		cpuData = nil
	}

	return cpuData
}

func (m *metrics) GetMemoryInfo() *process.MemoryInfoStat {
	// change to build
	processName := "go.exe"
	processes, _ := process.Processes()
	var memInfo *process.MemoryInfoStat

	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}
		if name == processName {
			memInfo, err = proc.MemoryInfo()
			if err != nil {
				continue
			}

		}
	}

	return memInfo
}
