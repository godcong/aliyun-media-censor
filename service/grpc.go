package service

import (
	"context"
	"fmt"
	"github.com/godcong/aliyun-media-censor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"syscall"
)

// GRPCServer ...
type GRPCServer struct {
	Type   string
	Port   string
	Path   string
	server *grpc.Server
}

func (s *GRPCServer) Validate(context.Context, *proto.ValidateRequest) (*proto.CensorReply, error) {
	panic("implement me")
}

type grpcBack struct {
	BackType string
	BackAddr string
}

// Result ...
func Result(detail *proto.CensorReplyDetail) *proto.CensorReply {
	return &proto.CensorReply{
		Code:    0,
		Message: "success",
		Detail:  detail,
	}
}

// NewGRPCServer ...
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		Type: DefaultString(config.GRPC.Type, Type),
		Port: DefaultString(config.GRPC.Port, ":7786"),
		Path: DefaultString(config.GRPC.Path, "/tmp/censor.sock"),
	}
}

// Start ...
func (s *GRPCServer) Start() {
	if !config.GRPC.Enable {
		return
	}
	s.server = grpc.NewServer()
	var lis net.Listener
	var port string
	var err error
	go func() {
		if s.Type == "unix" {
			_ = syscall.Unlink(s.Path)
			lis, err = net.Listen(s.Type, s.Path)
			port = s.Path
		} else {
			lis, err = net.Listen("tcp", s.Port)
			port = s.Port
		}

		if err != nil {
			panic(fmt.Sprintf("failed to listen: %v", err))
		}

		proto.RegisterCensorServiceServer(s.server, s)
		// Register reflection service on gRPC server.
		reflection.Register(s.server)
		log.Printf("Listening and serving TCP on %s\n", port)
		if err := s.server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

}

// Stop ...
func (s *GRPCServer) Stop() {
	s.server.Stop()
}
