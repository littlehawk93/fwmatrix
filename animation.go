package fwmatrix

import (
	"errors"
	"time"
)

const (
	// FrameTime60FPS frame time for 60 fps animation
	FrameTime60FPS int64 = 17

	// FrameTime30FPS frame time for 30 fps animation
	FrameTime30FPS int64 = 33

	// FrameTime24FPS frame time for 24 fps animation
	FrameTime24FPS int64 = 42
)

// ErrorStopAnimation error used to signal to the
var ErrorStopAnimation error = errors.New("fwmatrix - stop animation")

// RenderFrameFunc function for rendering a single frame of an animation.
// Delta is the time (in milliseconds) since the last frame render call was made
// Should return any errors encountered during rendering.
// Do not call r.Flush() inside the this function.
// Return ErrorStopAnimation to halt animation without throwing an error
type RenderFrameFunc func(m MatrixRenderer, delta int64) error

// StartAnimation begins an animation on the provided MatrixRenderer using the provided rendering function.
// Frame time defines how many milliseconds should pass between rendering each frame.
// Stops whenever the returned error of the render function is not null.
// If the returned error is ErrorStopAnimation, this function will return nil. Otherwise, the error is returned
func StartAnimation(render RenderFrameFunc, m MatrixRenderer, frameTime int64) error {

	lastFrame := time.Now().UnixMilli()

	for {
		currTime := time.Now().UnixMilli()
		frameDiff := lastFrame - currTime
		lastFrame = currTime
		m.Clear()
		if err := render(m, frameDiff); err != nil {
			if !errors.Is(err, ErrorStopAnimation) {
				return err
			} else {
				break
			}
		}

		if err := m.Flush(); err != nil {
			return err
		}

		renderTime := time.Now().UnixMilli() - currTime

		if renderTime < frameTime {
			time.Sleep(time.Millisecond * time.Duration(frameTime-renderTime))
		}
	}
	return nil
}
