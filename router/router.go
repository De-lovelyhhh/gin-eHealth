package router

import (
	"e_healthy/api"
	"e_healthy/middleware/jwt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterBasic() *gin.Engine {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Type", "token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	}))

	BasicGroup := e.Group("")
	{
		InitAccountRouter(BasicGroup) // 登录注册相关api
		InitCaseAuth(BasicGroup)      // oauth认证相关api
		InitWebsocket(BasicGroup)     // websocket服务
		InitCaseRouter(BasicGroup)    // 病历相关api
	}

	return e
}

func RouterAuth() *gin.Engine {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Type", "token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	}))

	AuthGroup := e.Group("")
	AuthGroup.Use(jwt.JWTAuth())
	{
		AuthGroup.GET("/req_auth_code", api.GetAuthCode) // 获取授权码
	}

	return e
}

func RouterClient() *gin.Engine {
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Type", "token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	}))

	ClientGroup := e.Group("")
	ClientGroup.Use(jwt.JWTAuth())
	{
		ClientGroup.GET("/receive_client_brand", api.ReceiveClientBrand) // 接收客户端记号并生成客户端凭证
	}

	return e
}
