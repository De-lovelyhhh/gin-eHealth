package models

import (
	"encoding/base64"
	"reflect"

	"e_healthy/pkg/util"

	"github.com/jinzhu/gorm"
)

type PrivateCaseHistory struct {
	PrivateCaseId  int `json:"private_case_id";grom:"primaryKey"`
	Work           string
	Sex            string
	Age            string
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

func (PrivateCaseHistory) TableName() string {
	return "private_case"
}

func InitPrivateCase(caseId int, name string) error {
	var aeskey = []byte("321423u9y8d2fwfl")
	aesName := base64.StdEncoding.EncodeToString(util.AesEncrypt([]byte(name), aeskey))
	privateCaseHistory := PrivateCaseHistory{
		PrivateCaseId: caseId,
		Name:          aesName,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&privateCaseHistory).Error; err != nil {
			return err
		}

		return nil
	})
}

func SetEncrypt(pch *PrivateCaseHistory) error {
	// 加密秘钥
	var aeskey = []byte("321423u9y8d2fwfl")
	t := reflect.TypeOf(pch).Elem()
	v := reflect.ValueOf(pch).Elem()
	for k := 0; k < v.NumField(); k++ {
		// 遍历私密病历结构体中的成员
		if t.Field(k).Name == "PrivateCaseId" {
			continue
		}
		var fieldsData string
		if v.Field(k).String() == "" {
			fieldsData = ""
		} else {
			// 加密
			var pass []byte
			if v.Field(k).String() == "--" {
				pass = []byte("")
			} else {
				pass = []byte(v.Field(k).String())
			}
			xpass := util.AesEncrypt(pass, aeskey)
			fieldsData = base64.StdEncoding.EncodeToString(xpass)
		}

		v.Field(k).SetString(fieldsData)
	}
	// 根据私密病历ID更新数据表中的记录
	res := db.Model(pch).Where("private_case_id = ?", pch.PrivateCaseId).Update(pch)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func GetPrivateData(privateId int) *PrivateCaseHistory {
	pch := &PrivateCaseHistory{
		PrivateCaseId: privateId,
	}

	res := db.Where("private_case_id = ?", pch.PrivateCaseId).Find(&pch)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil
	}
	var aeskey = []byte("321423u9y8d2fwfl")
	t := reflect.TypeOf(pch).Elem()
	v := reflect.ValueOf(pch).Elem()
	for k := 0; k < v.NumField(); k++ {
		if t.Field(k).Name == "PrivateCaseId" {
			continue
		}
		if v.Field(k).String() == "" {
			continue
		}
		bytesPass, _ := base64.StdEncoding.DecodeString(v.Field(k).String())
		tpass := util.AesDecrypt(bytesPass, aeskey)
		// if string(tpass) == "" {
		// 	continue
		// }
		v.Field(k).SetString(string(tpass))
		// fmt.Printf("%s -- %v \n", t.Filed(k).Name, v.Field(k).Interface())
	}
	return pch
}
