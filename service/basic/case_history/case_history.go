package case_history

import (
	"e_healthy/models"
	"e_healthy/pkg/util"

	// "encoding/base64"
	"fmt"
	"reflect"
	"strconv"
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

	if identity == 1 && res.PatientAccount == account {
		return models.GetCaseDetail(caseId), models.GetPrivateData(res.PrivateId)
	} else if identity != 1 {
		isTrue := models.FindRecord(res.PrivateId, account)
		if isTrue {
			return models.GetCaseDetail(caseId), models.GetPrivateData(res.PrivateId)
		}
	}

	return models.GetCaseDetail(caseId), nil
}

type CompleteCase struct {
	CaseHistory        *models.CaseHistory
	PrivateCaseHistory *models.PrivateCaseHistory
}

func GetMyCaseList(userAccount string, identity int) []*CompleteCase {
	ch := make([]*CompleteCase, 0)
	privateList := models.FindByDoctor(userAccount)
	for _, v := range privateList {
		newCh, newPch := GetCaseDetail(v.CaseId, userAccount, identity)
		ch = append(ch, &CompleteCase{
			CaseHistory:        newCh,
			PrivateCaseHistory: newPch,
		})
	}
	return ch
}

func GetPrivateField(caseId string) []string {
	return strings.Split(models.GetPrivateField(caseId), ";")
}

func PushCase(ch *models.CaseHistory, pch *models.PrivateCaseHistory, pchm map[string]interface{}) {
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

// 更新私密病历时联动更新两个表
func UpdateCaseField(priId int, caseId int, deleteData []string, newData []string) error {
	// var aeskey = []byte("321423u9y8d2fwfl")
	ch := models.GetCaseDetail(caseId)
	vCh := reflect.ValueOf(ch).Elem()
	chm := make(map[string]interface{}, 0)

	pch := models.GetPrivateData(priId)
	vPch := reflect.ValueOf(pch).Elem()
	pchm := make(map[string]interface{}, 0)

	for _, val := range newData {
		camel := util.Camel(val)
		var dataStr string
		var dataInt int
		switch vCh.FieldByName(camel).Type().String() {
		case "int":
			dataInt = int(vCh.FieldByName(camel).Int())
			data := strconv.Itoa(dataInt)
			pchm[camel] = data
			vPch.FieldByName(camel).SetString(data)
			chm[camel] = 0
		case "string":
			dataStr = vCh.FieldByName(camel).String()
			pchm[camel] = dataStr
			vPch.FieldByName(camel).SetString(dataStr)
			chm[camel] = ""
		}
	}

	for _, val := range deleteData {
		camel := util.Camel(val)
		if data := vPch.FieldByName(camel).String(); data != "" {
			fmt.Print(data)
			switch vCh.FieldByName(camel).Type().String() {
			case "int":
				chm[camel] = toInt(data)
			case "string":
				chm[camel] = data
			}
			pchm[camel] = ""
			vPch.FieldByName(camel).SetString("")
			// outPriData[camel] = string(tpass)
		}
	}

	if err := models.SetEncrypt(pch); err != nil {
		return err
	}
	models.PushCaseMap(chm, ch)
	return nil
}

func toInt(str string) int64 {
	res, _ := strconv.ParseInt(str, 10, 64)
	return res
}
