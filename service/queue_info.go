package service

import (
	"github.com/godcong/aliyun-media-censor/config"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
	"log"
	"path/filepath"
)

// QueueCallback ...
type QueueCallback interface {
	Callback(*QueueResult) error
}

// QueueInfo ...
type QueueInfo struct {
	config        *config.Configure
	ID            string
	ObjectKey     string
	ValidateType  string
	ProcessMethod string
	Callback      string
	FileSource    string
	FileDest      string
}

// NewStreamerWithConfig ...
func NewStreamerWithConfig(cfg *config.Configure, id string) *QueueInfo {
	return &QueueInfo{
		config:     cfg,
		ID:         config.DefaultString(id, uuid.NewV1().String()),
		FileSource: cfg.Media.Upload,
		FileDest:   cfg.Media.Transfer,
	}
}

// FileName ...
func (s *QueueInfo) FileName() string {
	_, file := filepath.Split(s.ObjectKey)
	return file
}

// JSON ...
func (s *QueueInfo) JSON() string {
	st, err := jsoniter.MarshalToString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	return st
}

// ParseInfo ...
func ParseInfo(s string) *QueueInfo {
	var st QueueInfo
	err := jsoniter.UnmarshalFromString(s, &st)
	if err != nil {
		return nil
	}
	return &st
}
