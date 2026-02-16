//go:build !headless

package display

import (
	"context"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EmulatorDisplay struct {
	renderer     Renderer
	inputHandler InputHandler
	ctx          context.Context
	width        int
	height       int
	scale        float64
}

// Ensure interface implementation
var _ Display = (*EmulatorDisplay)(nil)

func NewEmulator(scale float64) *EmulatorDisplay {
	if scale <= 0 {
		scale = 4.0 // Default scale 4x (128x64 -> 512x256)
	}
	return &EmulatorDisplay{
		width:  128,
		height: 64,
		scale:  scale,
	}
}

func (g *EmulatorDisplay) Update() error {
	if g.ctx != nil {
		select {
		case <-g.ctx.Done():
			return ebiten.Termination
		default:
		}
	}

	// Handle Input
	if g.inputHandler != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.inputHandler("next_mode")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.inputHandler("toggle_net")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.inputHandler("toggle_workload")
		}
	}

	return nil
}

func (g *EmulatorDisplay) Draw(screen *ebiten.Image) {
	if g.renderer == nil {
		return
	}
	img, err := g.renderer()
	if err != nil {
		log.Printf("Renderer error: %v", err)
		return
	}

	// Draw the image onto the screen
	op := &ebiten.DrawImageOptions{}
	ebImg := ebiten.NewImageFromImage(img)
	screen.DrawImage(ebImg, op)
}

func (g *EmulatorDisplay) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func (d *EmulatorDisplay) Run(ctx context.Context, renderer Renderer, inputHandler InputHandler) error {
	d.ctx = ctx
	d.renderer = renderer
	d.inputHandler = inputHandler

	ebiten.SetWindowSize(int(float64(d.width)*d.scale), int(float64(d.height)*d.scale))
	ebiten.SetWindowTitle("RackDisp-I2C Emulator")

	if err := ebiten.RunGame(d); err != nil && err != ebiten.Termination {
		return err
	}
	return nil
}

func (d *EmulatorDisplay) Close() error {
	return nil
}
