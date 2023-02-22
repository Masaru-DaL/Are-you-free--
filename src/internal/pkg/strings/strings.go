package strings

import (
	"regexp"
	"strings"
)

/* 文字列に空白が混入されているかのチェックを行う */
func CheckWhitespaceInString(stringData string) bool {
	// 空白の指定と初期化
	reg := "[ 　]"
	initializedReg := regexp.MustCompile(reg)

	// 文字列のチェック
	checkResult := initializedReg.MatchString(stringData)

	return checkResult
}

/* 文字列にクオーテーションが混入されているかのチェックを行う */
func CheckQuotationInString(stringData string) bool {
	// クオーテーションの指定と初期化
	reg := "[\"'`”’｀]"
	initializedReg := regexp.MustCompile(reg)

	// 文字列のチェック
	checkResult := initializedReg.MatchString(stringData)

	return checkResult
}

/* Emailの形式のチェックを行う */
func CheckEmailFormat(emailAddress string) bool {
	// email形式の指定と初期化
	emailFormat := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	initializedEmailFormat := regexp.MustCompile(emailFormat)

	// emailのチェック
	checkResult := initializedEmailFormat.MatchString(emailAddress)
	return checkResult
}

/* 日付文字列を年/月/日で分割する */
func SplitDateByHyphen(dateString string) (string, string, string) {
	dateArray := strings.Split(dateString, "-")

	yearStr := dateArray[0]
	monthStr := dateArray[1]
	dayStr := dateArray[2]

	return yearStr, monthStr, dayStr
}
