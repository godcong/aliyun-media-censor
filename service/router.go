package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/role-manager-server/model"
	"github.com/godcong/role-manager-server/util"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Router ...
func Router(eng *gin.Engine) {
	verV0 := "v0"
	eng.Use(AccessControlAllow)
	g0 := eng.Group(verV0)
	//登录
	g0.POST("login", LoginPOST(verV0))
	//组织注册
	g0.POST("register", RegisterPOST(verV0))

	g0.POST("genesis", GenesisGet(verV0))

	v0 := g0.Group("")
	v0.Use(LoginCheck(verV0))
	v0.POST("addUser", AddUserPOST(verV0))

}

// GenesisGet ...
func GenesisGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := model.NewGenesis()
		if model.FindOne(role, bson.M{
			"slug": role.Slug,
		}) != nil && role.ID != primitive.NilObjectID {
			failed(ctx, "genesis is created")
			return
		}
		passwd := util.GenerateRandomString(16)
		user := model.NewUser()
		user.Name = "genesis"
		user.SetPassword(passwd)
		err := model.Transaction(func() error {

			err := role.Create()
			if err != nil {
				return err
			}
			err = user.Create()
			if err != nil {
				return err
			}
			ru := model.NewRoleUser()
			ru.SetUser(user)
			ru.SetRole(role)
			err = ru.CreateIfNotExist()
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, gin.H{
			"Name":     user.Name,
			"Password": passwd,
		})
		return
	}
}

// RegisterPOST ...
func RegisterPOST(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.PostForm("applyName")     //商户名称
		ctx.PostForm("applyCode")     //社会统一信用代码
		ctx.PostForm("applyContact")  //商户联系人
		ctx.PostForm("applyPosition") //联系人职位
		ctx.PostForm("applyPhone")    //联系人手机号
		ctx.PostForm("applyMailbox")  //联系人邮箱
	}
}

// AddUserPOST ...
func AddUserPOST(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// User ...
func User(ctx *gin.Context) *model.User {
	if v, b := ctx.Get("user"); b {
		if v0, b := v.(*model.User); b {
			return v0
		}
	}
	return nil
}
