package fwmatrix

import (
	"github.com/tarm/serial"
)

const (
	gsMatrixBufferSize int = MatrixHeight * MatrixWidth
)

// GSMatrix a tool for drawing basic 8bit greyscale pixel data on a LED Matrix module
type GSMatrix struct {
	port       *serial.Port
	drawBuffer []uint8
	brightness uint8
}

// GetBrightness get the current pixel rendering brightness
func (me *GSMatrix) GetBrightness() uint8 {
	return me.brightness
}

// SetBrightness set the current pixel rendering brightness
func (me *GSMatrix) SetBrightness(brightness uint8) {
	me.brightness = brightness
}

// Close closes the underlying serial port used by this BWMatrix
func (me *GSMatrix) Close() error {

	return me.port.Close()
}

// Clear sets all pixels to zero
func (me *GSMatrix) Clear() {

	for i := range me.drawBuffer {
		me.drawBuffer[i] = 0
	}
}

// Flush writes the current stored pixel buffer to the LED Matrix module to display.
// Returns any errors encountered during serial communications
func (me *GSMatrix) Flush() error {

	for col := 0; col < MatrixWidth; col++ {
		if err := WriteCommand(me.port, CmdStageCol, me.getColumnCommand(col)); err != nil {
			return err
		}
	}
	return WriteCommand(me.port, CmdFlushCols, nil)
}

// SetPixel sets the pixel at the provided coordinate to the current brightness
func (me *GSMatrix) SetPixel(x, y int) {

	if x < 0 || x >= MatrixWidth || y < 0 || y >= MatrixHeight {
		return
	}

	index := x*MatrixHeight + y

	me.drawBuffer[index] = me.brightness
}

func (me *GSMatrix) getColumnCommand(col int) []uint8 {

	buf := []uint8{uint8(col)}

	return append(buf, me.drawBuffer[col*MatrixHeight:(col+1)*MatrixHeight]...)
}

// NewGS opens a serial port with the provided name and baud rate and initializes a new GSMatrix with it.
// Returns any errors that occurred when opening the serial port
func NewGS(names string, baud int) (*GSMatrix, error) {

	p, err := serial.OpenPort(&serial.Config{Name: names, Baud: baud})

	if err != nil {
		return nil, err
	}
	return &GSMatrix{port: p, drawBuffer: make([]uint8, gsMatrixBufferSize)}, nil
}

// NewGSWithPort create and initialize a new GSMatrix using the provided serial port
func NewGSWithPort(port *serial.Port) *GSMatrix {

	return &GSMatrix{port: port, drawBuffer: make([]uint8, gsMatrixBufferSize)}
}
