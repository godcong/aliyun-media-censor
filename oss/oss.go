package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/godcong/aliyun-media-censor/config"

	"log"
	"os"
	"path/filepath"
)

// OSS ...
type OSS struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	*oss.Bucket
}

// BucketServer ...
type BucketServer struct {
	config *config.Configure
	server []*OSS
	info   *DownloadInfo
}

var server *BucketServer

// NewOSS ...
func NewOSS(oss *config.OSS) *OSS {
	return &OSS{
		Endpoint:        config.DefaultString(oss.Endpoint, ""),
		AccessKeyID:     config.DefaultString(oss.AccessKeyID, ""),
		AccessKeySecret: config.DefaultString(oss.AccessKeySecret, ""),
		BucketName:      config.DefaultString(oss.BucketName, ""),
	}
}

// Connect ...
func (o *OSS) Connect() error {
	client, err := oss.New(o.Endpoint, o.AccessKeyID, o.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("failed to create new client: %s", err)
	}

	bucket, err := client.Bucket(o.BucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %s", err)
	}
	o.Bucket = bucket
	return nil
}

// NewBucketServer ...
func NewBucketServer(cfg *config.Configure) *BucketServer {
	var s BucketServer
	s.config = cfg
	for _, val := range cfg.OSS {
		oss := NewOSS(&val)
		err := oss.Connect()
		if err != nil {
			log.Println(err)
			panic(err)
		}
		s.server = append(s.server, oss)
	}
	return &s
}

// InitOSS ...
func InitOSS(cfg *config.Configure) {
	server = NewBucketServer(cfg)
}

// Server ...
func (s *BucketServer) Server(idx ...int) *OSS {
	if idx == nil {
		return s.server[0]
	}
	return s.server[idx[0]]
}

// Info ...
func (s *BucketServer) Info() *DownloadInfo {
	if s.info == nil {
		s.info = NewDownloadInfo()
	}
	return s.info
}

// SetInfo ...
func (s *BucketServer) SetInfo(info *DownloadInfo) {
	s.info = info
}

// DownloadInfo ...
type DownloadInfo struct {
	config     *config.Configure
	DirPath    string
	PartSize   int64
	Routines   oss.Option
	Checkpoint oss.Option
	Progress   oss.Option
}

// NewDownloadInfo ...
func NewDownloadInfo(cfg *config.Configure) *DownloadInfo {
	return &DownloadInfo{
		config:     cfg,
		DirPath:    config.DefaultString(cfg.Media.Download, "download"),
		PartSize:   100 * 1024 * 1024,
		Routines:   oss.Routines(5),
		Checkpoint: oss.Checkpoint(true, "./cp"),
		Progress:   oss.Progress(&progress{}),
	}
}

// RegisterListener ...
func (i *DownloadInfo) RegisterListener(lis Progress) {
	i.Progress = oss.Progress(lis)
}

// Progress ...
type Progress interface {
	ProgressChanged(event *oss.ProgressEvent)
	SetObjectKey(objectKey string)
	ObjectKey() string
	Path() string
	SetPath(path string)
	Option() oss.Option
}

type progress struct {
	objectKey string
	path      string
}

// Path ...
func (p *progress) Path() string {
	return p.path
}

// SetPath ...
func (p *progress) SetPath(path string) {
	p.path = path
}

// NewProgress ...
func NewProgress() Progress {
	return &progress{}
}

// Option ...
func (p *progress) Option() oss.Option {
	return oss.Progress(p)
}

// ObjectKey ...
func (p *progress) ObjectKey() string {
	return p.objectKey
}

// SetObjectKey ...
func (p *progress) SetObjectKey(objectKey string) {
	p.objectKey = objectKey
}

// 定义进度变更事件处理函数。
func (p *progress) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)

	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

// FileName ...
func FileName(objectKey string) string {
	_, file := filepath.Split(objectKey)
	return file
}

// Download ...
func (s *BucketServer) Download(p Progress) error {
	di := s.Info()
	path := di.DirPath
	if p.Path() != "" {
		path = p.Path()
	}

	fp := filepath.Join(path, FileName(p.ObjectKey()))
	dir, _ := filepath.Split(fp)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Println(err)
		//ignore error
	}
	err = s.Server().DownloadFile(p.ObjectKey(), fp, di.PartSize, di.Routines, p.Option(), di.Checkpoint)
	if err != nil {
		return err
	}
	return nil
}

// Upload ...
func (s *BucketServer) Upload(p Progress) error {
	di := s.Info()
	path := di.DirPath
	if p.Path() != "" {
		path = p.Path()
	}
	fp := filepath.Join(path, p.ObjectKey())
	err := s.Server().UploadFile(p.ObjectKey(), fp, di.PartSize, di.Routines, p.Option(), di.Checkpoint)
	if err != nil {
		return err
	}
	return nil
}

// URL ...
func (s *BucketServer) URL(p Progress) (string, error) {
	signedURL, err := s.Server().SignURL(p.ObjectKey(), oss.HTTPGet, 60*60*24)
	if err != nil {
		return "", err
	}
	return signedURL, err

}

// IsExist ...
func (s *BucketServer) IsExist(p Progress) bool {
	exist, err := s.Server().IsObjectExist(p.ObjectKey())
	if err != nil {
		log.Println(err)
		return false
	}
	return exist
}

// Server ...
func Server() *BucketServer {
	return server
}
