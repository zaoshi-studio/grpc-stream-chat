package server

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
	}

	var (
		configPath string
	)

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file path")

	cmd.Run = func(_ *cobra.Command, args []string) {
		cfgBytes, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("read config file failed: %v", err)
		}

		var config Config
		if err = json.Unmarshal(cfgBytes, &config); err != nil {
			log.Fatalf("unmarshal config file failed: %v", err)
		}

		baseCtx := context.Background()
		svr := New(config)

		if err = svr.Run(baseCtx); err != nil {
			log.Fatalf("server run failed: %v", err)
		}

		signalC := make(chan os.Signal, 1)
		signal.Notify(signalC, os.Interrupt)
		<-signalC

		if err = svr.Stop(); err != nil {
			log.Fatalf("server stop failed: %v", err)
		}

	}

	return cmd
}
