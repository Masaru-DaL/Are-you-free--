package strings

import "regexp"

/* 文字列に空白が混入されているかのチェックを行う */
func CheckWhitespaceInString(stringData string) bool {
	// 特殊文字/空白などの区切り文字の指定と初期化
	reg := "[ 　]"
	initializedReg := regexp.MustCompile(reg)

	// 文字列のチェック
	checkResult := initializedReg.MatchString(stringData)

	return checkResult
}
