package service

import (
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"github.com/satori/go.uuid"
	"testing"
)

// TestDownload ...
func TestDownload(t *testing.T) {
	queue = redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",              // no password set
		DB:       RedisQueueIndex, // use default DB
	})
	err := download(&QueueInfo{
		encrypt:     false,
		ID:          uuid.NewV1().String(),
		Key:         "",
		ObjectKey:   "origin/5c35cc6b5ec8a925a4143001/e84976d3567f339635eb0d49cccae72c/0050.mp4",
		KeyURL:      "",
		KeyName:     "",
		KeyInfoName: "",
		KeyDest:     "",
		FileSource:  "upload",
		FileDest:    "validating",
	})
	t.Log(err)
}

// TestTransfer ...
func TestTransfer(t *testing.T) {
	err := toM3U8("16118190-1ae9-11e9-b250-00155d33ca31", "upload/16118190-1ae9-11e9-b250-00155d33ca31/0050.mp4", "validating")
	t.Log(err)
}

// TestCommit ...
func TestCommit(t *testing.T) {
	detail, e := commitToIPNS("16118190-1ae9-11e9-b250-00155d33ca31", "validating/16118190-1ae9-11e9-b250-00155d33ca31")
	t.Log(detail, e)
	var cr QueueResult
	//
	err := mapstructure.Decode(detail, &cr)
	t.Log(err, cr)
}
