package case_history

import (
	"e_healthy/models"
	"strings"
)

func GetCaseList() []models.CaseHistory {
	return models.GetAllCase()
}

func GetCaseDetail(caseId int, account string, identity int) (*models.CaseHistory, *models.PrivateCaseHistory) {
	// 1 查case_private_map表，拿私密ID
	// 2 如果不是病患，那么找case_doctor_map
	// 3 如果上面不满足直接返回脱敏后的病历
	res := models.GetPrivateByCaseId(caseId)
	if res == nil {
		return nil, nil
	}
	ch := &models.CaseHistory{
		CaseId: caseId,
	}
	if identity == 1 && res.PatientAccount == account {
		return models.GetCaseDetail(ch), models.GetPrivateData(res.PrivateId)
	} else if identity != 1 {
		isTrue := models.FindRecord(res.PrivateId, account)
		if isTrue {
			return models.GetCaseDetail(ch), models.GetPrivateData(res.PrivateId)
		}
	}

	return models.GetCaseDetail(ch), nil
}

type CompleteCase struct {
	*models.CaseHistory
	*models.PrivateCaseHistory
}

func GetMyCaseList(userAccount string, identity int) []*CompleteCase {
	ch := make([]*CompleteCase, 0)
	privateList := models.FindByDoctor(userAccount)
	for _, v := range privateList {
		newCh, newPch := GetCaseDetail(v.CaseId, userAccount, identity)
		ch = append(ch, &CompleteCase{
			newCh,
			newPch,
		})
	}
	return ch
}

func GetPrivateField(caseId string) []string {
	return strings.Split(models.GetPrivateField(caseId), ";")
}

func PushCase(ch *models.CaseHistory, pch *models.PrivateCaseHistory) {
	models.PushCaseHistory(ch)
	// 获取私密病历ID
	cpm := models.GetPrivateByCaseId(ch.CaseId)
	pch.PrivateCaseId = cpm.PrivateId
	models.SetEncrypt(pch)
}

func GetCaseIdByAccount(account string) int {
	return models.GetPrivateId(account).CaseId
}

func SetPrivateFields(fields string, caseId int) error {
	err := models.SetPrivateFields(fields, caseId)
	if err != nil {
		return err
	}
	return nil
}

func GetModel(caseId int) *models.CasePrivateMap {
	return models.GetPrivateByCaseId(caseId)
}
