package entity

import "errors"

/* Errors for admin */
var (
	ErrInternalServerError = errors.New("internal Server Error")

	ErrNoUserFound           = errors.New("no user found")
	ErrSQLGetFailed          = errors.New("failed to get data from table")
	ErrSQLCreateStmt         = errors.New("failed to create prepared statement")
	ErrSQLExecFailed         = errors.New("failed to exec sql")
	ErrSQLResultNotDesired   = errors.New("sql result is not desired")
	ErrSQLLastInsertIdFailed = errors.New("failed to get last inserted id")
)

/* Errors for user */
const (
	ERR_INTERNAL_SERVER_ERROR = "もう一度入力してください。"
	ERR_INPUT_EMPTY           = "値を入力してください。"
	ERR_NO_WHITESPACE         = "空白は混入させないでください。"
	ERR_NO_QUOTATION          = "クォーテーションは含めないでください。"
	ERR_INPUT_MISSING         = "再確認用の値が一致しません。"
	ERR_FAILED_EMAIL_FORMAT   = "Emailの形式が正しくありません。"
	ERR_ALREADY_USER_EXISTS   = "入力された情報は既に存在しています。"
	ERR_NO_USER               = "入力された情報が存在しません。"
	ERR_INPUT_LIMIT_OVER      = "制限文字数以内で入力してください。"
)
