package internal

import (
	"fmt"
	"log"
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

	c.CPUStatusLoader.Init()

	c.renderCpu()
}

func (c *Shukuangkuang) renderCpu() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
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

	plots := make([]*widgets.Plot, count)

	for i := 0; i < count; i++ {
		// compute the position of each box
		x := (i % cols) * box_width
		y := (i / cols) * box_height

		plots[i] = widgets.NewPlot()
		plots[i].MaxVal = 100
		plots[i].Title = fmt.Sprintf("CPU %d", i)
		plots[i].SetRect(x, y, x+box_width, y+box_height)
		plots[i].AxesColor = ui.ColorWhite
		plots[i].LineColors = []ui.Color{ui.ColorGreen}
		plots[i].TitleStyle.Fg = ui.ColorWhite
	}

	draw := func() {
		percent, err := c.CPUStatusLoader.GetCPUStatus()
		if err != nil {
			log.Fatalf("failed to get cpu status: %v", err)
		}
		for i := 0; i < count; i++ {
			plots[i].Data = [][]float64{percent[i]}
			plots[i].Title = fmt.Sprintf("%d - %.2f%%", i, percent[i][0])
		}

		for i := 0; i < count; i++ {
			ui.Render(plots[i])
		}
	}

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			draw()
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
	}
}
