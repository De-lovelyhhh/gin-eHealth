package api

import (
	"e_healthy/middleware/jwt"
	"e_healthy/models"
	errorData "e_healthy/pkg/error"
	"e_healthy/pkg/setting"
	"e_healthy/pkg/util"
	"e_healthy/service/basic/case_history"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkClaim(c *jwt.CustomClaims, g *gin.Context) {
	if c == nil {
		g.JSON(http.StatusOK, gin.H{
			"code": errorData.UNKNOWN_ERROR,
		})
		return
	}
}

func PushCase(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	checkClaim(claims, c)

	if claims.Identity != 0 && claims.Identity != 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    errorData.NO_PUSH_CASE_AUTHORITY,
			"message": errorData.NO_PUSH_CASE_AUTHORITY_ERROR.Error(),
		})
		return
	}
	// 从数据库中获取敏感字段
	fields := case_history.GetPrivateField(c.PostForm("case_id"))
	// 敏感数据结构体
	pch := &models.PrivateCaseHistory{}
	// 利用反射获取私密病历表中的字段
	v := reflect.ValueOf(pch).Elem()
	for i, val := range fields {
		// 遍历敏感字段
		if i == len(fields)-1 {
			break
		}
		camel := camel(val)
		// 将敏感数据存入结构体中，以便后面存入数据库
		v.FieldByName(camel).SetString(c.PostForm(val))
	}

	// 普通病历结构体
	ch := &models.CaseHistory{}
	chv := reflect.ValueOf(ch).Elem()
	// 利用差集计算得到非敏感字段
	diffSet := util.SliceDiffStr(setting.CASE_FIELD, fields)
	for _, val := range diffSet {
		// 遍历非敏感字段
		camel := camel(val)
		var postData string
		if c.PostForm(val) == "--" {
			postData = ""
		} else {
			postData = c.PostForm(val)
		}
		// 根据字段类型赋值
		switch chv.FieldByName(camel).Type().String() {
		case "int":
			chv.FieldByName(camel).SetInt(toInt(postData))
		case "string":
			chv.FieldByName(camel).SetString(postData)
		}
	}
	case_history.PushCase(ch, pch)
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
	})
}

func toInt(str string) int64 {
	res, _ := strconv.ParseInt(str, 10, 64)
	return res
}

func camel(str string) string {
	name := strings.Replace(str, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func GetCaseList(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	checkClaim(claims, c)
	res := case_history.GetCaseList()
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
		"data": res,
	})
}

func GetCaseDetail(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	checkClaim(claims, c)
	caseId, _ := strconv.Atoi(c.Query("caseId"))
	caseData, priData := case_history.GetCaseDetail(caseId, claims.Account, claims.Identity)
	if caseData != nil && priData != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.SUCCESS,
			"data": gin.H{
				"case":    caseData,
				"private": priData,
			},
		})
	} else if caseData != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.SUCCESS,
			"data": gin.H{
				"case": caseData,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.UNKNOWN_ERROR,
		})
	}
}

func GetMyCaseList(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	checkClaim(claims, c)
	if claims.Identity == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.NO_GET_CASE_AUTH,
			"data": errorData.NO_GET_CASE_AUTH_ERROR.Error(),
		})
	}
	// 先找case_doctor_map表，再去一个个查
	myCaseList := case_history.GetMyCaseList(claims.Account, claims.Identity)
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
		"data": myCaseList,
	})
}

func GetCaseId(c *gin.Context) {
	account := c.Query("account")
	id := case_history.GetCaseIdByAccount(account)
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
		"data": id,
	})
}

func GetFields(c *gin.Context) {
	fields := case_history.GetPrivateField(c.Query("case_id"))
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
		"data": fields[:len(fields)-1],
	})
}

func ChangeCasePrivateField(c *gin.Context) {
	fields := c.Query("fields")
	caseId, _ := strconv.Atoi(c.Query("caseId"))
	err := case_history.SetPrivateFields(fields, caseId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.UNKNOWN_ERROR,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
	})
}

func GetAccountByCase(c *gin.Context) {
	caseId, _ := strconv.Atoi(c.Query("case_id"))
	account := case_history.GetModel(caseId).PatientAccount
	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
		"data": account,
	})
}
