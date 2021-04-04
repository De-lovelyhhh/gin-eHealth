package basic_auth_service

import (
	"e_healthy/models"
	errorData "e_healthy/pkg/error"
	"e_healthy/pkg/setting"
	"e_healthy/pkg/util"
)

type AgreeData struct {
	PatientAccount string
	SenderAccount  string
	AgreeCode      int
}

func UserAuthRequst(query string, queryData string, userAccount string, userIdentity int) int {
	err := models.UserAuthRequst(query, queryData)

	if err == errorData.QUERY_ERROR.Error() {
		return errorData.INVALID_PARAMS
	} else if err != "" {
		return errorData.UNKNOWN_ERROR
	}

	return errorData.SUCCESS
}

func PatientAgreeAuth() (string, string) {
	return util.RandStringBytes(15), setting.REDIRECT_AUTH_URL
}
