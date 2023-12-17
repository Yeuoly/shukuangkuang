package internal

type Shukuangkuang struct {
	mode     string // cpu or memory or process
	switched chan struct{}

	CPUStatusLoader    CPUStatusLoader
	MemoryStatusLoader MemoryStatusLoader
}

const (
	MODE_CPU     = "cpu"
	MODE_MEMORY  = "memory"
	MODE_PROCESS = "process"
)

type ShukuangkuangArgs struct {
	LogicCoreMode bool // use logic cpu cores mode
	Help          bool
}
