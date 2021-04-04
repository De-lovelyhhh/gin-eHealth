package router

import (
	"e_healthy/api"
	// "e_healthy/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func InitWebsocket(e *gin.RouterGroup) {
	WebsocketRouter := e.Group("/ws")
	// WebsocketRouter.Use(jwt.JWTAuth())

	{
		WebsocketRouter.GET("/patient_auth", api.WsAuthApi)
	}
}
