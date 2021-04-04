package router

import (
	"e_healthy/api"
	"e_healthy/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func InitCaseRouter(e *gin.RouterGroup) {
	CaseRouter := e.Group("/case")
	CaseRouter.Use(jwt.JWTAuth())
	{
		CaseRouter.POST("/push", api.PushCase)                      // 上传病历
		CaseRouter.GET("/case_list", api.GetCaseList)               // 获取病历列表
		CaseRouter.GET("/case_detail", api.GetCaseDetail)           // 获取病历详情（脱敏后）
		CaseRouter.GET("/my_case_list", api.GetMyCaseList)          // 获取已授权的病历
		CaseRouter.GET("/change_field", api.ChangeCasePrivateField) // 自定义敏感字段
		CaseRouter.GET("/get_id", api.GetCaseId)                    // 通过患者ID获取病历ID
		CaseRouter.GET("/get_fields", api.GetFields)                // 获取敏感字段
		CaseRouter.GET("/get_account_by_caseid", api.GetAccountByCase)
	}
}
