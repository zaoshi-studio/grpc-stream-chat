package main

import (
	"github.com/spf13/cobra"
	"github.com/zaoshi-studio/grpc-stream-chat/client"
	"github.com/zaoshi-studio/grpc-stream-chat/server"
	"log"
)

func main() {
	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(server.NewCmd())
	rootCmd.AddCommand(client.NewCmd())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
