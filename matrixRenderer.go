package fwmatrix

// MatrixRenderer renders pixels on the Framework 16 LED Matrix Module
type MatrixRenderer interface {
	Close() error
	Clear()
	Flush() error
	SetPixel(x, y int)
}
