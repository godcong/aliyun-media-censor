package censor

import (
	"github.com/godcong/aliyun-media-censor/util"
	"net/http"
	"net/url"
)

type ActionResponse struct {
	RequestID string `json:"RequestId"`
	JobID     string `json:"JobId"`
}

const AliSite = "http://mts.cn-hangzhou.aliyuncs.com"

func CensorRequest(values url.Values) (*ActionResponse, error) {
	resp, err := http.Get(URL(values))
	if err != nil {
		return nil, err
	}

	r := ActionResponse{}
	err = util.UnmarshalJSON(resp.Body, &r)
	return &r, err

}

func URL(values url.Values) string {
	return AliSite + "?" + values.Encode()
}
