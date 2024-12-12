package server

import (
	"github.com/zaoshi-studio/grpc-stream-chat/pb/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
)

func (svr *Server) Say(stream grpc.BidiStreamingServer[protocol.SayReq, protocol.SayRsp]) error {
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return nil
	}

	client := &Client{
		id:     p.Addr.String(),
		stream: stream,
	}

	svr.mu.Lock()
	if _, ok := svr.clients[client.id]; ok {
		delete(svr.clients, client.id)
	}
	svr.clients[client.id] = client
	svr.mu.Unlock()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		message := Message{
			Content: req.GetContent(),
			Sayer:   client.id,
			WaitC:   make(chan struct{}),
		}

		svr.contentC <- message
		<-message.WaitC
	}
}
