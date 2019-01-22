package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/godcong/aliyun-media-censor/green"
	"github.com/godcong/aliyun-media-censor/oss"
	"github.com/godcong/role-manager-server/util"
	"github.com/satori/go.uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"time"
)

// LoginPOST ...
func LoginPOST(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.PostForm("username")
	}
}

func result(ctx *gin.Context, code int, message string, detail interface{}) {
	h := gin.H{
		"code":    code,
		"message": message,
		"detail":  detail,
	}
	ctx.JSON(http.StatusOK, h)
}

func success(ctx *gin.Context, detail interface{}) {
	result(ctx, 0, "success", detail)
}

func failed(ctx *gin.Context, message string) {
	result(ctx, -1, message, nil)
}

// AccessControlAllow ...
func AccessControlAllow(ctx *gin.Context) {
	origin := ctx.Request.Header.Get("origin")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, XMLHttpRequest, "+
		"Accept-Encoding, X-CSRF-Token, Authorization")
	if ctx.Request.Method == "OPTIONS" {
		ctx.String(200, "ok")
		return
	}
	ctx.Next()
}

// EncryptJWT ...
func EncryptJWT(key []byte, sub []byte) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		panic(err)
	}
	cl := jwt.Claims{
		Subject:   string(sub),
		Issuer:    "godcong",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 14)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        util.GenerateRandomString(16),
	}

	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	return raw, err
}

// ValidatePOST ...
func ValidatePOST(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		qi := QueueInfo{
			ObjectKey: ctx.PostForm("object_key"),
			//CallbackURL:  ctx.PostForm("url"),
			ID:           ctx.PostForm("id"),
			ValidateType: ctx.PostForm("validate_type"),
		}
		var data []*green.ResultData
		var err error
		if qi.ValidateType == "frame" {
			Push(&qi)
			data = []*green.ResultData{}
		} else if qi.ValidateType == "jpg" {
			data, err = ParseValidateDo(&qi, func(url string) (data *green.ResultData, e error) {
				return green.ImageSyncScan(&green.BizData{
					Scenes: []string{"porn", "terrorism", "ad", "live", "sface"},
					Tasks: []green.Task{
						{
							DataID: uuid.NewV1().String(),
							URL:    url,
						},
					},
				})
			})
		} else {
			data, err = ParseValidateDo(&qi, func(url string) (data *green.ResultData, e error) {
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

		}

		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, data)

	}
}

// ParseValidateDo ...
func ParseValidateDo(info *QueueInfo, fn func(url string) (*green.ResultData, error)) ([]*green.ResultData, error) {
	server := oss.Server()
	p := oss.NewProgress()

	var dataList []*green.ResultData

	p.SetObjectKey(info.ObjectKey)

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

	return dataList, nil
}
