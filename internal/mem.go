package internal

import (
	"sort"
	"sync"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

const (
	MAX_MEM_COUNT = 256

	MAX_COCURRENT_PROCESS = 24
)

type ProcessInfo struct {
	Pid         int32
	Name        string
	MemoryUsage uint64
}

type byMemoryUsage []*ProcessInfo

func (a byMemoryUsage) Len() int           { return len(a) }
func (a byMemoryUsage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byMemoryUsage) Less(i, j int) bool { return a[i].MemoryUsage > a[j].MemoryUsage }

type MemoryStatusLoaderIface interface {
	// GetMemoryStatus returns the memory status, returns []usage, []topMemProc
	GetMemoryStatus() ([]uint64, []string, error)
	// GetMaxMemory returns the max memory
	GetMaxMemory() (uint64, error)
	// Init initializes the MemoryStatusLoader
	Init()
	// LaunchGoRoutine launches a goroutine to do the task
	LaunchGoRoutine()
	// StopGoRoutine stops the goroutine
	StopGoRoutine()
}

type MemoryStatusLoader struct {
	usage     []uint64
	maxMemory uint64
	tasks     chan func()
	idx       int
	processes []*ProcessInfo

	MemoryStatusLoaderIface
}

func (c *MemoryStatusLoader) Init() {
	c.usage = make([]uint64, MAX_MEM_COUNT)
	mem, err := mem.VirtualMemory()
	if err != nil {
		c.maxMemory = 2 << 30
	} else {
		c.maxMemory = mem.Total
	}
}

func (c *MemoryStatusLoader) LaunchGoRoutine() {
	c.tasks = make(chan func(), MAX_COCURRENT_PROCESS)
	for i := 0; i < MAX_COCURRENT_PROCESS; i++ {
		go func() {
			for task := range c.tasks {
				task()
			}
		}()
	}
}

func (c *MemoryStatusLoader) StopGoRoutine() {
	close(c.tasks)
}

func (c *MemoryStatusLoader) GetMemoryStatus() ([]uint64, []*ProcessInfo, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, nil, err
	}

	// move all the data to the right one step
	for i := MAX_MEM_COUNT - 1; i > 0; i-- {
		c.usage[i] = c.usage[i-1]
	}

	// update the first one
	c.usage[0] = memInfo.Used

	if c.idx%10 == 0 {
		// range over the usage to find the topMemProc
		processes, err := process.Processes()
		if err != nil {
			return nil, nil, err
		}

		wg := sync.WaitGroup{}
		topMemProc := make([]*ProcessInfo, 0)
		lock := sync.Mutex{}

		for _, p := range processes {
			wg.Add(1)
			p := p
			c.tasks <- func() {
				memInfo, err := p.MemoryInfo()
				defer wg.Done()
				if err != nil {
					return
				}

				name, err := p.Name()
				if err != nil {
					return
				}

				p := &ProcessInfo{
					Pid:         p.Pid,
					Name:        name,
					MemoryUsage: memInfo.RSS,
				}

				lock.Lock()
				topMemProc = append(topMemProc, p)
				lock.Unlock()
			}
		}

		wg.Wait()

		sort.Sort(byMemoryUsage(topMemProc))

		c.processes = topMemProc
	}

	c.idx++

	return c.usage[:], c.processes, nil
}

func (c *MemoryStatusLoader) GetMaxMemory() (uint64, error) {
	return c.maxMemory, nil
}

func NewMemoryStatusLoader() MemoryStatusLoader {
	return MemoryStatusLoader{}
}
