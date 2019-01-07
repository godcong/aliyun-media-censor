package censor

import (
	"github.com/godcong/aliyun-media-censor/util"
	"net/http"
	"net/url"
	"time"
)

type JobResponse struct {
	RequestID string `json:"RequestId"`
	JobID     string `json:"JobId"`
}

type JobDetailResponse struct {
	RequestID            string `json:"RequestId"`
	MediaCensorJobDetail struct {
		Input struct {
			Bucket   string `json:"Bucket"`
			Location string `json:"Location"`
			Object   string `json:"Object"`
		} `json:"Input"`
		Suggestion       string `json:"Suggestion"`
		DescCensorResult struct {
			Scene      string `json:"Scene"`
			Suggestion string `json:"Suggestion"`
			Label      string `json:"Label"`
			Rate       string `json:"Rate"`
		} `json:"DescCensorResult"`
		CreationTime      time.Time `json:"CreationTime"`
		State             string    `json:"State"`
		TitleCensorResult struct {
			Scene      string `json:"Scene"`
			Suggestion string `json:"Suggestion"`
			Label      string `json:"Label"`
			Rate       string `json:"Rate"`
		} `json:"TitleCensorResult"`
		PipelineID         string `json:"PipelineId"`
		ID                 string `json:"Id"`
		VensorCensorResult struct {
			NextPageToken      string `json:"NextPageToken"`
			CensorCensorResult struct {
				CensorResult []struct {
					Scene      string `json:"Scene"`
					Suggestion string `json:"Suggestion"`
					Label      string `json:"Label"`
					Rate       string `json:"Rate"`
				} `json:"CensorResult"`
			} `json:"CensorCensorResult"`
			VideoTimelines struct {
				VideoTimeline []struct {
					CensorCensorResult struct {
						CensorResult []struct {
							Scene      string `json:"Scene"`
							Suggestion string `json:"Suggestion"`
							Label      string `json:"Label"`
							Rate       string `json:"Rate"`
						} `json:"CensorResult"`
					} `json:"CensorCensorResult"`
					Timestamp string `json:"Timestamp"`
					Object    string `json:"Object"`
				} `json:"VideoTimeline"`
			} `json:"VideoTimelines"`
		} `json:"VensorCensorResult"`
		BarrageCensorResult struct {
			Scene      string `json:"Scene"`
			Suggestion string `json:"Suggestion"`
			Label      string `json:"Label"`
			Rate       string `json:"Rate"`
		} `json:"BarrageCensorResult"`
		VideoCensorConfig struct {
			VideoCensor bool `json:"VideoCensor"`
			OutputFile  struct {
				Bucket   string `json:"Bucket"`
				Location string `json:"Location"`
				Object   string `json:"Object"`
			} `json:"OutputFile"`
		} `json:"VideoCensorConfig"`
	} `json:"MediaCensorJobDetail"`
}

const AliSite = "http://mts.cn-hangzhou.aliyuncs.com"

func CensorRequest(values url.Values) (*JobResponse, error) {
	resp, err := http.Get(URL(values))
	if err != nil {
		return nil, err
	}

	r := JobResponse{}
	err = util.UnmarshalJSON(resp.Body, &r)
	return &r, err

}

func URL(values url.Values) string {
	return AliSite + "?" + values.Encode()
}
