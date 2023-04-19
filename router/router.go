package router

import (
	"Im/middlewares"
	"Im/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})
	//用户登录
	r.POST("/login", service.Login)
	//发送验证验证码
	r.POST("/send/code", service.SendCode)
	//注册
	r.POST("/register", service.Register)

	//用户组
	auth := r.Group("/u", middlewares.AuthCheck())

	//用户详情
	auth.GET("/user/detail", service.UserDetail)

	//发送接收消息
	auth.GET("/websocket/message", service.WebsocketMessage)

	//聊天消息列表
	auth.GET("/chat/list", service.ChatList)

	//通过账号查询用户
	auth.POST("/search/account", service.SearchByAccount)

	//添加好友
	auth.POST("/add/account", service.AddFriend)
	//添加好友
	auth.POST("/delete/account", service.DeleteFriend)
	return r
}
