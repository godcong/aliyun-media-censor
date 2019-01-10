package oss

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"
)

func TestStartQueue(t *testing.T) {
	// 根据命令行指定文件名创建 profile 文件
	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal(err)
	}
	// 开启 CPU profiling
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	RegisterCallback(func(strings chan<- string, info *QueueInfo) {
		println("queue", info.ObjectKey)
		//time.Sleep(1 * time.Millisecond)
		strings <- info.ObjectKey
	})
	StartQueue(context.Background(), 5)

	for i := 0; i < 100; i++ {
		Push(&QueueInfo{
			ObjectKey: strconv.Itoa(i),
		})
		time.Sleep(1 * time.Millisecond)
	}
	time.Sleep(30 * time.Second)
	StopQueue()
	time.Sleep(10 * time.Second)
}
