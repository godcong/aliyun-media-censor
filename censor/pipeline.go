package censor

import (
	"errors"
	"github.com/godcong/aliyun-media-censor/util"
	"net/http"
	"net/url"
)

type Pipeline struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	State        string `json:"State"`
	Speed        string `json:"Speed"`
	NotifyConfig struct {
		Topic string `json:"Topic"`
	} `json:"NotifyConfig"`
	Role string `json:"Role"`
}

type ResponseInfo struct {
	RequestID string   `json:"RequestId"`
	Pipeline  Pipeline `json:"Pipeline"`
}

type ResponseListInfo struct {
	RequestID    string `json:"RequestId"`
	TotalCount   int    `json:"TotalCount"`
	PageNumber   int    `json:"PageNumber"`
	PageSize     int    `json:"PageSize"`
	PipelineList struct {
		Pipeline []Pipeline `json:"Pipeline"`
	} `json:"PipelineList"`
}

func QueryPipeLine(ids ...string) (*ResponseListInfo, error) {
	if ids == nil {
		return nil, errors.New("no pipeline ids")
	}

	val := url.Values{
		"PipelineIds": ids[:],
	}
	val.Set("Action", "QueryPipelineList")

	resp, err := http.Get(URL(val))
	if err != nil {
		return nil, err
	}

	info := ResponseListInfo{}

	err = util.UnmarshalJSON(resp.Body, &info)
	return &info, err
}
