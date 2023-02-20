package time

import (
	"time"
)

/* 年月日の情報から曜日を取得する */
func GetWeekdayByDate(year int, month int, day int) string {
	// 日本語返却用
	jpWeekdayArray := []string{"日", "月", "火", "水", "木", "金", "土"}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	jpWeekday := jpWeekdayArray[date.Weekday()]

	return jpWeekday
}
