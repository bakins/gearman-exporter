package main

import (
	"fmt"
	"os"

	exporter "github.com/bakins/gearman-exporter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	addr        *string
	gearmanAddr *string
)

func serverCmd(cmd *cobra.Command, args []string) {

	logger, err := exporter.NewLogger()
	if err != nil {
		panic(err)
	}

	e, err := exporter.New(
		exporter.SetAddress(*addr),
		exporter.SetGearmanAddress(*gearmanAddr),
		exporter.SetLogger(logger),
	)

	if err != nil {
		logger.Fatal("failed to create exporter", zap.Error(err))
	}

	if err := e.Run(); err != nil {
		logger.Fatal("failed to run exporter", zap.Error(err))
	}
}

var rootCmd = &cobra.Command{
	Use:   "gearman-exporter",
	Short: "Gearman metrics exporter",
	Run:   serverCmd,
}

func main() {
	addr = rootCmd.PersistentFlags().StringP("addr", "", "127.0.0.1:8080", "listen address for metrics handler")
	gearmanAddr = rootCmd.PersistentFlags().StringP("gearmand", "", "127.0.0.1:4730", "address of gearmand")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("root command failed: %v", err)
		os.Exit(-2)
	}
}
