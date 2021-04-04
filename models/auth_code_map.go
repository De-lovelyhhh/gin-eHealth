package models

import "github.com/jinzhu/gorm"

type AuthCodeMap struct {
	UserAccount string `json:"user_account"`
	Code        string
}

func (AuthCodeMap) TableName() string {
	return "auth_code_map"
}

func SetCode(ac *AuthCodeMap) error {
	return db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_account = ?", ac.UserAccount).Find(&ac)
		if res.RowsAffected != 0 {
			if err := tx.Model(ac).Update(ac).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Create(ac).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func DeleteCode() {}
