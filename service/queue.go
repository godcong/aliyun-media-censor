package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/godcong/aliyun-media-censor/config"
	"github.com/godcong/aliyun-media-censor/ffmpeg"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/godcong/aliyun-media-censor/util"
	"github.com/json-iterator/go"
	"log"
	"math"
	"time"
)

// StatusUploaded 已上传
const StatusUploaded = "uploaded"

// StatusDownloading 正在下载
const StatusDownloading = "downloading"

// StatusDownloaded 已下载
const StatusDownloaded = "downloaded"

// StatusQueuing 队列中
const StatusQueuing = "queuing"

// StatusTransferring 转换中
const StatusTransferring = "transferring"

// StatusFileWrong 文件错误
const StatusFileWrong = "wrong file"

//StatusFailed 失败
const StatusFailed = "failed"

// StatusFinished 完成
const StatusFinished = "finished"

// HandleFunc ...
type HandleFunc func(name, key string) error

// QueueServer ...
type QueueServer struct {
	*redis.Client
	config    *config.Configure
	cancel    context.CancelFunc
	Processes int
}

var globalQueue *QueueServer

// Push ...
func Push(v *QueueInfo) {
	globalQueue.Push(v)
}

// Push ...
func (s *QueueServer) Push(v *QueueInfo) {
	s.RPush("censor_queue", v.JSON())
}

// Pop ...
func Pop() *QueueInfo {
	return globalQueue.Pop()
}

// Pop ...
func (s *QueueServer) Pop() *QueueInfo {
	pop := s.LPop("censor_queue").Val()
	return ParseInfo(pop)
}

func validating(ch chan<- string, info *QueueInfo) {
	var err error
	chanRes := info.FileName()
	defer func() {
		if err != nil {
			log.Println(err)
			globalQueue.Set(info.ID, StatusFailed, 0)
			chanRes = info.FileName() + ":[" + err.Error() + "]"
		}
		ch <- chanRes
	}()
	p := oss.NewProgress()
	p.SetObjectKey(info.ObjectKey)
	p.SetPath(info.ID)
	server := oss.Server()
	if !server.IsExist(p) {
		err = fmt.Errorf("object [%s] is not exist", info.ObjectKey)
		return
	}
	err = server.Download(p)
	ts, err := ffmpeg.TransferJPG(info.FileSource, info.ObjectKey)
	if err != nil {
		return
	}
	log.Println(ts)
	files, err := util.FileList(info.FileSource + info.ObjectKey)
	if err != nil {
		return
	}

	fileLen := len(files)
	steps := int(math.Ceil(float64(fileLen) / 64))

	rd := make(chan *green.ResultData, steps)

	sta, end := 0, 64
	for i := 0; i < int(steps); i++ {
		if end > fileLen {
			end = fileLen
		}
		go green.ProcessFrame(rd, files[sta:end], i*64+64, info.FileDest, info.ObjectKey)
		sta = i * 64
		end = sta + 64
	}

	var rds []*green.ResultData

	for i := 0; i < int(steps); i++ {
		select {
		case v := <-rd:
			if v != nil {
				rds = append(rds, v)
			}
		}
	}

	err = NewBack().Callback(info.ID, &ResultDataList{
		ID:     info.ID,
		Detail: rds,
	})
	return
}

type ResultDataList struct {
	ID     string              `json:"id"`
	Detail []*green.ResultData `json:"detail"`
}

func (list *ResultDataList) JSON() string {
	s, _ := jsoniter.MarshalToString(list)
	return s
}

func validateNothing(threads chan<- string) {
	time.Sleep(3 * time.Second)
	threads <- ""
}

// NewQueueServer ...
func NewQueueServer(cfg *config.Configure) *QueueServer {
	client := redis.NewClient(&redis.Options{
		Addr:     config.DefaultString(cfg.Queue.HostPort, ":6379"),
		Password: config.DefaultString(cfg.Queue.Password, ""), // no password set
		DB:       cfg.Queue.DB,                                 // use default DB
	})
	return &QueueServer{
		config: cfg,
		Client: client,
	}
}

// Start ...
func (s *QueueServer) Start() {
	pong, err := s.Ping().Result()
	if err != nil {
		panic(err)
	}
	globalQueue = s
	log.Println(pong)

	var c context.Context
	c, s.cancel = context.WithCancel(context.Background())
	//run with a new go channel
	go func() {
		threads := make(chan string, s.Processes)

		for i := 0; i < s.Processes; i++ {
			log.Println("start", i)
			go validateNothing(threads)
		}

		for {
			select {
			case v := <-threads:
				if v != "" {
					log.Println("success: ", v)
				}

				if s := Pop(); s != nil {
					go validating(threads, s)
				} else {
					go validateNothing(threads)
				}
			case <-c.Done():
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

// Stop ...
func (s *QueueServer) Stop() {
	if s.cancel == nil {
		return
	}
	s.cancel()
}
