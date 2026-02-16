package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"cryptofy.re/rackdisp/internal/config"
	"cryptofy.re/rackdisp/internal/display"
	"cryptofy.re/rackdisp/internal/hardware"
	"cryptofy.re/rackdisp/internal/stats"
	"cryptofy.re/rackdisp/internal/ui"
)

func Run(ctx context.Context, cfg config.Config) error {
	// 1. Hardware Detection
	hwConfig := hardware.GetDeviceConfig()
	log.Printf("Detected hardware: %s. Using I2C Bus: %s\n", hwConfig.Name, hwConfig.I2CBus)

	// 2. Initialize UI
	uiSystem := ui.NewUI()

	// 3. Initialize Display
	var disp display.Display
	var err error

	deviceType := cfg.DeviceType
	if deviceType == "auto" {
		if hwConfig.Name == "Unknown" {
			// If unknown hardware and on Windows (implied by execution environment check if needed, but 'unknown' is good proxy for dev machine)
			// Actually hard to detect dev machine purely from /proc, but let's assume if unknown we default to emulator if configured or fall back to SSD1306
			// For now, let's look at build tags or just simple logic.
			// Ideally detecting if I2C is available.
			deviceType = "emulator" // Default to emulator if hardware not detected? Or maybe "ssd1306" is safer for RPi.
			// Let's stick to explicit default or safe fallback.
			// Revert to legacy behavior: try SSD1306.
			// Will need to look into more device types later on.
			deviceType = "ssd1306"
		} else {
			deviceType = "ssd1306"
		}
	}

	switch deviceType {
	case "emulator":
		disp = display.NewEmulator(4.0)
	case "ssd1306":
		disp, err = display.NewSSD1306(hwConfig.I2CBus, cfg.Rotate)
		if err != nil {
			return fmt.Errorf("failed to init ssd1306: %v", err)
		}
	default:
		return fmt.Errorf("unknown device type: %s", deviceType)
	}
	defer disp.Close()

	// 4. Start Stats Collector (Async)
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.Refresh) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s, err := stats.GetSystemStats()
				if err != nil {
					log.Printf("Error getting stats: %v", err)
					continue
				}
				// Populate temp label from hw config if needed, or just raw temp
				// s.Temp already has "XXC"
				uiSystem.UpdateStats(s)
			}
		}
	}()

	// 5. Run Display Loop (Blocking)
	log.Printf("Starting RackDisp with device type: %s\n", deviceType)

	inputHandler := func(event string) {
		switch event {
		case "next_mode":
			uiSystem.NextMode()
		case "toggle_net":
			uiSystem.ToggleMockNet()
		case "toggle_workload":
			uiSystem.ToggleMockWorkload()
		}
	}

	return disp.Run(ctx, uiSystem.GetFrame, inputHandler)
}
