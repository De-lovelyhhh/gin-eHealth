package router

import (
	"e_healthy/api"

	"github.com/gin-gonic/gin"
)

func InitAccountRouter(e *gin.RouterGroup) {
	AccountRouter := e.Group("/account")
	{
		AccountRouter.POST("/login", api.Login)             // 登录
		AccountRouter.POST("/signin", api.Signin)           // 注册
		AccountRouter.POST("/check_token", api.CheckToken)  // 检查token合法性
		AccountRouter.GET("/check_user", api.CheckPatient)  // 检查患者身份
		AccountRouter.GET("/get_userinfo", api.GetUserInfo) // 获取用户信息
	}
}
