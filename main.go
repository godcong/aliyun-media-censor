//go:generate protoc --go_out=plugins=grpc:./proto censor.proto
package main

import (
	"flag"
	"fmt"
	"github.com/godcong/aliyun-media-censor/config"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/godcong/aliyun-media-censor/service"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("path", "config.toml", "load config file from path")

func main() {
	file, err := os.OpenFile("censor.log", os.O_SYNC|os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Lshortfile | log.Ldate)
	flag.Parse()
	err = config.Initialize(*configPath)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//ctx, cancel := context.WithCancel(context.Background())

	service.Start()
	//start
	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		//cancel()
		service.Stop()
		done <- true
	}()
	<-done
}
