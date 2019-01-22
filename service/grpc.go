package service

import (
	"context"
	"fmt"
	"github.com/godcong/aliyun-media-censor/config"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/proto"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"syscall"
)

// GRPCServer ...
type GRPCServer struct {
	config *config.Configure
	server *grpc.Server
	Type   string
	Port   string
	Path   string
}

// Validate ...
func (s *GRPCServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.CensorReply, error) {
	var rd []*green.Result

	qi := QueueInfo{
		ObjectKey:    req.ObjectKey,
		ID:           req.ID,
		ValidateType: req.ValidateType.String(),
	}
	var data []*green.ResultData
	var err error
	switch req.ValidateType {
	case proto.CensorValidateType_Frame:
		Push(&qi)
		data = []*green.ResultData{}
	case proto.CensorValidateType_JPG:
		data, err = ParseValidateDo(&qi, func(url string) (data *green.ResultData, e error) {
			return green.ImageSyncScan(&green.BizData{
				Scenes: []string{"porn", "terrorism", "ad", "live", "sface"},
				Tasks: []green.Task{
					{
						DataID: uuid.NewV1().String(),
						URL:    url,
					},
				},
			})
		})
	case proto.CensorValidateType_Video:
		data, err = ParseValidateDo(&qi, func(url string) (data *green.ResultData, e error) {
			return green.VideoAsyncScan(&green.BizData{
				Scenes:      []string{"porn", "terrorism", "ad", "live", "sface"},
				AudioScenes: []string{"antispam"},
				Tasks: []green.Task{
					{
						DataID:    uuid.NewV1().String(),
						URL:       url,
						Interval:  30,
						MaxFrames: 200,
					},
				}})
		})
	}

	m, _ := jsoniter.MarshalToString(rd)
	return Result(&proto.CensorReplyDetail{
		ID:   req.ID,
		Json: m,
	}), nil
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
func NewGRPCServer(cfg *config.Configure) *GRPCServer {
	return &GRPCServer{
		Type: config.DefaultString(cfg.GRPC.Type, Type),
		Port: config.DefaultString(cfg.GRPC.Port, ":7786"),
		Path: config.DefaultString(cfg.GRPC.Path, "/tmp/censor.sock"),
	}
}

// Start ...
func (s *GRPCServer) Start() {
	if !s.config.GRPC.Enable {
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
