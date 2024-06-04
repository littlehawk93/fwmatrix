package fwmatrix

import "github.com/tarm/serial"

const (
	magicByteH uint8 = 0x32
	magicByteL uint8 = 0xAC
)

// ShowPattern displays a pre-programmed pattern on a LED matrix module via the provided serial port.
// The PatPercentage pattern requires the per parameter to set the percentage to display.
// Returns any errors encountered during serial communications
func ShowPattern(p *serial.Port, pat Pattern, per uint8) error {

	buf := []uint8{uint8(pat)}

	if pat == PatPercentage {

		if per > 100 {
			per = 100
		}
		buf = append(buf, per)
	}

	return WriteCommand(p, CmdPattern, buf)
}

// WriteCommand send a command to a LED matrix module via the provided serial port.
// Set params to nil if no parameters are needed.
// Returns any errors encountered during serial communications
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

// SetBrightness sets the global brightness for all pixels on the LED Matrix module
// Returns any errors encountered during serial communications
func SetBrightness(p *serial.Port, brightness uint8) error {

	return WriteCommand(p, CmdSleep, []byte{brightness})
}

// Panic cause a FW panic.
// Returns any errors encountered during serial communications
func Panic(p *serial.Port) error {

	return WriteCommand(p, CmdPanic, nil)
}

// SetSleepState if sleep is true, sets the LED Matrix module to sleep until turned off (set to false)
// Returns any errors encountered during serial communications
func SetSleepState(p *serial.Port, sleep bool) error {

	state := uint8(0)

	if sleep {
		state = 1
	}
	return WriteCommand(p, CmdSleep, []byte{state})
}

// GetSleepState returns whether or not the LED matrix module is in sleep mode.
// Returns true if in sleep mode, false otherwise.
// Returns any errors encountered during serial communications
func GetSleepState(p *serial.Port) (bool, error) {

	return getCommandBoolResult(p, CmdSleep)
}

// SetAnimation if enabled, scrolls through various pre-programmed animations until disabled.
// Returns any errors encountered during serial communications
func SetAnimationEnabled(p *serial.Port, enable bool) error {

	state := uint8(0)

	if enable {
		state = 1
	}
	return WriteCommand(p, CmdAnimate, []byte{state})
}

// GetAnimateState returns whether or not the LED matrix module is in the middle of an animation
// Returns true if in an animation, false otherwise.
// Returns any errors encountered during serial communications
func GetAnimateState(p *serial.Port) (bool, error) {

	return getCommandBoolResult(p, CmdAnimate)
}

func getCommandBoolResult(p *serial.Port, c Command) (bool, error) {
	if err := WriteCommand(p, c, nil); err != nil {
		return false, err
	}

	buf := make([]byte, 1)

	if err := readAll(p, buf); err != nil {
		return false, err
	}

	return buf[0] != 0, nil
}

func readAll(p *serial.Port, buf []byte) error {
	bytesRead := 0

	for bytesRead < len(buf) {
		cnt, err := p.Read(buf)

		bytesRead += cnt
		if err != nil {
			return err
		}
	}
	return nil
}
