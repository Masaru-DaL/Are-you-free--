package validation

// 日付の文字列をチェックする
func StringDateCheck(dateStr string) bool {
	// dateの文字列の長さ（例: 2023-01-30[ハイフン込みで10文字分])
	lenDateStr := 10
	var result bool

	if len(dateStr) == lenDateStr {
		result = true
	} else {
		result = false
	}

	return result
}
