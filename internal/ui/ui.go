package ui

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"cryptofy.re/rackdisp/internal/stats"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// UI manages the state and rendering of the display.
// It handles the "Business Logic" of what gets drawn to the screen.
type UI struct {
	mu           sync.RWMutex
	currentStats *stats.SystemStats
	lastUpdate   time.Time
	ticks        uint64
	width        int
	height       int

	// Mode state controls which "Page" is currently visible.
	// 0: Overview, 1: Network, 2: System, 3: AI/Ollama
	mode       int
	modeTicker int // Logic tick counter to handle auto-rotation

	// Simulation flags for testing UI elements without real hardware/workloads
	mockNet      bool
	mockWorkload bool
	useMocks     bool
}

func NewUI() *UI {
	return &UI{
		currentStats: &stats.SystemStats{
			IP:       "Init...",
			Hostname: "Loading...",
			Temp:     "--",
		},
		width:  128,
		height: 64,
		mode:   0,
	}
}

// UpdateStats is called by the main loop to inject fresh system data.
func (u *UI) UpdateStats(s *stats.SystemStats) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.currentStats = s
	u.lastUpdate = time.Now()
}

// --- Input Methods (for Emulator/Interactive usage) ---

func (u *UI) NextMode() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.mode = (u.mode + 1) % 4
	u.modeTicker = 0
}

func (u *UI) ToggleMockNet() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.useMocks = true
	u.mockNet = !u.mockNet
}

func (u *UI) ToggleMockWorkload() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.useMocks = true
	u.mockWorkload = !u.mockWorkload
}

// GetFrame generates the actual image to be displayed on the OLED.
// It is called ~30 times a second.
func (u *UI) GetFrame() (image.Image, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.ticks++
	u.modeTicker++

	// Page Cycle Timer
	// We want to switch pages every 20 seconds.
	// Assuming ~30fps -> 600 ticks.
	cycleTicks := 600
	if u.modeTicker > cycleTicks {
		u.mode = (u.mode + 1) % 4
		u.modeTicker = 0
	}

	img := image.NewGray(image.Rect(0, 0, u.width, u.height))

	// --- STATUS BAR (Top 16px - Yellow Zone) ---

	// 1. Hostname on the Left
	hostStr := u.currentStats.Hostname
	if len(hostStr) > 10 {
		hostStr = hostStr[:10]
	}
	drawText(img, 2, 10, hostStr)

	// 2. Tab Indicators (Center-Right)
	// These show which of the 4 pages is active.
	u.drawTabs(img, 74, 8, 4, u.mode)

	// 3. Status Icons (Far Right)
	// Network Icon
	netStatus := u.currentStats.IP != "Offline" && u.currentStats.IP != ""
	if u.useMocks {
		netStatus = u.mockNet
	}
	u.drawNetIcon(img, 108, 2, netStatus)

	// Workload Icon (Hammer)
	// Activates if CPU > 50% OR if we detect an AI model running.
	workloadActive := u.currentStats.CPUPercent > 50.0
	if len(u.currentStats.OllamaModels) > 0 {
		workloadActive = true
	}

	// Handle Mock Data Injection for Emulation
	if u.useMocks {
		if u.mockWorkload {
			workloadActive = true
			// If we are forcing workload in the emulator, inject a fake model
			// so the user can see the detailed view on Page 4.
			if len(u.currentStats.OllamaModels) == 0 {
				u.currentStats.OllamaModels = []stats.ModelInfo{
					{Name: "llama3:latest", ParamSize: "8B", Quantization: "Q4_0", VRAM: 5046586572},
				}
			}
		} else {
			// Clear the fake model if we toggle it off
			if len(u.currentStats.OllamaModels) == 1 && u.currentStats.OllamaModels[0].Name == "llama3:latest" {
				u.currentStats.OllamaModels = []stats.ModelInfo{}
			}
		}
	}
	if workloadActive {
		u.drawHammerIcon(img, 119, 2)
	}

	// 4. Progress Bar
	// A thin line at the bottom of the header showing time until next page switch.
	progress := float64(u.modeTicker) / float64(cycleTicks)
	u.drawProgressBar(img, 0, 14, u.width, 2, progress)

	// --- CONTENT AREA (Bottom 48px - Blue Zone, y=16+) ---
	offsetY := 16

	switch u.mode {
	case 0: // Overview Mode (CPU/RAM/Temp)
		drawText(img, 2, offsetY+10, "CPU")
		u.drawBar(img, 28, offsetY+4, 80, 6, u.currentStats.CPUPercent)

		drawText(img, 2, offsetY+24, "RAM")
		u.drawBar(img, 28, offsetY+18, 80, 6, u.currentStats.MemPercent)

		// Temp bottom center
		tempStr := fmt.Sprintf("Temp: %s", u.currentStats.Temp)
		drawText(img, 28, offsetY+40, tempStr)

	case 1: // Network Mode (Big IP)
		drawText(img, 2, offsetY+12, "IP Address:")
		drawText(img, 2, offsetY+28, u.currentStats.IP)

	case 2: // System Mode (Disk/Uptime)
		drawText(img, 2, offsetY+10, "Disk")
		u.drawBar(img, 34, offsetY+4, 74, 6, u.currentStats.DiskPercent)

		drawText(img, 2, offsetY+28, "Uptime:")
		drawText(img, 2, offsetY+42, u.currentStats.Uptime)

	case 3: // AI/Ollama Mode
		drawText(img, 2, offsetY+10, "AI Workload")
		if len(u.currentStats.OllamaModels) > 0 {
			// Show first model info
			m := u.currentStats.OllamaModels[0]
			modelName := m.Name
			drawText(img, 2, offsetY+20, modelName)

			// Details: "8B / Q4_0"
			details := fmt.Sprintf("%s / %s", m.ParamSize, m.Quantization)
			drawText(img, 2, offsetY+32, details)

			// VRAM: "4.7GB VRAM"
			vramGB := float64(m.VRAM) / 1024 / 1024 / 1024
			vramStr := fmt.Sprintf("%.1fGB VRAM", vramGB)
			drawText(img, 2, offsetY+44, vramStr)
		} else {
			drawText(img, 2, offsetY+24, "Idle")
		}
	}

	// Global Animation (Heartbeat dot in corner to show the app hasn't frozen)
	if (u.ticks/15)%2 == 0 {
		img.SetGray(124, 60, color.Gray{Y: 255})
	}

	return img, nil
}

