package internal

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func (c *Shukuangkuang) Run(args ShukuangkuangArgs) {
	if c.stop != nil {
		close(c.stop)
	}
	c.stop = make(chan struct{}, 1)

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

	plots := make([]*widgets.SparklineGroup, count)

	for i := 0; i < count; i++ {
		// compute the position of each box
		x := (i % cols) * box_width
		y := (i / cols) * box_height

		sparkline := widgets.NewSparkline()
		sparkline.MaxVal = 100
		sparkline.Title = fmt.Sprintf("CPU %d", i)
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

	// cause PollEvents is not implemented in sync mode, so we need to use goroutine to avoid blocking, it's too slow
	go func() {
		uiEvents := ui.PollEvents()
		for event := range uiEvents {
			switch event.ID {
			case "q", "<C-c>":
				c.stop <- struct{}{}
			}
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			draw()
		case <-c.stop:
			return
		}
	}
}
