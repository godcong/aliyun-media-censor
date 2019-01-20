package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// RestServer ...
type RestServer struct {
	*gin.Engine
	Server *http.Server
}

// NewRestServer ...
func NewRestServer() *RestServer {
	eng := gin.Default()
	return &RestServer{
		Engine: eng,
		Server: &http.Server{
			Addr:    ":7787",
			Handler: eng,
		},
	}
}

// Start ...
func (s *RestServer) Start() {
	Router(s.Engine)

	go func() {
		log.Printf("[GIN-debug] Listening and serving HTTP on %s\n", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

}

// Stop ...
func (s *RestServer) Stop() error {
	return s.Server.Shutdown(nil)
}
