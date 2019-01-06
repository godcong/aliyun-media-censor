package censor

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ActionResponse struct {
	RequestID string `json:"RequestId"`
	JobID     string `json:"JobId"`
}

const AliSite = "http://mts.cn-hangzhou.aliyuncs.com"

func CensorRequest(values url.Values) (*ActionResponse, error) {
	link := AliSite + "?" + values.Encode()
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := ActionResponse{}
	err = jsoniter.Unmarshal(bytes, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil

}
