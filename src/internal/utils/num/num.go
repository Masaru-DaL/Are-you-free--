package num

import (
	"fmt"
	"strconv"
)

/* 整数型で10以下の場合は先頭に0を付けて文字列で返す */
func NumToFormattedString(num int) string {
	var str string
	if num < 10 {
		str = fmt.Sprintf("0%d", num)
	} else {
		str = strconv.Itoa(num)
	}

	return str
}
