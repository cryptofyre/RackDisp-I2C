# RackDisp-I2C 🖥️

> **A high-performance, animated system monitor for 1U Racks, Raspberry Pis, and Jetson Nanos.**

RackDisp-I2C is a robust CLI application designed to display critical system metrics on SSD1306 OLED displays. Optimized for I2C communication, it features smooth 30FPS animations, detailed AI workload monitoring (via Ollama), and a rotating status page system.

It is built with [periph.io](https://periph.io/) for hardware abstraction and [gopsutil](https://github.com/shirou/gopsutil) for system stats.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21%2B-cyan)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20ARM-orange)

<img width="320" height="160" alt="rackdisp_5scxi4cBWJ" src="https://github.com/user-attachments/assets/f204cf1d-7fcb-48bc-bdbc-52f5b1600063" />
<img width="320" height="160" alt="rackdisp_nuXsmpE0ln" src="https://github.com/user-attachments/assets/c6c09b99-b1ee-4fd9-88f6-f3d0e093ba9c" />
<img width="320" height="160" alt="rackdisp_Nrvj5YBp4W" src="https://github.com/user-attachments/assets/347d8197-1aa7-4d26-8bc4-d4bc018e9d95" />
<img width="320" height="160" alt="chrome_XJCasBxtWB" src="https://github.com/user-attachments/assets/24bb2cd5-5bad-495e-8636-b3602f0ed6a1" />

## ✨ Features

- **📊 Multi-Page Dashboard**: Automatically cycles every 20 seconds between:
    1.  **Overview**: CPU, RAM, and Temperature.
    2.  **Network**: Large, readable IP Address.
    3.  **Storage & Uptime**: Disk usage (Root) and system uptime.
    4.  **AI Workload**: Real-time Ollama model tracking (Name, Quantization, VRAM usage).
- **🎨 High-Contrast UI**: Optimized for monochrome OLEDs (1-bit color).
- **🔌 Hardware Support**: 
    - Raspberry Pi 4/5
    - NVIDIA Jetson Nano / Orin Nano
    - Any Linux device with I2C exposed (`/dev/i2c-*`).
- **🪟 Desktop Emulator**: Test the UI on Windows/Mac/Linux without hardware.
- **🛠️ Production Ready**: Includes a built-in installer for systemd service generation.

## 🚀 Installation

### 1. Pre-requisites (Linux)
Ensure I2C is enabled on your device.
```bash
# On Raspberry Pi
sudo raspi-config # Interface Options -> I2C -> Enable

# On Jetson Nano -> I2C is usually enabled by default on bus 1 or 7.
```

### 2. Build from Source
```bash
git clone https://github.com/cryptofyre/RackDisp-I2C.git
cd RackDisp-I2C

# Install dependencies
go mod tidy

# Build the binary
go build -o rackdisp ./cmd/rackdisp
```

### 3. Install as Service
RackDisp includes a helper command to install itself as a systemd service.
```bash
sudo ./rackdisp install
```
Follow the on-screen instructions to enable and start the service.

## 📖 Usage

### Running on Hardware
Simply run the binary. It will attempt to auto-detect the I2C bus.
```bash
./rackdisp
```

**Options:**
- `--device=ssd1306`: Force hardware mode.
- `--refresh=2s`: Change how often system stats (CPU/RAM) are updated (default 2s).
- `--rotate`: Rotate the display 180 degrees (useful for rack mounting).

### Running the Emulator
Want to test the UI on your dev machine? Use the emulator!
```bash
go run . --device=emulator
# OR if built:
./rackdisp --device=emulator
```

**Emulator Controls:**
- **Space**: Cycle to the next page immediately.
- **N**: Toggle "Internet Connected" icon.
- **W**: Toggle "Workload" icon (and inject mock AI data).

## 🧠 AI Workload Monitoring (Ollama)

RackDisp automatically detects if [Ollama](https://ollama.com/) is running and querying its API.
- **Config**: By default, it queries `http://10.10.70.159:11434`. 
- **Modify**: Edit `internal/stats/ollama.go` to change the IP address if your Ollama instance is elsewhere.

## 📁 Project Structure

```
├── cmd/rackdisp/       # Main entry point
├── internal/
│   ├── display/        # Hardware drivers (SSD1306) & Emulator
│   ├── stats/          # System data collection (CPU, Ollama, etc.)
│   └── ui/             # Graphics context & drawing logic
└── Systemd Service     # Generated via 'install' command
```

## 🖨️ 3D Printing

For a 1U chassis, I've developed a part for converting a 10 inch rack slot to a 19 inch rack slot; But in the blank space I've added a spot just for mounting the SSD1306 display. This can go into any of your devices that support such a thing.

Check it out on MakerWorld: tbd

## 🤝 Contributing

Contributions are welcome! Feel free to open an issue or submit a Pull Request.

## 📄 License

MIT License. Free to use and modify, I'd love to see what projects you come up with!
