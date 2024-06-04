package fwmatrix

type Command uint8

const (
	// CmdBrightness command to set LED brightness
	CmdBrightness Command = 0x00

	// CmdPattern command to display a pre-programmed pattern
	CmdPattern Command = 0x01

	// CmdBootloader command to jump to the bootloader
	CmdBootloader Command = 0x02

	// CmdBootloader command to get or set the sleep state
	CmdSleep Command = 0x03

	// CmdAnimate command to get or set an animation
	CmdAnimate Command = 0x04

	// CmdPanic command to cause a FW panic
	CmdPanic Command = 0x05

	// CmdDrawBW command to draw a 1bit B&W image
	CmdDrawBW Command = 0x06

	// CmdStageCol command to set a greyscale column of pixels
	CmdStageCol Command = 0x07

	// CmdFlushCols command to display greyscale columns
	CmdFlushCols Command = 0x08
)
