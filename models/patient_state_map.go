package models

import "github.com/jinzhu/gorm"

type PatientStateMap struct {
	PatientAccount string `json:"patient_account"`
	State          string
}

func (PatientStateMap) TableName() string {
	return "patient_state_map"
}

func SetMap(psm *PatientStateMap) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := db.Model(psm).Update(psm).Error; err != nil {
			return err
		}
		return nil
	})
}

func CheckState(psm *PatientStateMap) error {
	res := db.Where(psm)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
