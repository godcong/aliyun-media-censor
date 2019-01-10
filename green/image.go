package green

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	"github.com/json-iterator/go"
	"log"
)

func ImageAsyncScan(data *BizData) (*ResultData, error) {
	req := green.CreateImageAsyncScanRequest()
	req.Content = []byte(data.JSON())
	response, err := DefaultClient.ImageAsyncScan(req)
	if err != nil {
		return &ResultData{}, err
	}
	var d ResultData
	err = jsoniter.Unmarshal(response.GetHttpContentBytes(), &d)
	fmt.Printf("%+v", d)
	log.Println(response.String())
	return &d, err
}

func ImageAsyncResult(request ...string) (*ResultData, error) {
	req := green.CreateImageAsyncScanResultsRequest()
	bytes, err := jsoniter.Marshal(request[:])
	if err != nil {
		return &ResultData{}, err
	}
	req.Content = bytes
	response, err := DefaultClient.ImageAsyncScanResults(req)
	if err != nil {
		return &ResultData{}, err
	}
	var d ResultData
	err = jsoniter.Unmarshal(response.GetHttpContentBytes(), &d)
	fmt.Printf("%+v", d)
	log.Println(response.String())
	return &d, err
}
