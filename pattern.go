package fwmatrix

type Pattern uint8

const (
	// PatPercentage displays a progress indicator using a provided percetange
	PatPercentage Pattern = 0x00

	// PatGradient displays a brightness gradient from top to bottom
	PatGradient Pattern = 0x01

	// PatDoubleGradient displays a brightness gradient from the middle to both top and bottom
	PatDoubleGradient Pattern = 0x02

	// PatLotusHor displays the text "LOTUS" horizontally across the matrix
	PatLotusHor Pattern = 0x03

	// PatZigZag displays a zigzag pattern
	PatZigZag Pattern = 0x04

	// PatFullBrightness displays all LEDs at 100% brightness
	PatFullBrightness Pattern = 0x05

	// PatPanic displays the text "PANIC" across the matrix
	PatPanic Pattern = 0x06

	// PatLotusVert displays the text "LOTUS" vertically across the matrix
	PatLotusVert Pattern = 0x07
)
