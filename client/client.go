package client

import (
	"context"
	"github.com/zaoshi-studio/grpc-stream-chat/pb/protocol"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

type Client struct {
	serverAddr string
	contentC   chan string
	sayStream  grpc.BidiStreamingClient[protocol.SayReq, protocol.SayRsp]
	ctx        context.Context
	cancel     context.CancelFunc
	eg         errgroup.Group
}

func New(serverAddr string) *Client {
	client := &Client{
		serverAddr: serverAddr,
		contentC:   make(chan string, 1),
	}
	return client
}

func (c *Client) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	client, err := grpc.NewClient(c.serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	chatClient := protocol.NewChatClient(client)

	sayStream, err := chatClient.Say(c.ctx)
	if err != nil {
		return err
	}
	c.sayStream = sayStream

	c.eg.Go(func() error {
		for {
			rsp, err := c.sayStream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				continue
			}
			log.Printf("recv: %s", rsp)
		}
	})

	c.eg.Go(func() error {
		for {
			select {
			case <-c.ctx.Done():
				return nil
			case content := <-c.contentC:
				log.Printf("content: %v", content)
				if err := c.sayStream.SendMsg(&protocol.SayReq{
					Content: content,
				}); err != nil {
					log.Printf("send msg failed: %v", err)
					continue
				}
			}
		}
	})
	return nil
}

func (c *Client) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}

	return nil
}
