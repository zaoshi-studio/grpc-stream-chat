package client

import (
	"bufio"
	"context"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strings"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "client",
	}

	var (
		serverAddr string
	)

	cmd.PersistentFlags().StringVarP(&serverAddr, "server", "s", "localhost:8080", "server address")

	cmd.Run = func(_ *cobra.Command, args []string) {

		baseCtx := context.Background()
		client := New(serverAddr)

		if err := client.Run(baseCtx); err != nil {
			log.Fatalf("server run failed: %v", err)
		}

		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}

			content := strings.Trim(s, "\r\n")

			client.contentC <- content
		}

		if err := client.Stop(); err != nil {
			log.Fatalf("client stop failed: %v", err)
		}

	}

	return cmd
}
