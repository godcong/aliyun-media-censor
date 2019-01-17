package main

import (
	"context"
	"fmt"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/godcong/aliyun-media-censor/service"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	file, err := os.OpenFile("censor.log", os.O_SYNC|os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Lshortfile | log.Ldate)
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	oss.StartQueue(ctx, 5)
	service.Start()
	//start
	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		cancel()
		done <- true
	}()
	<-done
}
