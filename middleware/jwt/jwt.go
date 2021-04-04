package jwt

import (
	errorData "e_healthy/pkg/error"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": errorData.WITHOUT_TOKEN,
				"msg":    errorData.WITHOUT_TOKEN_ERROR.Error(),
			})
			c.Abort()
			return
		}

		log.Print("get token: ", token)

		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == errorData.TOKEN_EXPIRED_ERROR {
				c.JSON(http.StatusOK, gin.H{
					"status": errorData.TOKEN_EXPIRED,
					"msg":    errorData.TOKEN_EXPIRED_ERROR.Error(),
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": errorData.UNKNOWN_ERROR,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	SignKey string = "newtrekWang"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	Identity int
	Account  string
	Password string
	jwt.StandardClaims
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errorData.TOKEN_MALFORMED_ERROR
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errorData.TOKEN_EXPIRED_ERROR
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errorData.TOKEN_NOT_VALID_YET_ERROR
			} else {
				return nil, errorData.TOKEN_INVAILD_ERROR
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errorData.TOKEN_INVAILD_ERROR
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", errorData.TOKEN_INVAILD_ERROR
}
