package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

const serviceTemplate = `[Unit]
Description=RackDisp-I2C Display Service
After=network.target

[Service]
Type=simple
ExecStart={{.ExecPath}}
Restart=always
RestartSec=5
User=root
Group=root

[Install]
WantedBy=multi-user.target
`

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Generate and install systemd service",
	Run:   runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		return
	}
	execPath, _ = filepath.Abs(execPath)

	data := struct {
		ExecPath string
	}{
		ExecPath: execPath,
	}

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	fileName := "rackdisp.service"
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", fileName, err)
		return
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Printf("Error writing template: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated %s\n\n", fileName)
	fmt.Println("To install, run:")
	fmt.Printf("  sudo mv %s /etc/systemd/system/\n", fileName)
	fmt.Println("  sudo systemctl daemon-reload")
	fmt.Println("  sudo systemctl enable --now rackdisp.service")
}
