package api

import (
	"e_healthy/middleware/jwt"
	errorData "e_healthy/pkg/error"
	"e_healthy/service/basic/account_service"
	"net/http"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

type AccountDataForm struct {
	Account      string `binding:"required"`
	Password     string `binding:"required"`
	Identity     int
	Organization string
	Name         string
	Sex          int
}

// LoginResult 登录结果结构
type LoginResult struct {
	Token string `json:"token"`
	User  AccountDataForm
}

func Login(c *gin.Context) {
	var form AccountDataForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, errorData.INVALID_PARAMS)
		return
	}

	accountService := account_service.User{
		Account:  form.Account,
		Password: form.Password,
		Identity: form.Identity,
	}

	statusCode, errorCode := accountService.Login()

	if errorCode == 0 {
		token := generateToken(c, form)
		c.JSON(statusCode, gin.H{
			"code":  errorCode,
			"token": token,
		})
		return
	}

	c.JSON(statusCode, gin.H{
		"code": errorCode,
	})
}

func Signin(c *gin.Context) {
	var form AccountDataForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.INVALID_PARAMS,
		})
		return
	}

	accountService := account_service.User{
		Account:      form.Account,
		Password:     form.Password,
		Identity:     form.Identity,
		Organization: form.Organization,
		Name:         form.Name,
		Sex:          form.Sex,
	}

	statusCode, errorCode := accountService.Signin()
	if errorCode == 0 {
		token := generateToken(c, form)
		c.JSON(statusCode, gin.H{
			"code":  errorCode,
			"token": token,
		})
		return
	}
	c.JSON(statusCode, gin.H{
		"code": errorCode,
	})
}

// 生成令牌
func generateToken(c *gin.Context, user AccountDataForm) string {
	j := &jwt.JWT{
		[]byte("newtrekWang"),
	}
	claims := jwt.CustomClaims{
		user.Identity,
		user.Account,
		user.Password,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),    // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600*24), // 过期时间 一小时
			Issuer:    "newtrekWang",                      //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		return ""
	}

	return token
}

func CheckToken(c *gin.Context) {
	tokenString := c.Request.Header.Get("token")
	if tokenString == "" {
		return
	}
	j := &jwt.JWT{[]byte("newtrekWang")}

	token, err := jwtgo.ParseWithClaims(tokenString, &jwt.CustomClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwtgo.ValidationError); ok {
			if ve.Errors&jwtgo.ValidationErrorMalformed != 0 {
				c.JSON(http.StatusOK, gin.H{
					"code":    errorData.TOKEN_MALFORMED,
					"message": errorData.TOKEN_MALFORMED_ERROR.Error(),
				})
				return
			} else if ve.Errors&jwtgo.ValidationErrorExpired != 0 {
				// Token is expired
				c.JSON(http.StatusOK, gin.H{
					"code":    errorData.TOKEN_EXPIRED,
					"message": errorData.TOKEN_EXPIRED_ERROR.Error(),
				})
				return
			} else if ve.Errors&jwtgo.ValidationErrorNotValidYet != 0 {
				c.JSON(http.StatusOK, gin.H{
					"code":    errorData.TOKEN_NOT_VALID_YET,
					"message": errorData.TOKEN_NOT_VALID_YET_ERROR.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    errorData.TOKEN_INVAILD,
					"message": errorData.TOKEN_INVAILD_ERROR.Error(),
				})
				return
			}
		}
	}
	if claims, ok := token.Claims.(*jwt.CustomClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.SUCCESS,
			"data": claims,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": errorData.UNKNOWN_ERROR,
	})
}

func CheckPatient(c *gin.Context) {
	patientAccount := c.Query("patient_account")
	code, err := account_service.CheckPatient(patientAccount)

	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
		})
	}
}

func GetUserInfo(c *gin.Context) {
	userAccount := c.Query("user_account")
	data := account_service.GetUserInfo(userAccount)

	if data == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.UNKNOWN_ERROR,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": errorData.SUCCESS,
			"data": gin.H{
				"name":         data.Name,
				"identity":     data.Identity,
				"organization": data.Organization,
			},
		})
	}
}
