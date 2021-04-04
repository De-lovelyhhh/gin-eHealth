package models

import (
	errorData "e_healthy/pkg/error"
	"log"

	"github.com/jinzhu/gorm"
)

type CasePrivateMap struct {
	CaseId         int    `json:"case_id";gorm:"primary_key"`
	PrivateId      int    `json:"private_id"`
	PatientAccount string `json:"patient_account"`
	PrivateFields  string `json:"private_fields"`
}

func (CasePrivateMap) TableName() string {
	return "case_privite_map"
}

func UserAuthRequst(query string, queryData string) string {
	queryMap := map[string]string{
		"case_id":         "CaseId",
		"patient_account": "PatientAccount",
	}

	result := db.Where(queryMap[query]+" = ?", queryData)

	if result.Error != nil {
		log.Println(result.Error)
		return result.Error.Error()
	}

	if result.RowsAffected != 1 {
		return errorData.QUERY_ERROR.Error()
	}
	return ""
}

func InitMap(id int, patientAccount string, privateFields string, name string) error {
	if privateFields == "" {
		privateFields = "work;name;patient_account;"
	}
	privateId := id*4 + 3
	privateMap := CasePrivateMap{
		CaseId:         id,
		PatientAccount: patientAccount,
		PrivateFields:  privateFields,
		PrivateId:      privateId,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&privateMap).Error; err != nil {
			return err
		}

		if err := InitPrivateCase(privateMap.PrivateId, name); err != nil {
			return err
		}
		return nil
	})
}

func GetPrivateId(patientAccount string) *CasePrivateMap {
	cpm := &CasePrivateMap{
		PatientAccount: patientAccount,
	}
	res := db.Where("patient_account = ?", cpm.PatientAccount).Find(&cpm)
	if res.RowsAffected != 1 {
		return nil
	}
	return cpm
}

func GetPrivateByCaseId(caseId int) *CasePrivateMap {
	cpm := &CasePrivateMap{
		CaseId: caseId,
	}
	res := db.Where("case_id = ?", caseId).Find(&cpm)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil
	}
	return cpm
}

func GetPrivateField(caseId string) string {
	cpm := &CasePrivateMap{}
	if res := db.Where("case_id = ?", caseId).Find(&cpm); res.Error != nil {
		log.Println(res.Error)
		return ""
	}
	return cpm.PrivateFields
}

func SetPrivateFields(fields string, caseId int) error {
	cpm := &CasePrivateMap{}
	db.Where("case_id = ?", caseId).Find(&cpm)
	cpm.PrivateFields = fields
	res := db.Model(cpm).Where("case_id = ?", caseId).Update(cpm)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
