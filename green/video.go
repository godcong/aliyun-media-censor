package green

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	"github.com/json-iterator/go"
	"log"
	"net/url"
	"strings"
)

const ContentTypeJSON = "application/json"
const GMTTimeFormmat = "Mon, 02 Jan 2006 15:04:05 -0700 GMT"

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

func VideoAsyncScan(data *BizData) (string, error) {
	req := green.CreateVideoAsyncScanRequest()
	req.Content = []byte(data.JSON())
	response, err := DefaultClient.VideoAsyncScan(req)
	b := response.GetHttpContentBytes()
	var d ResultData
	err = jsoniter.Unmarshal(b, &d)
	fmt.Printf("%+v", d)
	log.Println(response.String())
	return "", err
}

func VideoResults(request ...string) (string, error) {
	data, err := jsoniter.Marshal(request[:])
	if err != nil {
		log.Println(data, err)
	}

	req := green.CreateVideoAsyncScanResultsRequest()
	req.Content = data
	response, err := DefaultClient.VideoAsyncScanResults(req)
	b := response.GetHttpContentBytes()
	var d ResultData
	err = jsoniter.Unmarshal(b, &d)
	fmt.Printf("%+v", d)
	log.Println(response.String())
	return "", err
}

//
//func (r *VideoRequest) Reader() io.Reader {
//	b, err := jsoniter.Marshal(r)
//	if err != nil {
//		log.Println(err)
//		return nil
//	}
//	return bytes.NewBuffer(b)
//
//}
