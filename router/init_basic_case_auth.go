package router

import (
	"e_healthy/api"
	"e_healthy/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func InitCaseAuth(e *gin.RouterGroup) {
	CaseAuthRouter := e.Group("/auth")
	CaseAuthRouter.Use(jwt.JWTAuth())
	{
		CaseAuthRouter.GET("/get_auth", api.UserAuthRequst)     // 用户向业务服务器发起授权请求
		CaseAuthRouter.GET("/agree_auth", api.PatientAgreeAuth) // 资源拥有者确认授权接口
		CaseAuthRouter.GET("/check_voucher", api.CheckVoucher)  // 资源拥有者通过业务服务器验证客户端凭证是否被更改
		CaseAuthRouter.GET("/to_auth", api.ToAuth)              // 去认证服务器验证客户端凭证
	}
}
