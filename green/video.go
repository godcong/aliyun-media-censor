package green

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/json-iterator/go"
	"log"
	"net/url"
	"strings"
)

// ContentTypeJSON ...
const ContentTypeJSON = "application/json"

// AliSite ...
const AliSite = "https://green.cn-shanghai.aliyuncs.com"

// URL ...
func URL(values url.Values) string {
	return AliSite + "?" + values.Encode()
}

// Link ...
func Link(url string) string {
	if strings.Index(url, "/") == 0 {
		return AliSite + url
	}

	return AliSite + "/" + url
}

// VideoSyncScanWithCallback ...
func VideoSyncScanWithCallback(data *BizData, fn func(response *green.VideoSyncScanResponse, err error)) <-chan int {
	req := green.CreateVideoSyncScanRequest()
	req.Content = []byte(data.JSON())
	return DefaultClient.VideoSyncScanWithCallback(req, fn)

}

// VideoSyncScan ...
func VideoSyncScan(data *BizData) (*ResultData, error) {
	req := green.CreateVideoSyncScanRequest()
	req.Content = []byte(data.JSON())
	resp, err := DefaultClient.VideoSyncScan(req)
	if err != nil {
		return &ResultData{}, err
	}
	return ResponseToResultData(resp)
}

// VideoAsyncScan ...
func VideoAsyncScan(data *BizData) (*ResultData, error) {
	req := green.CreateVideoAsyncScanRequest()
	req.Content = []byte(data.JSON())
	resp, err := DefaultClient.VideoAsyncScan(req)
	if err != nil {
		return &ResultData{}, err
	}
	return ResponseToResultData(resp)
}

// VideoResults ...
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

// QueueProcessFrame ...
func ProcessFrame(output chan<- *ResultData, files []string, offset int, dest string, objectKey string) {
	var err error

	defer func() {
		if err != nil {
			output <- nil
			return
		}
	}()

	server := oss.Server()
	p := oss.NewProgress()
	var frames []Frame
	for _, file := range files {
		p.SetObjectKey(objectKey + "/" + file)
		p.SetPath(dest)
		if !server.IsExist(p) {
			err := server.Upload(p)
			if err != nil {
				return
			}
		}
		u, err := server.URL(p)
		if err != nil {
			return
		}

		frames = append(frames, Frame{
			URL:    u,
			Offset: offset,
		})
		offset += 15
	}

	data := &BizData{
		Scenes: []string{"porn", "terrorism", "ad", "live", "sface"},
		Tasks: []Task{
			{
				Frames: frames,
			}},
	}

	res, err := VideoSyncScan(data)
	if err != nil {
		return
	}

	output <- res
}

// Process ...
func Process(tasks []Task) (*ResultData, error) {
	data, err := VideoSyncScan(&BizData{
		Scenes:      []string{"porn", "terrorism", "ad", "live", "sface"},
		AudioScenes: []string{"antispam"},
		Tasks:       tasks,
	})
	return data, err
}
