package display

import (
	"context"
	"image"
)

// Renderer is a function that generates the current frame image.
type Renderer func() (image.Image, error)

// InputHandler is a function that processes input events (e.g. key presses)
type InputHandler func(event string)

// Display is the interface for different display backends (Hardware vs Emulator)
type Display interface {
	// Run starts the display loop. It blocks until context is cancelled or error occurs.
	// It is responsible for calling the Renderer at the appropriate refresh rate.
	// InputHandler is optional and may be nil if the display doesn't support input.
	Run(ctx context.Context, renderer Renderer, inputHandler InputHandler) error
	Close() error
}
