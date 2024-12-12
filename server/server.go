package server

import (
	"context"
	"fmt"
	"github.com/zaoshi-studio/grpc-stream-chat/pb/protocol"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type Message struct {
	Content string
	Sayer   string
	WaitC   chan struct{}
}

type Server struct {
	protocol.UnimplementedChatServer
	clients  map[string]*Client
	mu       sync.Mutex
	contentC chan Message
	config   Config
	grpcSvr  *grpc.Server
	ctx      context.Context
	cancel   context.CancelFunc
	eg       errgroup.Group
}

func New(config Config) *Server {
	svr := &Server{
		config:   config,
		clients:  make(map[string]*Client),
		contentC: make(chan Message, 10),
	}

	svr.grpcSvr = grpc.NewServer()
	protocol.RegisterChatServer(svr.grpcSvr, svr)
	return svr
}

func (svr *Server) Run(ctx context.Context) error {
	svr.ctx, svr.cancel = context.WithCancel(ctx)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", svr.config.Port))
	if err != nil {
		return err
	}

	svr.eg.Go(func() error {
		if err = svr.grpcSvr.Serve(listener); err != nil {
			return err
		}
		return nil
	})

	svr.eg.Go(func() error {
		for {
			select {
			case <-svr.ctx.Done():
				return nil
			case message := <-svr.contentC:

				log.Printf("client: %v say: %v", message.Sayer, message.Content)
				svr.mu.Lock()
				for _, client := range svr.clients {
					if message.Sayer == client.id {
						continue
					}
					if err := client.stream.SendMsg(&protocol.SayRsp{
						Content: message.Content,
						Sayer:   message.Sayer,
					}); err != nil {
						log.Printf("send message to client: %v error: %v", client.id, err)
					}
				}
				svr.mu.Unlock()
				close(message.WaitC)
			}
		}
	})

	return nil
}

func (svr *Server) Stop() error {
	if svr.cancel != nil {
		svr.cancel()
	}
	svr.grpcSvr.GracefulStop()
	if err := svr.eg.Wait(); err != nil {
		return err
	}
	return nil
}
