package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/config"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// ContentTypeJSON ...
const ContentTypeJSON = "application/json"

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

func (b *restBack) Callback(id string, res *ResultDataList) error {
	back := filepath.Join(CheckPrefix(b.BackURL), b.Version, "censor/callback")
	log.Println(back)

	resp, err := http.Post(back, ContentTypeJSON, strings.NewReader(res.JSON()))
	if err != nil {
		log.Println(err)
		return err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes), err)
	return err
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

// CheckPrefix ...
func CheckPrefix(url string) string {
	if strings.Index(url, "http") != 0 {
		return "http://" + url
	}
	return url
}
