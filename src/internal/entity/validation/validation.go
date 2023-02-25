package validation

import "regexp"

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
