package service

import "github.com/godcong/aliyun-media-censor/config"

// service ...
type service struct {
	grpc *GRPCServer
	rest *RestServer
}

var server *service

// Start ...
func Start() {
	cfg := config.Config()

	server = &service{
		grpc: NewGRPCServer(cfg),
		rest: NewRestServer(cfg),
		//queue: NewQueueServer(cfg),
	}

	server.rest.Start()
	server.grpc.Start()
}

// Stop ...
func Stop() {
	server.rest.Stop()
	server.grpc.Stop()
}
