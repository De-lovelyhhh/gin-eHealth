package api

import (
	errorData "e_healthy/pkg/error"
	"net/http"

	"e_healthy/service/auth"

	"github.com/gin-gonic/gin"
)

func GetAuthCode(c *gin.Context) {
	state := c.Query("state")
	userAccount := c.Query("patientAccount")

	// 检查state
	if err := auth.CheckState(state, userAccount); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    errorData.UNKNOWN_ERROR,
			"message": err.Error(),
		})
		return
	}

	// 生成授权码
	code := auth.GetAuthCode(userAccount)
	if code != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.SUCCESS,
			"data": gin.H{
				"code": code,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.UNKNOWN_ERROR,
		})
	}

}
