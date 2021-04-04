package errorData

const (
	SUCCESS        = 0
	INVALID_PARAMS = 1

	UNKNOWN_ERROR = 100000
	ACCOUNT_ERROR = 100001
	ACCOUNT_EXITS = 100002
	NO_ACCOUNT    = 100003

	TOKEN_EXPIRED       = 100100
	WITHOUT_TOKEN       = 100101
	TOKEN_NOT_VALID_YET = 100102
	TOKEN_MALFORMED     = 100103
	TOKEN_INVAILD       = 100104

	NO_PUSH_CASE_AUTHORITY = 100201
	NO_GET_CASE_AUTH       = 100202

	MAN_ATTACK = 200001
)
