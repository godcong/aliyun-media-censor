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
	config      *config.Configure
	encrypt     bool
	ID          string
	Key         string
	ObjectKey   string
	KeyURL      string
	KeyName     string
	KeyInfoName string
	KeyDest     string
	FileSource  string
	FileDest    string
	Callback    string
}

// NewQueueInfo ...
func NewQueueInfo() *QueueInfo {
	return &QueueInfo{
		encrypt:     false,
		ID:          uuid.NewV1().String(),
		Key:         "",
		KeyURL:      "",
		KeyName:     "",
		KeyInfoName: "",
		KeyDest:     "",
		//FileName:    util.GenerateRandomString(64),
		FileSource: "",
		FileDest:   "",
	}
}

// NewStreamerWithConfig ...
func NewStreamerWithConfig(cfg *config.Configure, id string) *QueueInfo {
	return &QueueInfo{
		encrypt:     false,
		ID:          config.DefaultString(id, uuid.NewV1().String()),
		KeyURL:      cfg.Media.KeyURL,
		KeyName:     cfg.Media.KeyFile,
		KeyInfoName: cfg.Media.KeyInfoFile,
		KeyDest:     cfg.Media.KeyDest,
		FileSource:  cfg.Media.Upload,
		FileDest:    cfg.Media.Transfer,
	}
}

// FileName ...
func (s *QueueInfo) FileName() string {
	_, file := filepath.Split(s.ObjectKey)
	return file
}

// Encrypt ...
func (s *QueueInfo) Encrypt() bool {
	return s.encrypt
}

// SetEncrypt ...
func (s *QueueInfo) SetEncrypt(encrypt bool) {
	s.encrypt = true
	s.KeyURL = s.config.Media.KeyURL
	s.KeyName = s.config.Media.KeyFile
	s.KeyInfoName = s.config.Media.KeyInfoFile
	s.KeyDest = s.config.Media.KeyDest
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
