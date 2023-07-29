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

/* 日付文字列を年/月/日で分割する */
func SplitDateByHyphen(dateString string) (string, string, string) {
	dateArray := strings.Split(dateString, "-")

	yearStr := dateArray[0]
	monthStr := dateArray[1]
	dayStr := dateArray[2]

	return yearStr, monthStr, dayStr
}

// 日付の文字列をチェックする
func IsDateString(dateStr string) bool {
	// dateの文字列の長さ（例: 2023-01-30[ハイフン込みで10文字分])
	lenDateStr := 10

	if len(dateStr) == lenDateStr {
		return true
	} else {
		return false
	}
}

// 年月日が別々に送られた来た場合の文字列をチェックする
func IsYearString(dateStr string) bool {
	lenYearStr := 4
	yearRegex := regexp.MustCompile(`^[0-9]+$`)

	if len(dateStr) == lenYearStr {
		return yearRegex.MatchString(dateStr)
	} else {
		return false
	}
}
func IsMonthDayString(dateStr string) bool {
	lenMonthAndDayStr := 2
	monthAndDayRegex := regexp.MustCompile(`^[0-9]+$`)

	if len(dateStr) == lenMonthAndDayStr {
		return monthAndDayRegex.MatchString(dateStr)
	} else {
		return false
	}
}

// 時間が送られて来た場合の文字列をチェックする
func IsTimeString(dateStr string) bool {
	lenTimeStr := 2
	timeRegex := regexp.MustCompile(`^[0-9]+$`)

	if len(dateStr) == lenTimeStr {
		return timeRegex.MatchString(dateStr)
	} else {
		return false
	}
}
