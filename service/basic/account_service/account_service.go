package account_service

import (
	"e_healthy/models"
	errorData "e_healthy/pkg/error"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Account      string
	Name         string
	Password     string
	Identity     int
	Organization string
	Sex          int
}

func (u *User) Login() (int, int) {
	user := map[string]interface{}{
		"account":  u.Account,
		"password": u.Password,
		"identity": u.Identity,
	}

	hashPassword := models.Login(user)
	isRight := comparePasswords(hashPassword, []byte(u.Password))

	if isRight {
		return http.StatusOK, errorData.SUCCESS
	} else {
		return http.StatusOK, errorData.UNKNOWN_ERROR
	}
}

func (u *User) Signin() (int, int) {
	user := map[string]interface{}{
		"account":      u.Account,
		"password":     PasswordEncryption([]byte(u.Password)),
		"identity":     u.Identity,
		"name":         u.Name,
		"organization": u.Organization,
		"sex":          u.Sex,
	}

	if err := models.Signin(user); err != nil && err == errorData.ACCOUNT_EXIT_ERROR {
		return http.StatusOK, errorData.ACCOUNT_EXITS
	} else if err != nil {
		return http.StatusOK, errorData.UNKNOWN_ERROR
	}

	return http.StatusOK, errorData.SUCCESS
}

func PasswordEncryption(pwd []byte) string {
	res, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		return ""
	}

	return string(res)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CheckPatient(patientAccount string) (int, string) {
	code, err := models.CheckPatient(patientAccount)
	if err != nil {
		return code, err.Error()
	}
	return code, ""
}

func GetUserInfo(userAccount string) *models.User {
	data := models.GetUserInfo(userAccount)
	return data
}
