package internal

import "math"

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
