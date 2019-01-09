package oss

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestStartQueue(t *testing.T) {
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
