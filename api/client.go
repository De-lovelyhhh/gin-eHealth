package api

import (
	"crypto/md5"
	errorData "e_healthy/pkg/error"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveClientBrand(c *gin.Context) {
	clientBrand := c.Query("clientBrand")
	h := md5.New()
	h.Write([]byte(clientBrand))
	ClientCertificate := hex.EncodeToString(h.Sum(nil))

	c.JSON(http.StatusOK, gin.H{
		"code": errorData.SUCCESS,
		"data": gin.H{
			"clientCertificate": ClientCertificate,
		},
	})
}
