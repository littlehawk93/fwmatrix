package fwmatrix

import "math"

// DrawLine draws a line with the provided matrix renderer from a point to another
func DrawLine(r MatrixRenderer, x0, y0, x1, y1 int) {

	dx := x1 - x0
	dy := y1 - y0

	mx := dx

	if mx < 0 {
		mx = -mx
	}

	my := dy

	if my < 0 {
		my = -my
	}

	steps := mx

	if my > mx {
		steps = my
	}

	xi := float64(dx) / float64(steps)
	yi := float64(dy) / float64(steps)

	x := float64(x0)
	y := float64(y0)

	for i := 0; i <= steps; i++ {
		r.SetPixel(int(math.Round(x)), int(math.Round(y)))
		x += xi
		y += yi
	}
}

// DrawRect draws an empty rectangle with the provided matrix with opposite corners at x0,y0 and x1,y1
func DrawRect(r MatrixRenderer, x0, y0, x1, y1 int) {

	DrawLine(r, x0, y0, x1, y0)
	DrawLine(r, x0, y0, x0, y1)
	DrawLine(r, x1, y0, x1, y1)
	DrawLine(r, x0, y1, x1, y1)
}

// DrawFillRect draws a filled rectangle with the provided matrix with opposite corners at x0,y0 and x1,y1
func DrawFillRect(r MatrixRenderer, x0, y0, x1, y1 int) {

	dx := x1 - x0
	dy := y1 - y0

	if dx < 0 {
		dx = -dx
	}

	if dy < 0 {
		dy = -dy
	}

	if dy > dx {
		for x := x0; x <= x1; x++ {
			DrawLine(r, x, y0, x, y1)
		}
	} else {
		for y := y0; y <= y1; y++ {
			DrawLine(r, x0, y, x1, y)
		}
	}
}
