package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/config"
	"log"
	"net/http"
)

// RestServer ...
type RestServer struct {
	*gin.Engine
	config *config.Configure
	server *http.Server
	Port   string
}

type restBack struct {
	config  *config.Configure
	BackURL string
	Version string
}

func (b *restBack) Callback(*QueueResult) error {
	panic("implement me")
}

// NewRestBack ...
func NewRestBack(cfg *config.Configure) QueueCallback {
	return &restBack{
		config:  cfg,
		BackURL: config.DefaultString(cfg.Callback.BackAddr, "localhost:7780"),
		Version: config.DefaultString(cfg.Callback.Version, "v0"),
	}
}

// NewRestServer ...
func NewRestServer(cfg *config.Configure) *RestServer {
	s := &RestServer{
		Engine: gin.Default(),
		config: cfg,
		Port:   config.DefaultString(cfg.REST.Port, ":7785"),
	}
	return s
}

// Start ...
func (s *RestServer) Start() {
	if !s.config.REST.Enable {
		return
	}

	Router(s.Engine)

	s.server = &http.Server{
		Addr:    s.Port,
		Handler: s.Engine,
	}
	go func() {
		log.Printf("Listening and serving HTTP on %s\n", s.Port)
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

}

// Stop ...
func (s *RestServer) Stop() {
	if err := s.server.Shutdown(nil); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
