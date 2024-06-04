package fwmatrix

import (
	"math"

	"github.com/tarm/serial"
)

const (
	bwMatrixBufferSize int = 39
)

// BWMatrix a tool for drawing basic 1bit pixel data on a LED Matrix module
type BWMatrix struct {
	port       *serial.Port
	drawBuffer []uint8
}

// Close closes the underlying serial port used by this BWMatrix
func (me *BWMatrix) Close() error {

	return me.port.Close()
}

// Clear sets all pixels to off
func (me *BWMatrix) Clear() {

	for i := range me.drawBuffer {
		me.drawBuffer[i] = 0
	}
}

// DrawPixel turns a pixel at the provided coordinate on
func (me *BWMatrix) DrawPixel(x, y int) {

	if x < 0 || x >= MatrixWidth || y < 0 || y >= MatrixHeight {
		return
	}

	index := y*MatrixWidth + x

	me.drawBuffer[index/8] |= (1 << (index % 8))
}

// DrawLine draws a line from a point to another
func (me *BWMatrix) DrawLine(x0, y0, x1, y1 int) {

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
		me.DrawPixel(int(math.Round(x)), int(math.Round(y)))
		x += xi
		y += yi
	}
}

// DrawRect draws an empty rectangle with opposite corners at x0,y0 and x1,y1
func (me *BWMatrix) DrawRect(x0, y0, x1, y1 int) {

	me.DrawLine(x0, y0, x1, y0)
	me.DrawLine(x0, y0, x0, y1)
	me.DrawLine(x1, y0, x1, y1)
	me.DrawLine(x0, y1, x1, y1)
}

// DrawFillRect draws a filled rectangle with opposite corners at x0,y0 and x1,y1
func (me *BWMatrix) DrawFillRect(x0, y0, x1, y1 int) {
	for x := x0; x <= x1; x++ {
		me.DrawLine(x, y0, x, y1)
	}
}

// Flush writes the current stored pixel buffer to the LED Matrix module to display.
// Returns any errors encountered during serial communications
func (me *BWMatrix) Flush() error {

	return WriteCommand(me.port, CmdDrawBW, me.drawBuffer)
}

// NewBW opens a serial port with the provided name and baud rate and initializes a new BWMatrix with it.
// Returns any errors that occurred when opening the serial port
func NewBW(names string, baud int) (*BWMatrix, error) {

	p, err := serial.OpenPort(&serial.Config{Name: names, Baud: baud})

	if err != nil {
		return nil, err
	}
	return &BWMatrix{port: p, drawBuffer: make([]uint8, bwMatrixBufferSize)}, nil
}

// NewBWWithPort create and initialize a new BWMatrix using the provided serial port
func NewBWWithPort(port *serial.Port) *BWMatrix {

	return &BWMatrix{port: port, drawBuffer: make([]uint8, bwMatrixBufferSize)}
}
