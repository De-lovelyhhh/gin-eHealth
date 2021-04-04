package models

import "github.com/jinzhu/gorm"

type CaseDoctorMap struct {
	PrivateCaseId int    `json:"private_case_id"`
	CaseId        int    `json:"case_id"`
	UserAccount   string `json:"user_account"`
	UserIdentity  int    `json:"user_identity"`
}

func (CaseDoctorMap) TableName() string {
	return "case_doctor_map"
}

func SetCaseDoctorMap(cdm *CaseDoctorMap) error {
	cdmm := &CaseDoctorMap{
		PrivateCaseId: cdm.PrivateCaseId,
		UserAccount:   cdm.UserAccount,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("private_case_id = ? and user_account = ?", cdmm.PrivateCaseId, cdmm.UserAccount).Find(&cdmm)
		if res.RowsAffected != 0 {
			return nil
		} else {
			if err := tx.Create(cdm).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func FindRecord(pci int, ua string) bool {
	cdm := &CaseDoctorMap{
		PrivateCaseId: pci,
		UserAccount:   ua,
	}
	res := db.Where("private_case_id = ? and user_account = ?", cdm.PrivateCaseId, cdm.UserAccount).Find(&cdm)
	if res.Error == nil && res.RowsAffected == 1 {
		return true
	}
	return false
}

func FindByDoctor(userAccount string) []*CaseDoctorMap {
	cdm := &CaseDoctorMap{
		UserAccount: userAccount,
	}
	res := make([]*CaseDoctorMap, 0)

	db.Where("user_account = ?", cdm.UserAccount).Find(&res)
	return res
}
