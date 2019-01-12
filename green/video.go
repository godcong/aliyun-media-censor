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

func QueueProcessJPG(output chan<- string, info *oss.QueueInfo) {
	server := oss.Server2()
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
		count := 1
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

		data := &BizData{
			Tasks: []Task{
				{
					ClientInfo:  nil,
					DataID:      "",
					URL:         "",
					Frames:      frames,
					FramePrefix: "",
					Interval:    0,
					MaxFrames:   0,
				}},
		}

		resultData, err := VideoSyncScan(data)

		resp, err := http.PostForm(info.CallbackURL, url.Values{
			"request_key": []string{info.RequestKey},
			"code":        []string{"-1"},
			"msg":         []string{err.Error()},
			"detail":      []string{string(resultData.JSON())},
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

func Process(tasks []Task) (*ResultData, error) {
	data, err := VideoSyncScan(&BizData{
		Scenes:      []string{"porn", "terrorism", "ad", "live", "sface"},
		AudioScenes: []string{"antispam"},
		Tasks:       tasks,
	})
	return data, err
}
