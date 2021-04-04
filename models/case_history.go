package models

import (
	"github.com/jinzhu/gorm"
)

type CaseHistory struct {
	CaseId         int `json:"case_id";grom:"primaryKey"`
	Work           string
	Sex            int
	Age            int
	From           string
	RecordDate     string `json:"record_date"`
	Allergy        string
	Recorder       string
	PresentIllness string `json:"present_illness"`
	Illness        string
	DoctorAccount  string `json:"doctor_account"`
	FamilyHistory  string `json:"family_history"`
	Diagnosis      string
	Name           string
	PatientAccount string `json:"patient_account"`
}

func (CaseHistory) TableName() string {
	return "case_history"
}

func InitCase(id int, patientAccount string, sex int) error {
	user := CaseHistory{
		CaseId:         id,
		PatientAccount: patientAccount,
		Sex:            sex,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetAllCase() (ch []CaseHistory) {
	db.Find(&ch)
	return ch
}

func GetCaseDetail(ch *CaseHistory) *CaseHistory {
	res := db.Where("case_id = ?", ch.CaseId).Find(&ch)
	if res.Error != nil {
		return nil
	}
	return ch
}

func PushCaseHistory(ch *CaseHistory) error {
	res := db.Model(ch).Where("case_id = ?", ch.CaseId).Update(ch)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
