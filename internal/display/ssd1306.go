package display

import (
	"context"
	"fmt"
	"image"
	"log"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
)

type SSD1306Display struct {
	dev    *ssd1306.Dev
	bus    i2c.BusCloser
	stop   chan struct{}
	rotate bool
}

// Ensure interface implementation
var _ Display = (*SSD1306Display)(nil)

func NewSSD1306(i2cBus string, rotate bool) (*SSD1306Display, error) {
	// Initialize periph.io host if not already done (it's safe to call multiple times internally usually, but good to ensure)
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize periph: %v", err)
	}

	b, err := i2creg.Open(i2cBus)
	if err != nil {
		return nil, fmt.Errorf("failed to open I2C bus %s: %v", i2cBus, err)
	}

	opts := ssd1306.DefaultOpts
	opts.Rotated = rotate

	dev, err := ssd1306.NewI2C(b, &opts)
	if err != nil {
		b.Close()
		return nil, fmt.Errorf("failed to initialize ssd1306: %v", err)
	}

	return &SSD1306Display{
		dev:    dev,
		bus:    b,
		stop:   make(chan struct{}),
		rotate: rotate,
	}, nil
}

func (d *SSD1306Display) Run(ctx context.Context, renderer Renderer, inputHandler InputHandler) error {
	// 30Hz refresh for smooth animations locally, even if stats are slower
	ticker := time.NewTicker(33 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-d.stop:
			return nil
		case <-ticker.C:
			img, err := renderer()
			if err != nil {
				log.Printf("Renderer error: %v", err)
				continue
			}
			if err := d.dev.Draw(img.Bounds(), img, image.Point{}); err != nil {
				log.Printf("Draw error: %v", err)
			}
		}
	}
}

func (d *SSD1306Display) Close() error {
	close(d.stop)
	return d.bus.Close()
}
