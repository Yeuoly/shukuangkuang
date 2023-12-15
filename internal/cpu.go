package internal

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type CPUStatusLoader interface {
	// GetCPUStatus returns the cpu status
	GetCPUStatus() ([][]float64, error)
	// GetCPUCount returns the cpu count
	GetCPUCount() (int, error)
	// Init initializes the CPUStatusLoader
	Init()
}

const (
	MAX_CPU_COUNT = 64
)

type SingleCPU struct {
	percent [MAX_CPU_COUNT]float64
	CPUStatusLoader
}

type LogicCPU struct {
	percent [][MAX_CPU_COUNT]float64
	CPUStatusLoader
}

func (c *SingleCPU) GetCPUStatus() ([][]float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	// move all the data to the right one step
	for i := MAX_CPU_COUNT - 1; i > 0; i-- {
		c.percent[i] = c.percent[i-1]
	}
	// update the first one
	c.percent[0] = percent[0]

	result := make([][]float64, 1)
	result[0] = c.percent[:]
	return result, err
}

func (c *SingleCPU) GetCPUCount() (int, error) {
	return 1, nil
}

func (c *SingleCPU) Init() {
	c.percent = [MAX_CPU_COUNT]float64{}
}

func (c *LogicCPU) GetCPUStatus() ([][]float64, error) {
	percent, err := cpu.Percent(time.Second, true)
	// move all the data of per cpu to the right one step
	for i := len(percent) - 1; i > 0; i-- {
		for j := MAX_CPU_COUNT - 1; j > 0; j-- {
			c.percent[i][j] = c.percent[i][j-1]
		}
	}
	// update the first one of per cpu
	for i := 0; i < len(percent); i++ {
		c.percent[i][0] = percent[i]
	}

	result := make([][]float64, len(percent))
	for i := 0; i < len(percent); i++ {
		result[i] = c.percent[i][:]
	}

	return result, err
}

func (c *LogicCPU) GetCPUCount() (int, error) {
	return cpu.Counts(true)
}

func (c *LogicCPU) Init() {
	count, _ := c.GetCPUCount()
	c.percent = make([][MAX_CPU_COUNT]float64, count)
}

func NewCPUStatusLoader(logicCoreMode bool) CPUStatusLoader {
	if logicCoreMode {
		return &LogicCPU{}
	}
	return &SingleCPU{}
}
