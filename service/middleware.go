package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/role-manager-server/model"
	"github.com/json-iterator/go"
)

// LoginCheck ...
func LoginCheck(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token == "" {
			failed(ctx, "token is null")
			ctx.Abort()
			return
		}
		t := Token{}
		e := jsoniter.Unmarshal([]byte(token), &t)
		if e != nil {
			failed(ctx, e.Error())
			ctx.Abort()
			return
		}

		user := model.NewUser()
		user.ID = model.ID(t.OID)
		err := user.Find()
		if err != nil {
			failed(ctx, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
