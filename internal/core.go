package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func (c *Shukuangkuang) Run(args ShukuangkuangArgs) {
	if args.LogicCoreMode {
		c.CPUStatusLoader = NewCPUStatusLoader(true)
	} else {
		c.CPUStatusLoader = NewCPUStatusLoader(false)
	}

	c.MemoryStatusLoader = NewMemoryStatusLoader()

	c.CPUStatusLoader.Init()
	c.MemoryStatusLoader.Init()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	c.mode = MODE_CPU
	c.switched = make(chan struct{}, 1)

	c.event()
	c.run()
}

func (c *Shukuangkuang) event() {
	go func() {
		uiEvents := ui.PollEvents()
		for event := range uiEvents {
			switch event.ID {
			case "q", "<C-c>":
				ui.Close()
				os.Exit(0)
			case "m":
				c.mode = MODE_MEMORY
				c.switched <- struct{}{}
			case "c":
				c.mode = MODE_CPU
				c.switched <- struct{}{}
			}
		}
	}()
}

func (c *Shukuangkuang) run() {
	for {
		switch c.mode {
		case MODE_CPU:
			c.renderCpu()
		case MODE_MEMORY:
			c.renderMemory()
		case MODE_PROCESS:
		}
	}
}

func (c *Shukuangkuang) renderLoop(timeout time.Duration, draw func()) {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	draw()
	for {
		select {
		case <-ticker.C:
			draw()
		case <-c.switched:
			return
		}
	}
}

func (c *Shukuangkuang) renderCpu() {
	// get total width and height
	w, h := ui.TerminalDimensions()
	// compute the count of cpu cores
	count, err := c.CPUStatusLoader.GetCPUCount()
	if err != nil {
		log.Fatalf("failed to get cpu count: %v", err)
	}

	// compute the width and height of each box
	box_width, box_height := CaculateBoxSizes(count, w, h)

	// compute the count of rows and cols
	cols := w / box_width

	plots := make([]*widgets.SparklineGroup, count)

	for i := 0; i < count; i++ {
		// compute the position of each box
		x := (i % cols) * box_width
		y := (i / cols) * box_height

		sparkline := widgets.NewSparkline()
		sparkline.MaxVal = 100
		sparkline.LineColor = ui.ColorGreen
		sparkline.TitleStyle.Fg = ui.ColorWhite
		plots[i] = widgets.NewSparklineGroup(sparkline)
		plots[i].SetRect(x, y, x+box_width, y+box_height)
	}

	draw := func() {
		percent, err := c.CPUStatusLoader.GetCPUStatus()
		if err != nil {
			log.Fatalf("failed to get cpu status: %v", err)
		}
		for i := 0; i < count; i++ {
			plots[i].Sparklines[0].Data = percent[i]
			plots[i].Title = fmt.Sprintf("%d - %.2f%%", i, percent[i][0])
			if percent[i][0] > 80 {
				plots[i].Sparklines[0].LineColor = ui.ColorRed
			} else if percent[i][0] > 50 {
				plots[i].Sparklines[0].LineColor = ui.ColorYellow
			} else {
				plots[i].Sparklines[0].LineColor = ui.ColorGreen
			}
		}

		for i := 0; i < count; i++ {
			ui.Render(plots[i])
		}
	}

	defer ui.Clear()

	c.renderLoop(time.Second, draw)
}

func (c *Shukuangkuang) renderMemory() {
	c.MemoryStatusLoader.LaunchGoRoutine()
	defer c.MemoryStatusLoader.StopGoRoutine()
	// top from [0, 5] showes the memory usage of the system
	w, h := ui.TerminalDimensions()
	totalMemoryBarHeight := 10
	if h < totalMemoryBarHeight {
		totalMemoryBarHeight = h
	}
	processBarHeight := h - totalMemoryBarHeight

	sparkline := widgets.NewSparkline()
	sparkline.MaxVal = 100
	sparkline.LineColor = ui.ColorGreen
	sparkline.TitleStyle.Fg = ui.ColorWhite
	totalMemoryBar := widgets.NewSparklineGroup(sparkline)
	totalMemoryBar.Title = "Memory Usage"
	totalMemoryBar.SetRect(0, 0, w, totalMemoryBarHeight)

	drawables := make([]ui.Drawable, 1)
	drawables[0] = totalMemoryBar

	if processBarHeight > 0 {
		processTable := widgets.NewTable()
		processTable.Title = "Process Usage"
		processTable.SetRect(0, totalMemoryBarHeight, w, h)
		processTable.Rows = [][]string{}
		processTable.TextStyle = ui.NewStyle(ui.ColorWhite)
		processTable.RowSeparator = false
		drawables = append(drawables, processTable)
	}

	draw := func() {
		usages, procs, err := c.MemoryStatusLoader.GetMemoryStatus()
		if err != nil {
			log.Fatalf("failed to get memory status: %v", err)
		}

		maxMemory, err := c.MemoryStatusLoader.GetMaxMemory()
		if err != nil {
			log.Fatalf("failed to get max memory: %v", err)
		}

		data := make([]float64, len(usages))
		for i := 0; i < len(usages); i++ {
			data[i] = float64(usages[i]) / float64(maxMemory) * 100
		}

		totalMemoryBar.Sparklines[0].Data = data
		totalMemoryBar.Sparklines[0].Title = fmt.Sprintf("%v/%v", Bytes2Human(usages[0]), Bytes2Human(maxMemory))

		if processBarHeight > 0 {
			processTable := drawables[1].(*widgets.Table)
			processTable.Rows = [][]string{
				{"PID", "Name", "Usage"},
			}
			for i := 0; i < len(procs); i++ {
				processTable.Rows = append(processTable.Rows, []string{fmt.Sprintf("%v", procs[i].Pid), procs[i].Name, fmt.Sprintf("%v", Bytes2Human(procs[i].MemoryUsage))})
			}
		}

		for _, drawable := range drawables {
			ui.Render(drawable)
		}
	}

	defer ui.Clear()

	c.renderLoop(time.Second, draw)
}
