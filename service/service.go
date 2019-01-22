package service

import (
	"github.com/godcong/aliyun-media-censor/config"
	"github.com/godcong/aliyun-media-censor/oss"
	"log"
)

// service ...
type service struct {
	grpc  *GRPCServer
	rest  *RestServer
	queue *QueueServer
}

var server *service

// Start ...
func Start() {
	cfg := config.Config()

	server = &service{
		grpc:  NewGRPCServer(cfg),
		rest:  NewRestServer(cfg),
		queue: NewQueueServer(cfg),
	}

	log.Println("run main")
	oss.InitOSS(config.Config())

	server.rest.Start()
	server.grpc.Start()

	server.queue.Processes = 5
	server.queue.Start()

}

// Stop ...
func Stop() {
	server.rest.Stop()
	server.grpc.Stop()
	server.queue.Stop()
}

// NewBack ...
func NewBack() StreamerCallback {
	cfg := config.Config()
	if cfg != nil && cfg.Callback.Type == "grpc" {
		return NewGRPCBack(cfg)
	}
	return NewRestBack(cfg)
}
