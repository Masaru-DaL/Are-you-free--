package entity

import "errors"

/* Errors for admin */
var (
	ErrInternalServerError = errors.New("internal Server Error")

	ErrNoUserFound           = errors.New("no user found")
	ErrNoDateFreeTimeFound   = errors.New("no date-free-time found")
	ErrNoFreeTimeFound       = errors.New("no free-time found")
	ErrSQLGetFailed          = errors.New("failed to get data from table")
	ErrSQLCreateStmt         = errors.New("failed to create prepared statement")
	ErrSQLExecFailed         = errors.New("failed to exec sql")
	ErrSQLResultNotDesired   = errors.New("sql result is not desired")
	ErrSQLLastInsertIdFailed = errors.New("failed to get last inserted id")
	ErrSQLTransactionError   = errors.New("sql transaction error")
)

/* Errors for user */
var (
	MESSAGE_INTERNAL_SERVER_ERROR = "エラーが発生しました。"
	MESSAGE_NO_FREE_TIME_FOUND    = "作成したfree-timeがありません。"
)

/* ERROR CODE */
const (
	ERR_INTERNAL_SERVER_ERROR    = "ERR_INTERNAL_SERVER_ERROR"
	ERR_INPUT_EMPTY              = "値を入力してください。"
	ERR_NO_WHITESPACE            = "空白は混入させないでください。"
	ERR_NO_QUOTATION             = "クォーテーションは含めないでください。"
	ERR_NO_CHOICE                = "選択されていない箇所があります。"
	ERR_INPUT_MISSING            = "再確認用の値が一致しません。"
	ERR_FAILED_EMAIL_FORMAT      = "Emailの形式が正しくありません。"
	ERR_ALREADY_USER_EXISTS      = "入力された情報は既に存在しています。"
	ERR_NO_USER                  = "入力された情報が存在しません。"
	ERR_INPUT_LIMIT_OVER         = "制限文字数以内で入力してください。"
	ERR_CHOICE_TIME              = "作成するfree-timeの指定に誤りがあります。"
	ERR_CHOICE_DATE              = "現在の日付より前に作成することはできません。"
	ERR_ALREADY_FREE_TIME_EXISTS = "設定された時刻の範囲に既にfree-timeが存在します。"
	ERR_NO_DATE_FREE_TIME_FOUND  = "ERR_NO_DATE_FREE_TIME_FOUND"
)
