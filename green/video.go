package green

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	"github.com/godcong/aliyun-media-censor/ffmpeg"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/godcong/aliyun-media-censor/util"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
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

// QueueProcessJPG ...
func QueueProcessJPG(output chan<- string, info *oss.QueueInfo) {
	server := oss.Server()
	p := oss.NewProgress()
	p.SetObjectKey(info.ObjectKey)

	if server.IsExist(p) {
		err := server.Download(p)
		ts, err := ffmpeg.TransferJPG("./download/", info.ObjectKey)
		if err != nil {
			log.Println(err)
			output <- err.Error()
			return
		}
		log.Println(ts)
		files, err := util.FileList("transferred/" + info.ObjectKey)
		if err != nil {
			log.Println(err)
			output <- err.Error()
			return
		}

		var frames []Frame
		count := 0
		for _, file := range files {
			p.SetObjectKey(info.ObjectKey + "/" + file)
			p.SetPath("transferred")
			if !server.IsExist(p) {
				err := server.Upload(p)
				if err != nil {
					log.Println(err)
					output <- err.Error()
					return
				}
			}
			u, err := server.URL(p)
			if err != nil {
				log.Println(err)
				output <- err.Error()
				return
			}

			frames = append(frames, Frame{
				URL:    u,
				Offset: count,
			})
			count += 15
		}
		log.Println("request:", info.RequestKey)
		log.Println("frames:", frames)
		data := &BizData{
			Scenes: []string{"porn", "terrorism", "ad", "live", "sface"},
			Tasks: []Task{
				{
					Frames: frames,
				}},
		}

		resultData, err := VideoSyncScan(data)
		msg := "success"
		code := "0"
		if err != nil {
			msg = err.Error()
			code = "-1"
		}

		resp, err := http.PostForm(info.CallbackURL, url.Values{
			"request_key": []string{info.RequestKey},
			"code":        []string{code},
			"message":     []string{msg},
			"detail":      []string{string(resultData.ArrayedJSON())},
		})

		if err != nil {
			log.Println(err)
			output <- err.Error()
			return
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			output <- err.Error()
			return
		}
		log.Println(string(bytes))
	}

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
