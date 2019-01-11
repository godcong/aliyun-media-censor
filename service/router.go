package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/ffmpeg"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	g0.GET("list/:path", func(ctx *gin.Context) {
		path := ctx.Param("path")
		path = strings.Replace(path, ".", "/", -1)
		files, err := ioutil.ReadDir(path)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		var fileNames []string

		for _, file := range files {
			if !file.IsDir() {
				fileNames = append(fileNames, file.Name())
			}
		}

		success(ctx, fileNames)
	})

	g0.POST("upload", func(ctx *gin.Context) {
		filePath := ctx.PostForm("name")
		ts, err := ffmpeg.TransferSplit("./download/" + filePath)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		log.Println(ts)
		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey(filePath)

		if !server.IsExist(p) {
			err := server.Upload(p)
			if err != nil {
				failed(ctx, err.Error())
				return
			}
		}

		u, err := server.URL(p)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		log.Println(u)
		success(ctx, u)
	})

	g0.POST("validate/:name/pic", func(ctx *gin.Context) {
		name := ctx.Param("name")
		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey(name)

		if !server.IsExist(p) {
			failed(ctx, "obejct key is not exist")
			return
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

	g0.GET("status/:id/pic", func(ctx *gin.Context) {
		data, err := green.ImageAsyncResult(ctx.Param("id"))
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, data)
	})

	g0.POST("validate/:name/video", func(ctx *gin.Context) {
		name := ctx.Param("name")

		server := oss.Server2()
		p := oss.NewProgress()
		p.SetObjectKey(name)

		if !server.IsExist(p) {
			failed(ctx, "obejct key is not exist")
			return
		}

		u, err := server.URL(p)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		data, err := green.VideoAsyncScan(&green.BizData{
			Scenes:      []string{"porn", "terrorism", "ad", "live", "sface"},
			AudioScenes: []string{"antispam"},
			Tasks: []green.Task{
				{
					DataID:    uuid.NewV1().String(),
					URL:       u,
					Interval:  30,
					MaxFrames: 200,
				},
			},
		})
		if err != nil {
			failed(ctx, err.Error())
			return
		}

		success(ctx, data)
	})

	g0.GET("status/:id/video", func(ctx *gin.Context) {
		data, err := green.VideoResults(ctx.Param("id"))
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, data)
	})

}
