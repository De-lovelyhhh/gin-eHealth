package setting

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

var DatabaseSetting = &Database{"mysql", "root", "Jh1515451232", "127.0.0.1", "3306", "ehealth"}

var (
	IP_ADDR           = "127.0.0.1"
	REDIRECT_AUTH_URL = IP_ADDR + ":7002/req_auth_code"
	CASE_FIELD        = []string{"case_id",
		"name",
		"age",
		"from",
		"sex",
		"present_illness",
		"illness",
		"doctor_account",
		"patient_account",
		"work",
		"family_history",
		"record_date",
		"diagnosis",
		"allergy",
		"recorder",
	}
)
