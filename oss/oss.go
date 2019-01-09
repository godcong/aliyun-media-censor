package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"path/filepath"
)

type OSS struct {
	Config Config
	Bucket *oss.Bucket
}

func newOSS(config Config, bucket *oss.Bucket) *OSS {
	return &OSS{Config: config, Bucket: bucket}
}

type Config struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	downloadInfo    *DownloadInfo
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
		Progress:   nil,
	}
}

//type ProgressListenFunc func(event *oss.ProgressEvent)

func (i *DownloadInfo) RegisterListener(lis ProgressListener) {
	i.Progress = oss.Progress(lis)
}

type ProgressListener interface {
	ProgressChanged(event *oss.ProgressEvent)
}

type progress struct {
}

// 定义进度变更事件处理函数。
func (listener *progress) ProgressChanged(event *oss.ProgressEvent) {
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

func (o *OSS) Download(objectKey string) error {
	di := o.Config.DownloadInfo()
	fp := filepath.Join(di.DirPath, objectKey)
	err := o.Bucket.DownloadFile(objectKey, fp, di.PartSize, di.Routines, di.Progress, di.Checkpoint)
	if err != nil {
		return err
	}
	return nil
}
