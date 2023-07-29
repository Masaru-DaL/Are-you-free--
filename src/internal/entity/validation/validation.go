package validation

import "regexp"

/* Emailの形式のチェックを行う */
func IsEmail(emailAddress string) bool {
	// email形式の指定と初期化
	emailFormat := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	initializedEmailFormat := regexp.MustCompile(emailFormat)

	// emailのチェック
	checkResult := initializedEmailFormat.MatchString(emailAddress)
	return checkResult
}
