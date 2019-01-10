package green

import (
	"errors"
	"github.com/json-iterator/go"
	"log"
)

type Responder interface {
	GetHttpContentBytes() []byte
	GetHttpContentString() string
	IsSuccess() bool
}

func ResponseToResultData(r Responder) (*ResultData, error) {
	if !r.IsSuccess() {
		return &ResultData{}, errors.New("no success")
	}
	log.Println(r.GetHttpContentString())
	var d ResultData
	err := jsoniter.Unmarshal(r.GetHttpContentBytes(), &d)
	if err != nil {
		return &ResultData{}, nil
	}
	return &d, nil
}
