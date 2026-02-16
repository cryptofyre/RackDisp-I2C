package config

type Config struct {
	DeviceType string // "auto", "ssd1306", "emulator"
	Rotate     bool
	Refresh    int // Refresh rate in Hz (or just interval in ms)
}
