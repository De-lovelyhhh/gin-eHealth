package api

import (
	"e_healthy/handler/wshandler"
	"e_healthy/middleware/jwt"
	"e_healthy/models"
	errorData "e_healthy/pkg/error"
	"e_healthy/service/auth"
	"e_healthy/service/basic/basic_auth_service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserAuthRequst(c *gin.Context) {
	var (
		query     string
		queryData string
	)

	claims := c.MustGet("claims").(*jwt.CustomClaims)
	if claims == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.UNKNOWN_ERROR,
		})
		return
	}

	if c.Query("patient_account") != "" {
		query = "patient_account"
		queryData = c.Query("patient_account")
	} else if c.Query("case_id") != "" {
		query = "case_id"
		queryData = c.Query("case_id")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.INVALID_PARAMS,
		})
		return
	}

	basic_auth_service.UserAuthRequst(query, queryData, claims.Account, claims.Identity)
}

func PatientAgreeAuth(c *gin.Context) {
	agreeCode, _ := strconv.Atoi(c.Query("agree_code"))
	agreeData := &basic_auth_service.AgreeData{
		PatientAccount: c.Query("patient_account"),
		SenderAccount:  c.Query("sender_account"),
		AgreeCode:      agreeCode,
	}

	if agreeCode == 1 {
		state, redirect_url := basic_auth_service.PatientAgreeAuth()
		// 生成state之后存起来
		psm := &models.PatientStateMap{
			PatientAccount: agreeData.PatientAccount,
			State:          state,
		}
		models.SetMap(psm)
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.SUCCESS,
			"data": gin.H{
				"state":        state,
				"redirect_url": redirect_url,
			},
		})
		return
	}
	// 拒绝授权
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
	})
}

func CheckVoucher(c *gin.Context) {}

func ToAuth(c *gin.Context) {
	identity, _ := strconv.Atoi(c.Query("identity"))
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	ad := &auth.AuthData{
		Account:           claims.Account,
		Code:              c.Query("code"),
		ClientCertificate: c.Query("clientCertificate"),
		ReqAccount:        c.Query("reqAccount"),
		Identity:          identity,
	}

	isEqual, token := auth.ToClientServer(ad)
	if !isEqual {
		c.JSON(http.StatusOK, gin.H{
			"code":    errorData.MAN_ATTACK,
			"message": errorData.MAN_ATTACK_ERROR.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.SUCCESS,
		})

		rm := map[string]string{
			"target": ad.ReqAccount,
			"token":  token,
		}
		rmString, _ := json.Marshal(rm)
		wshandler.Manager.Broadcast <- rmString
	}
}
