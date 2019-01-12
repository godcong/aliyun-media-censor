package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/ffmpeg"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/godcong/aliyun-media-censor/util"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"path/filepath"
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

	g0.POST("download", func(ctx *gin.Context) {
		server := oss.Server2()
		p := oss.NewProgress()

		name := ctx.PostForm("name")

		p.SetObjectKey(name)
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
		success(ctx, p.ObjectKey())
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

	g0.POST("upload", func(ctx *gin.Context) {
		filePath := ctx.PostForm("name")
		tp := ctx.PostForm("type")

		server := oss.Server2()
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

	g0.POST("validate/frame", func(ctx *gin.Context) {

	})

	g0.POST("validate/pic", func(ctx *gin.Context) {
		data, err := ParseValidateDo(ctx, func(url string) (data *green.ResultData, e error) {
			return green.ImageAsyncScan(&green.BizData{
				Scenes: []string{"porn"},
				Tasks: []green.Task{
					{
						DataID: uuid.NewV1().String(),
						URL:    url,
					},
				},
			})
		})

		if err != nil {
			failed(ctx, err.Error())
		}

		success(ctx, data)

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
		data, err := ParseValidateDo(ctx, func(url string) (data *green.ResultData, e error) {
			return green.VideoAsyncScan(&green.BizData{
				Scenes:      []string{"porn", "terrorism", "ad", "live", "sface"},
				AudioScenes: []string{"antispam"},
				Tasks: []green.Task{
					{
						DataID:    uuid.NewV1().String(),
						URL:       url,
						Interval:  30,
						MaxFrames: 200,
					},
				}})
		})

		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, data)
	})

	g0.GET("status/video", func(ctx *gin.Context) {
		id := ctx.QueryArray("id")
		data, err := green.VideoResults(id...)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, data)
	})

}

// ParseValidateDo ...
func ParseValidateDo(ctx *gin.Context, fn func(url string) (*green.ResultData, error)) ([]*green.ResultData, error) {
	server := oss.Server2()
	p := oss.NewProgress()

	tp := ctx.PostForm("type")
	names := ctx.PostFormArray("names")
	name := ctx.PostForm("name")

	if tp != "list" {
		names = []string{name}
	}

	var dataList []*green.ResultData

	for _, name := range names {
		p.SetObjectKey(name)

		if !server.IsExist(p) {
			return nil, errors.New("obejct key is not exist")
		}

		u, err := server.URL(p)
		if err != nil {
			return nil, err
		}

		resultData, err := fn(u)

		if err != nil {
			return nil, err
		}
		dataList = append(dataList, resultData)
	}
	return dataList, nil
}
