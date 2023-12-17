package internal

import (
	"fmt"
	"math"
)

func CaculateBoxSizes(N, w, h int) (int, int) {
	h = h * 3

	k := int(math.Sqrt(float64(N)))

	rows := k * h / w
	cols := k * w / h

	width := 0

	if rows >= k {
		width = w / k
	} else {
		width = w / cols
	}

	rows = int(math.Ceil(float64(N) / float64(cols)))
	height := h / rows

	return width, height / 3
}

func Bytes2Human(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2fKB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(bytes)/1024/1024)
	} else {
		return fmt.Sprintf("%.2fGB", float64(bytes)/1024/1024/1024)
	}
}
