package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"path/filepath"
	"sync"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	downloadInfo    *DownloadInfo
}

type OSS struct {
	Config Config
	Bucket *oss.Bucket
}

var server1 *OSS
var server2 *OSS

func init() {
	var err error
	once := sync.Once{}
	once.Do(func() {
		//server init

		if err != nil {
			panic(err)
		}

	})
}

func newOSS(config Config, bucket *oss.Bucket) *OSS {
	return &OSS{Config: config, Bucket: bucket}
}

func (c *Config) DownloadInfo() *DownloadInfo {
	if c.downloadInfo == nil {
		c.downloadInfo = NewDownloadInfo()
	}
	return c.downloadInfo
}

func (c *Config) SetDownloadInfo(downloadInfo *DownloadInfo) {
	c.downloadInfo = downloadInfo
}

type DownloadInfo struct {
	DirPath    string
	PartSize   int64
	Routines   oss.Option
	Checkpoint oss.Option
	Progress   oss.Option
}

func NewDownloadInfo() *DownloadInfo {
	return &DownloadInfo{
		DirPath:    "./download",
		PartSize:   100 * 1024 * 1024,
		Routines:   oss.Routines(5),
		Checkpoint: oss.Checkpoint(true, "./cp"),
		Progress:   oss.Progress(&progress{}),
	}
}

func (i *DownloadInfo) RegisterListener(lis Progress) {
	i.Progress = oss.Progress(lis)
}

type Progress interface {
	ProgressChanged(event *oss.ProgressEvent)
	SetObjectKey(objectKey string)
	ObjectKey() string
	Option() oss.Option
}

type progress struct {
	objectKey string
}

func NewProgress() Progress {
	return &progress{}
}

func (p *progress) Option() oss.Option {
	return oss.Progress(p)
}

func (p *progress) ObjectKey() string {
	return p.objectKey
}

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

func NewOSS(config Config) (*OSS, error) {
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %s", err)
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %s", err)
	}
	return newOSS(config, bucket), nil
}

func (o *OSS) Download(p Progress) error {
	di := o.Config.DownloadInfo()
	fp := filepath.Join(di.DirPath, p.ObjectKey())
	err := o.Bucket.DownloadFile(p.ObjectKey(), fp, di.PartSize, di.Routines, p.Option(), di.Checkpoint)
	if err != nil {
		return err
	}
	return nil
}

func (o *OSS) Upload(p Progress) error {
	di := o.Config.DownloadInfo()
	fp := filepath.Join(di.DirPath, p.ObjectKey())
	err := o.Bucket.UploadFile(p.ObjectKey(), fp, di.PartSize, di.Routines, p.Option(), di.Checkpoint)
	if err != nil {
		return err
	}
	return nil
}

func (o *OSS) URL(p Progress) (string, error) {
	signedURL, err := o.Bucket.SignURL(p.ObjectKey(), oss.HTTPGet, 60*60*24)
	if err != nil {
		return "", err
	}
	return signedURL, err

}

func Server1() *OSS {
	return server1
}

func Server2() *OSS {
	return server2
}

func Server3() {

}
