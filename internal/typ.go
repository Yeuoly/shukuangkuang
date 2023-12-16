package internal

type Shukuangkuang struct {
	stop            chan struct{}
	CPUStatusLoader CPUStatusLoader
}

type ShukuangkuangArgs struct {
	LogicCoreMode bool // use logic cpu cores mode
}
