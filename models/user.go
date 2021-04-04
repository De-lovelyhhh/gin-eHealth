package models

import (
	errorData "e_healthy/pkg/error"
	"log"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID           int `gorm:"primary_key;autoIncrement"`
	Account      string
	Name         string
	Password     string
	Identity     int
	Organization string
	Sex          int
}

func Login(data map[string]interface{}) string {
	user := User{
		Account:  data["account"].(string),
		Password: data["password"].(string),
		Identity: data["identity"].(int),
	}
	result := db.Select("password").Where("account = ?", user.Account).Find(&user)
	if err := result.Error; err != nil {
		log.Println(err)
		return err.Error()
	}

	if result.RowsAffected == 0 || result.RowsAffected > 1 {
		return errorData.LOGIN_DATA_ERROR.Error()
	}

	return user.Password
}

func Signin(data map[string]interface{}) error {
	user := User{
		Account:      data["account"].(string),
		Name:         data["name"].(string),
		Password:     data["password"].(string),
		Identity:     data["identity"].(int),
		Organization: data["organization"].(string),
		Sex:          data["sex"].(int),
	}

	isExit := db.Where(user)
	if isExit.RowsAffected > 0 {
		return errorData.ACCOUNT_EXIT_ERROR
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if user.Identity == 1 {
			err1 := InitCase(user.ID, user.Account, user.Sex)
			err2 := InitMap(user.ID, user.Account, "", user.Name)
			if err1 != nil {
				return err1
			} else if err2 != nil {
				return err2
			}
		}

		return nil
	})
}

func CheckPatient(patientAccount string) (int, error) {
	user := User{
		Account:  patientAccount,
		Identity: 1,
	}
	result := db.Where("account = ? and identity = 1", user.Account).Find(&user)

	if result.RowsAffected != 1 {
		return errorData.NO_ACCOUNT, errorData.NO_ACCOUNT_ERROR
	}
	if result.Error != nil {
		return errorData.UNKNOWN_ERROR, result.Error
	}

	return errorData.SUCCESS, nil
}

func GetUserInfo(userAccount string) *User {
	user := User{
		Account: userAccount,
	}
	res := db.Where("account = ?", user.Account).Find(&user)

	if res.Error != nil || res.RowsAffected != 1 {
		return nil
	}
	return &user
}
