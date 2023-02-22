package validation

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

	if len(dateStr) == lenYearStr {
		return true
	} else {
		return false
	}
}
func IsMonthAndDayString(dateStr string) bool {
	lenMonthAndDayStr := 2

	if len(dateStr) == lenMonthAndDayStr {
		return true
	} else {
		return false
	}
}
