package fwmatrix

import "github.com/tarm/serial"

const (
	magicByteH uint8 = 0x32
	magicByteL uint8 = 0xAC
)

// ShowPattern displays a pre-programmed pattern on a LED matrix module via the provided serial port. The PatPercentage pattern requires the per parameter to set the percentage to display.
// Returns any errors encountered during serial writing
func ShowPattern(p *serial.Port, pat Pattern, per uint8) error {

	buf := make([]uint8, 1)

	if pat == PatPercentage {
		buf[0] = per
	} else {
		buf = nil
	}
	return WriteCommand(p, CmdPattern, buf)
}

// WriteCommand send a command to a LED matrix module via the provided serial port. Set params to nil if no parameters are needed.
// Returns any errors encountered during serial writing
func WriteCommand(p *serial.Port, c Command, params []byte) error {

	buf := []uint8{magicByteH, magicByteL, uint8(c)}

	if len(params) > 0 {
		buf = append(buf, params...)
	}

	if _, err := p.Write(buf); err != nil {
		return err
	}
	return nil
}
