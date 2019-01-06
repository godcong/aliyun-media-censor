package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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

func (o *OSS) Put() {

}
