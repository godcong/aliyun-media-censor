package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/ffmpeg"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/godcong/aliyun-media-censor/util"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// Router ...
func Router(eng *gin.Engine) {
	verV0 := "v0"
	eng.Use(AccessControlAllow, func(ctx *gin.Context) {
		log.Println("visit", ctx.Request.URL.String())
	})
	g0 := eng.Group(verV0)
	//登录
	g0.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	g0.POST("download", func(ctx *gin.Context) {
		server := oss.Server()
		p := oss.NewProgress()

		key := ctx.PostForm("objectKey")
		p.SetObjectKey(key)
		if !server.IsExist(p) {
			failed(ctx, "obejct key is not exist")
			return
		}
		err := server.Download(p)
		if err != nil {
			log.Println(err)
			failed(ctx, err.Error())
			return
		}
		success(ctx, key)
	})

	g0.GET("list/:path", func(ctx *gin.Context) {
		path := ctx.Param("path")
		path = strings.Replace(path, ".", "/", -1)

		files, err := util.FileList(path)
		if err != nil {
			log.Println(err)
			failed(ctx, err.Error())
			return
		}
		success(ctx, files)
	})

	g0.GET("url", func(ctx *gin.Context) {
		server := oss.Server()
		key := ctx.Query("key")
		p := oss.NewProgress()
		p.SetObjectKey(key)

		if !server.IsExist(p) {
			failed(ctx, "object is not exist")
			return

		}
		url, err := server.URL(p)
		if !server.IsExist(p) {
			failed(ctx, err.Error())
			return
		}
		success(ctx, url)
	})

	g0.POST("fileupload", func(ctx *gin.Context) {
		filePath := ctx.PostForm("name")

		server := oss.Server()
		p := oss.NewProgress()
		p.SetObjectKey(filePath)

		if !server.IsExist(p) {
			err := server.Upload(p)
			if err != nil {
				log.Println(err)
				failed(ctx, err.Error())
				return
			}
		}
		url, err := server.URL(p)
		if err != nil {
			log.Println(err)
			failed(ctx, err.Error())
			return
		}
		success(ctx, url)

	})

	g0.POST("upload", func(ctx *gin.Context) {
		filePath := ctx.PostForm("name")
		tp := ctx.PostForm("type")

		server := oss.Server()
		p := oss.NewProgress()
		var urls []string
		files := []string{filepath.Join("./download", filePath)}
		if tp == "pic" {
			p.SetObjectKey(filePath)
			if !server.IsExist(p) {
				err := server.Upload(p)
				if err != nil {
					log.Println(err)
					failed(ctx, err.Error())
					return
				}
			}
			u, err := server.URL(p)
			if err != nil {
				log.Println(err)
				failed(ctx, err.Error())
				return
			}
			urls = append(urls, u)
			success(ctx, urls)
			return
		}

		ts, err := ffmpeg.TransferSplit("./download/", filePath)
		if err != nil {
			log.Println(err)
			failed(ctx, err.Error())
			return
		}
		log.Println(ts)

		files, err = util.FileList("transferred/" + filePath)
		if err != nil {
			log.Println(err)
			failed(ctx, err.Error())
			return
		}

		for _, file := range files {
			p.SetObjectKey(filePath + "/" + file)
			p.SetPath("transferred")
			if !server.IsExist(p) {
				err := server.Upload(p)
				if err != nil {
					log.Println(err)
					failed(ctx, err.Error())
					return
				}
			}
			u, err := server.URL(p)
			if err != nil {
				log.Println(err)
				failed(ctx, err.Error())
				return
			}
			urls = append(urls, u)
		}

		success(ctx, urls)
	})
	g0.POST("validate", ValidatePOST(verV0))

	g0.POST("validate/frame", func(ctx *gin.Context) {
		failed(ctx, "please use /validate")
	})

	g0.POST("validate/pic", func(ctx *gin.Context) {
		failed(ctx, "please use /validate")
	})

	g0.GET("status/pic", func(ctx *gin.Context) {
		id := ctx.QueryArray("id")
		data, err := green.ImageAsyncResult(id...)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, data)
	})

	g0.POST("validate/video", func(ctx *gin.Context) {
		failed(ctx, "please use /validate")
	})

	g0.GET("status", func(ctx *gin.Context) {
		id := ctx.QueryArray("id")
		tp := ctx.Query("type")

		if tp == "video" {
			data, err := green.VideoResults(id...)
			if err != nil {
				failed(ctx, err.Error())
				return
			}
			success(ctx, data)
		}

	})

}
