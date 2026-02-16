package cli

import (
	"context"
	"fmt"
	"os"

	"cryptofy.re/rackdisp/internal/app"
	"cryptofy.re/rackdisp/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg config.Config
)

var rootCmd = &cobra.Command{
	Use:   "rackdisp",
	Short: "RackDisp-I2C displays system stats on I2C OLED screens",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.Run(context.Background(), cfg)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfg.DeviceType, "device", "d", "auto", "Device type: auto, ssd1306, emulator")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Rotate, "rotate", "r", false, "Rotate display 180 degrees")
	rootCmd.PersistentFlags().IntVar(&cfg.Refresh, "refresh", 2, "Refresh stats interval in seconds")
}