// Helpers

func (u *UI) drawProgressBar(img *image.Gray, x, y, w, h int, progress float64) {
	if progress > 1.0 {
		progress = 1.0
	}
	fillW := int(float64(w) * progress)
	c := color.Gray{Y: 255} // White
	for i := 0; i < fillW; i++ {
		for j := 0; j < h; j++ {
			img.SetGray(x+i, y+j, c)
		}
	}
}

func (u *UI) drawTabs(img *image.Gray, x, y, count, active int) {
	spacing := 8
	for i := 0; i < count; i++ {
		cx := x + (i * spacing)
		c := color.Gray{Y: 255} // White
		if i == active {
			// Active Page: Filled 3x3 box
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					img.SetGray(cx+dx, y+dy, c)
				}
			}
		} else {
			// Inactive Page: Outline 3x3 box
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == -1 || dx == 1 || dy == -1 || dy == 1 {
						img.SetGray(cx+dx, y+dy, c)
					}
				}
			}
		}
	}
}

func (u *UI) drawNetIcon(img *image.Gray, x, y int, active bool) {
	if !active {
		return // Do not draw anything if offline
	}
	c := color.Gray{Y: 255} // White

	// Icon: Globe with Cursor (Pixel Art)
	// Center around x+5, y+4
	cx, cy := x+5, y+4

	// 1. Draw Globe (Grid Pattern)
	for i := -4; i <= 4; i++ {
		for j := -4; j <= 4; j++ {
			dist := i*i + j*j
			// Outline (Radius ~4)
			if dist >= 12 && dist <= 20 {
				img.SetGray(cx+i, cy+j, c)
			}
			// Inner Grid Lines (Equator/Meridian)
			if dist < 16 {
				if i == 0 || j == 0 {
					img.SetGray(cx+i, cy+j, c)
				}
				// Optional: Latitudes?
				if j == -2 || j == 2 {
					img.SetGray(cx+i, cy+j, c)
				}
			}
		}
	}

	// 2. Draw Cursor Overlay (Bottom Right)
	cursorX, cursorY := x+6, y+5
	cursorColor := c

	points := []struct{ dx, dy int }{
		{0, 0}, {0, 1}, {0, 2}, {0, 3},
		{1, 1}, {2, 2},
		{1, 3}, // tail
	}
	for _, p := range points {
		img.SetGray(cursorX+p.dx, cursorY+p.dy, cursorColor)
	}
}

func (u *UI) drawHammerIcon(img *image.Gray, x, y int) {
	// Icon: Sledgehammer (Vertical)
	c := color.Gray{Y: 255} // White

	// 1. Handle
	for j := 4; j < 9; j++ {
		img.SetGray(x+4, y+j, c)
	}

	// 2. Head (Heavy Block)
	for i := 1; i < 8; i++ {
		for j := 0; j < 4; j++ {
			img.SetGray(x+i, y+j, c)
		}
	}

	// 3. Animation: "Strike" effect (Spark)
	if (u.ticks/15)%2 == 0 {
		img.SetGray(x+8, y+8, c)
		img.SetGray(x+9, y+7, c)
	}
}

func (u *UI) drawBar(img *image.Gray, x, y, w, h int, percent float64) {
	borderColor := color.Gray{Y: 255}
	fillColor := color.Gray{Y: 255}

	// Clamp values
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	// Draw Border
	for i := x; i < x+w; i++ {
		img.SetGray(i, y, borderColor)
		img.SetGray(i, y+h-1, borderColor)
	}
	for j := y; j < y+h; j++ {
		img.SetGray(x, j, borderColor)
		img.SetGray(x+w-1, j, borderColor)
	}

	// Fill based on percentage
	fillW := int((percent / 100.0) * float64(w-2))
	for i := x + 1; i < x+1+fillW; i++ {
		for j := y + 1; j < y+h-1; j++ {
			img.SetGray(i, j, fillColor)
		}
	}
}

func drawText(img *image.Gray, x, y int, text string) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(text)
}
