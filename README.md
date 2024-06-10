# fwmatrix

[![GoDoc](https://godoc.org/github.com/littlehawk93/fwmatrix?status.svg)](https://godoc.org/github.com/littlehawk93/fwmatrix)

 fwmatrix is a library for drawing pixels on the [Framework 16 LED Matrix module](https://frame.work/products/16-led-matrix). 

 It provides simple libraries to easily allow programmers to control LEDs on the module and even perform animations, while provided low-level access to raw module commands for 100% custom interactions.

## Table of Contents

- [Installation](#Installation)
- [Features](#Features)
    - [Basic Commands](#basic-commands)
    - [Showing Patterns](#showing-patterns)
    - [Matrix Renderers](#matrix-renderers)
    - [Black & White Rendering](#black-&-white-rendering)
    - [Greyscale Rendering](#greyscale-rendering)
    - [Animations](#animations)
- [Feedback](#feedback)

### Installation

To install, simply use the following `go get` command:

```bash
go get github.com/littlehawk93/fwmatrix
```

### Features

fwmatrix provides tools for easily drawing 1-bit black and white images, 8bit greyscale images, displaying pre-programmed patterns, or sending raw commands to the LED matrix module. Additionally, there is a light framework for rendering simple animations smoothly.

#### Basic Commands

To send raw commands to the LED matrix module, you must first connect to the module via a serial port:

```go
port, err := serial.OpenPort(&serial.Config{Name: "COM3", Baud: 115200})
```

Commands can then be sent using the `WriteCommand` function:

```go
if err := fwmatrix.WriteCommand(port, fwmatrix.CmdDrawBW, paramBytes); err != nil {
    log.Fatal(err)
}
```

Some commands do not require additional bytes to be sent as parameters, in that case, simply send `nil`:

```go
if err := fwmatrix.WriteCommand(port, fwmatrix.CmdDrawBW, nil); err != nil {
    log.Fatal(err)
}
```

Visit [this page](https://github.com/FrameworkComputer/inputmodule-rs/blob/main/commands.md) to view documentation on the various commands the LED module supports. Currently supported commands are also defined using the `Command` type.

#### Showing Patterns

The function `ShowPattern` provides an easy shortcut to writing raw pattern commands to the module. Simply send the pattern you want to display. Supported patterns are included as constants of the type `Pattern`:

```go
if err := fwmatrix.ShowPattern(port, fwmatrix.PatZigZag, 0); err != nil {
    log.Fatal(err)
}
```

The third parameter is only used for the `PatPercentage` pattern, which displays a progress bar based on the value from 0 - 100.

The functions: `GetSleepState`, `SetSleepState`, `GetBrightness`, and `SetBrightness` are other shortcuts for sending simple commands for setting / retrieving the sleep or brightness status of the module without sending raw commands.

#### Matrix Renderers

fwmatrix uses matrix renderers for drawing graphics on an LED matrix module. The renderers handle all of the logic for setting LED values and sending various commands to the module in order to draw the graphics, making drawing shapes and various things easy.

There are some basic graphics tools for drawing lines and shapes that work with any `MatrixRenderer`:

```go
fwmatrix.DrawLine(renderer, 0, 0, 34, 9) // draw a line

fwmatrix.DrawRect(renderer, 5, 0, 3, 4) // draw a rectangle

fwmatrix.DrawFillRect(renderer, 5, 0, 3, 4) // draw a filled rectangle
```

#### Black & White Rendering

A matrix renderer for simple 1 bit (Black & White) graphics. Use the global `SetBrightness` function to set uniform brightness for all LEDs, and then draw graphics with this renderer for higher performance (more FPS) graphics, with no brightness control.

Example:

```go
port, err := serial.OpenPort(&serial.Config{Name: "COM3", Baud: 115200})

if err != nil {
    log.Fatal(err)
}

mat := fwmatrix.NewBWWithPort(port)

fwmatrix.SetBrightness(port, 100)

fwmatrix.DrawRect(mat, 0, 0, 5, 5)

if err = mat.Flush(); err != nil {
    log.Fatal(err)
}
```

#### Greyscale Rendering

A matrix renderer that supports brightness control for individual pixels. Use the matrix's `SetBrightness` function to set the current pixel brightness and each subsequent draw function will draw pixels at that brightness.

Example:

```go
mat, err := fwmatrix.NewGS("COM3", 115200)

if err != nil {
    log.Fatal(err)
}

mat.SetBrightness(100)

fwmatrix.DrawRect(mat, 0, 0, 5, 5)

mat.SetBrightness(50)

fwmatrix.DrawLine(mat, 6, 5, 10, 5)

if err = mat.Flush(); err != nil {
    log.Fatal(err)
}
```

#### Animations

fwmatrix also provides a simple animation engine for easily updating graphics at a fixed frame rate. Simply call the `StartAnimation` function with a rendering callback function of type `RenderFrameFunc`. This function will be passed in the current matrix renderer, along with a delta of the time passed (in milliseconds) since the last frame was rendered. If the callback function returns an error value other than `nil`, the animation will stop. Return the error `ErrorStopAnimation` to stop animation without throwing an error. You do not need to call `Clear` or `Flush` inside the `RenderFrameFunc`.

Example:

```go
boxWidth := 3
boxHeight := 3

port, err := serial.OpenPort(&serial.Config{Name: "COM3", Baud: 115200})

if err != nil {
    log.Fatal(err)
}

mat := fwmatrix.NewBWWithPort(port)

fwmatrix.SetBrightness(port, 100)

rng := rand.New(rand.NewSource(time.Now().UnixMicro()))
mat.Clear()

boxX := 3.0
boxY := 0.0

scale := 0.4

boxSpeedX := rng.Float64() * scale / float64(fwmatrix.FrameTime30FPS)
boxSpeedY := rng.Float64() * scale / float64(fwmatrix.FrameTime30FPS) * float64(fwmatrix.MatrixHeight) / float64(fwmatrix.MatrixWidth)

mat.SetBrightness(25)

bounces := 0

err = fwmatrix.StartAnimation(func(m fwmatrix.MatrixRenderer, delta int64) error {

    boxX += boxSpeedX * float64(delta)
    boxY += boxSpeedY * float64(delta)

    if boxX < 0 {
        boxX = 0
        boxSpeedX = -boxSpeedX
        bounces++
    } else if boxX+float64(boxWidth) > float64(fwmatrix.MatrixWidth-1) {
        boxX = float64(fwmatrix.MatrixWidth - boxWidth - 1)
        boxSpeedX = -boxSpeedX
        bounces++
    }

    if boxY < 0 {
        boxY = 0
        boxSpeedY = -boxSpeedY
        bounces++
    } else if boxY+float64(boxHeight) > float64(fwmatrix.MatrixHeight-1) {
        boxY = float64(fwmatrix.MatrixHeight) - float64(boxHeight) - 1
        boxSpeedY = -boxSpeedY
        bounces++
    }

    x := int(math.Round(boxX))
    y := int(math.Round(boxY))

    fwmatrix.DrawFillRect(m, x, y, x+boxWidth, y+boxHeight)
    if bounces >= 20 {
        return fwmatrix.ErrorStopAnimation
    }
    return nil
}, mat, fwmatrix.FrameTime30FPS)

if err != nil {
    log.Fatal(err)
}
```

### Custom Matrix Renderers

You can implement the `MatrixRenderer` interface to have your own custom rendering engine that will work with the various graphical and animation tools in fwmatrix. To do so, implement the following methods:

`Close` - close the underlying serial port used by the matrix renderer
`Clear` - set all LEDs to off in the matrix module
`Flush` - commit and send LED changes to the matrix module using any number of commands via serial
`SetPixel` - set an individual LED in the matrix module on

### Feedback

Thanks for checking out fwmatrix! For any feedback or suggestions on how to improve the library, please create an issue in the GitHub project.