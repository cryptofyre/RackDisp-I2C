//go:build headless

package display

import (
	"context"
	"fmt"
)

type EmulatorDisplay struct {
}

// Ensure interface implementation
var _ Display = (*EmulatorDisplay)(nil)

func NewEmulator(scale float64) *EmulatorDisplay {
	return &EmulatorDisplay{}
}

func (g *EmulatorDisplay) Run(ctx context.Context, renderer Renderer, inputHandler InputHandler) error {
	return fmt.Errorf("emulator support not compiled in this build")
}

func (g *EmulatorDisplay) Close() error {
	return nil
}
