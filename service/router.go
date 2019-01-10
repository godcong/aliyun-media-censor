package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/satori/go.uuid"
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

	g0.GET("validatepic", func(ctx *gin.Context) {
		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey("img20120822103826045350.jpg")
		err := server.Upload(p)
		if err != nil {
			log.Println(err)
		}
		u, err := server.URL(p)

		data, err := green.ImageAsyncScan(&green.BizData{
			Scenes: []string{"porn"},
			Tasks: []green.Task{
				{
					DataID: uuid.NewV1().String(),
					URL:    u,
				},
			},
		})
		if err != nil {
			failed(ctx, err.Error())
		}

		success(ctx, data)

	})

	g0.GET("statuspic/:id", func(ctx *gin.Context) {
		data, err := green.ImageAsyncResult(ctx.Param("id"))
		log.Println(data, err)
	})

	g0.GET("validate", func(ctx *gin.Context) {
		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey("5F6688B5.mp4")

		u, err := server.URL(p)

		data, err := green.VideoAsyncScan(&green.BizData{
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

		success(ctx, data)
	})

	g0.GET("status/:id", func(ctx *gin.Context) {
		data, err := green.VideoResults(ctx.Param("id"))
		log.Println(data, err)
	})

}
