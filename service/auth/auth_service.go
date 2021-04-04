package auth

import (
	"e_healthy/middleware/jwt"
	"e_healthy/models"
	"e_healthy/pkg/util"
	client_service "e_healthy/service/client"
	"encoding/base64"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

func GetAuthCode(userAccount string) string {
	code := util.RandStringBytes(10)

	ac := &models.AuthCodeMap{
		Code:        code,
		UserAccount: userAccount,
	}
	if err := models.SetCode(ac); err != nil {
		return ""
	} else {
		return code
	}
}

func CheckState(state string, userAccount string) error {
	psm := &models.PatientStateMap{
		State:          state,
		PatientAccount: userAccount,
	}
	err := models.CheckState(psm)
	return err
}

type AuthData struct {
	Account           string
	Code              string
	ClientCertificate string
	ReqAccount        string
	Identity          int
}

func ToClientServer(ad *AuthData) (bool, string) {
	clientBrand := base64.StdEncoding.EncodeToString([]byte(ad.Account))
	isEqual := client_service.CheckCertificate(clientBrand) == ad.ClientCertificate

	if !isEqual {
		return isEqual, ""
	}
	// 返回token
	token := generateToken(ad.Account, ad.Identity)

	// 绑定医生和私密病历
	privateData := models.GetPrivateId(ad.Account)
	if privateData == nil {
		return true, ""
	} else {
		cdm := &models.CaseDoctorMap{
			PrivateCaseId: privateData.PrivateId,
			UserAccount:   ad.ReqAccount,
			UserIdentity:  ad.Identity,
			CaseId:        privateData.CaseId,
		}
		if err := models.SetCaseDoctorMap(cdm); err != nil {
			return true, ""
		}
	}

	return true, token
}

func generateToken(account string, identity int) string {
	j := &jwt.JWT{
		[]byte("newtrekWang"),
	}
	claims := jwt.CustomClaims{
		identity,
		account,
		"",
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
