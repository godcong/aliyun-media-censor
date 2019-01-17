package oss

import (
	"context"
	"log"
	"sync"
	"time"
)

// DefaultSleepTime ...
const DefaultSleepTime = 5 * time.Second

// QueueInfo ...
type QueueInfo struct {
	ObjectKey    string
	CallbackURL  string
	RequestKey   string
	CallbackFunc func(chan<- string, *QueueInfo) `json:"-"`
}

// Queue ...
type Queue struct {
	callback CallbackFunc
	infos    []*QueueInfo
	lock     sync.RWMutex
}

// NewStreamQueue ...
func NewStreamQueue() *Queue {
	return &Queue{
		infos: []*QueueInfo{},
	}
}

// Push ...
func (s *Queue) Push(info *QueueInfo) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.infos = append(s.infos, info)
}

// Pop ...
func (s *Queue) Pop() *QueueInfo {
	s.lock.Lock()
	defer s.lock.Unlock()
	info := s.infos[0]
	s.infos = s.infos[1:len(s.infos)]

	return info
}

// Front ...
func (s *Queue) Front() *QueueInfo {
	s.lock.RLock()
	defer s.lock.RUnlock()
	info := s.infos[0]
	return info
}

// IsEmpty ...
func (s *Queue) IsEmpty() bool {
	return len(s.infos) == 0
}

// Size ...
func (s *Queue) Size() int {
	return len(s.infos)
}

// Clear ...
func (s *Queue) Clear() {
	s.infos = []*QueueInfo{}
}

// CallbackFunc ...
type CallbackFunc func(chan<- string, *QueueInfo)

var (
	globalCallback CallbackFunc
	globalCancel   context.CancelFunc
	queue          = NewStreamQueue()
)

// Push ...
func Push(v *QueueInfo) {
	queue.Push(v)
}

// Pop ...
func Pop() *QueueInfo {
	if !queue.IsEmpty() {
		return queue.Pop()
	}
	log.Println("nothing pop")
	return nil
}

// RegisterCallback ...
func RegisterCallback(fn CallbackFunc) {
	globalCallback = fn
}

func nilQueue(s chan<- string) {
	println("start queue")
	time.Sleep(DefaultSleepTime)
	s <- "nil"
}

// StartQueue ...
func StartQueue(ctx context.Context, process int) {
	var c context.Context
	c, globalCancel = context.WithCancel(ctx)
	//run with a new go channel
	go func() {
		threads := make(chan string, process)

		for i := 0; i < process; i++ {
			go nilQueue(threads)
		}

		for {
			select {
			case v := <-threads:
				println("success:", v)
				if s := Pop(); s != nil {
					if s.CallbackFunc != nil {
						go s.CallbackFunc(threads, s)
					} else if globalCallback != nil {
						go globalCallback(threads, s)
					} else {

					}
				} else {
					time.Sleep(DefaultSleepTime)
					go nilQueue(threads)
				}
			case <-c.Done():
				return
				//default:
				//	println("default sleep 5 seconds")
				//	time.Sleep(DefaultSleepTime)
			}
		}
	}()
}

// StopQueue ...
func StopQueue() {
	globalCancel()
}
