package errorData

import "errors"

var (
	LOGIN_DATA_ERROR   = errors.New("Login data error!")
	ACCOUNT_EXIT_ERROR = errors.New("Account already exits!")
	SIGN_DATA_ERROR    = errors.New("Sign data error!")
	QUERY_ERROR        = errors.New("Query error!") // 请求参数错误
	NO_ACCOUNT_ERROR   = errors.New("Has no such account!")

	TOKEN_EXPIRED_ERROR       = errors.New("Token is expired!")      // token已过期
	WITHOUT_TOKEN_ERROR       = errors.New("Request without token!") // 请求未携带token
	TOKEN_NOT_VALID_YET_ERROR = errors.New("Token not active yet")
	TOKEN_MALFORMED_ERROR     = errors.New("That's not even a token")
	TOKEN_INVAILD_ERROR       = errors.New("Couldn't handle this token:")

	NO_PUSH_CASE_AUTHORITY_ERROR = errors.New("You have no push authority!") // 没有上传病历的权限
	NO_GET_CASE_AUTH_ERROR       = errors.New("You have no authority to get other case")

	MAN_ATTACK_ERROR = errors.New("There could be man-in-the-middle attacks") // 可能存在中间人攻击
)
