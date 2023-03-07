package router

import (
	"go_http/api"
	"go_http/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	baseGroup := r.Group("/my_chat")
	//登录和注册接口
	baseGroup.POST("login", middleware.SHAMiddleWare(), api.LoginHandler)
	baseGroup.POST("register", middleware.SHAMiddleWare(), api.RegisterHandler)
	//好友相关接口
	baseGroup.POST("friend", middleware.JWTMiddleWare(), api.AddFriend)
	baseGroup.GET("friend", middleware.JWTMiddleWare(), api.QueryList)
	baseGroup.DELETE("friend", middleware.JWTMiddleWare(), api.DeleteFriend)
	//用户信息相关接口
	baseGroup.GET("user_info", api.GetUserInfoByName)
	baseGroup.GET("user_status", api.CheckUserState)

	//websocket接口
	baseGroup.GET("ws", middleware.JWTMiddleWare(), api.WebSocketRegister)
	return r
}
