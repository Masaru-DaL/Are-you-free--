package num

import (
	"fmt"
	"strconv"
)

func NumToFormattedString(num int) string {
	var str string
	if num < 10 {
		str = fmt.Sprintf("0%d", num)
	} else {
		str = strconv.Itoa(num)
	}

	return str
}
