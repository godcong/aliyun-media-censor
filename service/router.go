package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"log"
	"net/http"
)

// Router ...
func Router(eng *gin.Engine) {
	verV0 := "v0"
	eng.Use(AccessControlAllow)
	g0 := eng.Group(verV0)
	//登录
	g0.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	g0.GET("list", func(ctx *gin.Context) {
		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey("5F6688B5.mp4")
		err := server.Upload(p)
		if err != nil {
			log.Println(err)
		}
	})
	g0.GET("url", func(ctx *gin.Context) {
		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey("5F6688B5.mp4")

		u, err := server.URL(p)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		log.Println(u)
		success(ctx, u)
	})

	g0.GET("validate", func(ctx *gin.Context) {
		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey("5F6688B5.mp4")

		u, err := server.URL(p)
		data, err := green.GreenVideoAsyncscan(&green.VideoRequest{
			Scenes:      []string{"porn", "terrorism", "ad", "live", "sface"},
			AudioScenes: []string{"antispam"},
			Tasks: []green.Task{
				{
					DataID:    "dataid 00001",
					URL:       u,
					Interval:  1,
					MaxFrames: 200,
				},
			},
		})
		log.Println(data, err)

		success(ctx, data)
	})

	g0.GET("status/:id", func(ctx *gin.Context) {

	})

}
