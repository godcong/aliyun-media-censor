package service

// service ...
// service ...
type service struct {
	grpc *GRPCServer
	rest *RestServer
}

var server *service

func init() {
	server = &service{
		grpc: NewGRPCServer(),
		rest: NewRestServer(),
	}

}

// Start ...
func Start() {
	server.rest.Start()
	server.grpc.Start()
}

// Stop ...
func Stop() {
	server.rest.Stop()
	server.grpc.Stop()
}
