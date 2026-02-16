package stats

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemStats struct {
	IP           string
	Hostname     string
	CPUPercent   float64
	MemPercent   float64
	DiskPercent  float64
	Uptime       string
	Temp         string
	OllamaModels []ModelInfo
}

func GetSystemStats() (*SystemStats, error) {
	ip := GetLocalIP()
	hostName, _ := os.Hostname()
	cpuPercent, _ := cpu.Percent(0, false)
	memStat, _ := mem.VirtualMemory()

	// Disk Usage (Root)
	diskStat, _ := disk.Usage("/")
	diskVal := 0.0
	if diskStat != nil {
		diskVal = diskStat.UsedPercent
	}

	// Uptime
	uptimeStr := "N/A"
	hostStat, err := host.Info()
	if err == nil {
		uptime := time.Duration(hostStat.Uptime) * time.Second
		days := int(uptime.Hours()) / 24
		hours := int(uptime.Hours()) % 24
		uptimeStr = fmt.Sprintf("%dd %dh", days, hours)
	}

	// Ollama (optional)
	ollamaModels, _ := GetRunningOllamaModels()

	sysTemp := GetTemp()

	cpuVal := 0.0
	if len(cpuPercent) > 0 {
		cpuVal = cpuPercent[0]
	}

	return &SystemStats{
		IP:           ip,
		Hostname:     hostName,
		CPUPercent:   cpuVal,
		MemPercent:   memStat.UsedPercent,
		DiskPercent:  diskVal,
		Uptime:       uptimeStr,
		Temp:         sysTemp,
		OllamaModels: ollamaModels,
	}, nil
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "Offline"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func GetTemp() string {
	// thermal_zone0 handles the primary temperature on both Jetson and Pi architectures
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return "N/A"
	}
	tempStr := strings.TrimSpace(string(data))
	if tempInt, err := strconv.Atoi(tempStr); err == nil {
		return fmt.Sprintf("%dCR", tempInt/1000)
	}
	return "N/A"
}
