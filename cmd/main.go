package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ozeemandias/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const grpcPort = 50052

type server struct {
	chat_v1.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, req *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	log.Printf("%v", req)

	return &chat_v1.CreateResponse{}, nil
}

func (s *server) Delete(_ context.Context, req *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("%v", req)

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("%v", req)

	return &emptypb.Empty{}, nil
}

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chat_v1.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", ln.Addr())

	if err = s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
