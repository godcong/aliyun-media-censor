package main

import (
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/satori/go.uuid"
	"log"
)

func main() {
	server := oss.Server2()
	p := oss.NewProgress()
	p.SetObjectKey("videorandom.mp4")
	//err := server.Upload(p)
	//if err != nil {
	//	log.Println(err)
	//}
	u, err := server.URL(p)

	data, err := green.VideoSyncScan(&green.BizData{
		Scenes:      []string{"porn", "terrorism", "ad", "live", "sface"},
		AudioScenes: []string{"antispam"},
		Tasks: []green.Task{
			{
				DataID:    uuid.NewV1().String(),
				URL:       u,
				Interval:  1,
				MaxFrames: 200,
			},
		},
	})
	log.Println(data, err)
}
