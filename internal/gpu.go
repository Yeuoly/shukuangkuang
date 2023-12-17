package internal

const (
	MAX_GPU_COUNT = 64
)

type GPUStatusLoader interface {
	// GetGPUStatus returns the gpu status, returns [][]usage, [][]memory, [][]temperature
	GetGPUStatus() ([][]float64, [][]uint64, [][]float64, error)
	// GetGPUCount returns the gpu count
	GetGPUCount() (int, error)
}

type SingleGPU struct {
	percent     [MAX_GPU_COUNT]float64
	memory      [MAX_GPU_COUNT]uint64
	temperature [MAX_GPU_COUNT]float64
	MaxMemory   uint64

	GPUStatusLoader
}

type LogicGPU struct {
	percent     [][MAX_GPU_COUNT]float64
	memory      [][MAX_GPU_COUNT]uint64
	temperature [][MAX_GPU_COUNT]float64
	MaxMemory   uint64

	GPUStatusLoader
}

func getGPUStatus() ([][]float64, [][]uint64, [][]float64, error) {
	return nil, nil, nil, nil
}
