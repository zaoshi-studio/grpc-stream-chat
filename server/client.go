package server

import (
	"github.com/zaoshi-studio/grpc-stream-chat/pb/protocol"
	"google.golang.org/grpc"
)

type Client struct {
	id     string
	stream grpc.BidiStreamingServer[protocol.SayReq, protocol.SayRsp]
}
