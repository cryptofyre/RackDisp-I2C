package hardware

import (
	"os"
	"strings"
)

// DeviceConfig holds the specific variables for the detected hardware
type DeviceConfig struct {
	Name    string
	I2CBus  string
	TempLbl string
}

// GetDeviceConfig reads the Linux device tree to determine the board
func GetDeviceConfig() DeviceConfig {
	data, err := os.ReadFile("/proc/device-tree/model")
	if err == nil {
		model := strings.ToLower(string(data))
		if strings.Contains(model, "jetson") {
			return DeviceConfig{Name: "Jetson Orin Nano", I2CBus: "/dev/i2c-7", TempLbl: "SoC Temp:"}
		} else if strings.Contains(model, "raspberry pi 5") {
			return DeviceConfig{Name: "Raspberry Pi 5", I2CBus: "/dev/i2c-1", TempLbl: "CPU Temp:"}
		}
	}
	// Fallback to Pi default if detection fails
	return DeviceConfig{Name: "Unknown", I2CBus: "/dev/i2c-1", TempLbl: "Temp:"}
}
