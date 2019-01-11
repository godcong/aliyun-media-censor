package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// service ...
type service struct {
	*gin.Engine
	Server *http.Server
}

var server *service

func init() {
	server = defaultEngine()
}

// Start ...
func Start() {
	Router(server.Engine)

	go func() {
		log.Printf("[GIN-debug] Listening and serving HTTP on %s\n", server.Server.Addr)
		if err := server.Server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

}

// Stop ...
func Stop() error {
	return server.Server.Shutdown(nil)
}

func defaultEngine() *service {
	eng := gin.Default()
	return &service{
		Engine: eng,
		Server: &http.Server{
			Addr:    ":7789",
			Handler: eng,
		},
	}
}
