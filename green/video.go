package green

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	"github.com/json-iterator/go"
	"log"
	"net/url"
	"strings"
)

const ContentTypeJSON = "application/json"

const AliSite = "https://green.cn-shanghai.aliyuncs.com"

func URL(values url.Values) string {
	return AliSite + "?" + values.Encode()
}

func Link(url string) string {
	if strings.Index(url, "/") == 0 {
		return AliSite + url
	}

	return AliSite + "/" + url
}

func VideoSyncScanWithCallback(data *BizData, fn func(response *green.VideoSyncScanResponse, err error)) <-chan int {
	req := green.CreateVideoSyncScanRequest()
	req.Content = []byte(data.JSON())
	return DefaultClient.VideoSyncScanWithCallback(req, fn)

}

func VideoSyncScan(data *BizData) (*ResultData, error) {
	req := green.CreateVideoSyncScanRequest()
	req.Content = []byte(data.JSON())
	resp, err := DefaultClient.VideoSyncScan(req)
	if err != nil {
		return &ResultData{}, err
	}
	return ResponseToResultData(resp)
}

func VideoAsyncScan(data *BizData) (*ResultData, error) {
	req := green.CreateVideoAsyncScanRequest()
	req.Content = []byte(data.JSON())
	resp, err := DefaultClient.VideoAsyncScan(req)
	if err != nil {
		return &ResultData{}, err
	}
	return ResponseToResultData(resp)
}

func VideoResults(request ...string) (*ResultData, error) {
	data, err := jsoniter.Marshal(request[:])
	if err != nil {
		log.Println(data, err)
	}

	req := green.CreateVideoAsyncScanResultsRequest()
	req.Content = data
	resp, err := DefaultClient.VideoAsyncScanResults(req)

	if err != nil {
		return &ResultData{}, err
	}
	return ResponseToResultData(resp)
}
