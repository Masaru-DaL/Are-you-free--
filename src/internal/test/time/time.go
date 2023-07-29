package time

import (
	"src/internal/entity"
	"strconv"
	"time"
)

/* 年月日の情報から曜日を取得する */
func GetWeekdayByDate(yearStr, monthStr, dayStr string) string {
	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)
	day, _ := strconv.Atoi(dayStr)

	// 日本語返却用
	jpWeekdayArray := []string{"日", "月", "火", "水", "木", "金", "土"}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	jpWeekday := jpWeekdayArray[date.Weekday()]

	return jpWeekday
}

/* 送信された時間情報が正しいかチェックする */
func CheckInputTime(startFreeTimeHour int, startFreeTimeMinute int, endFreeTimeHour int, endFreeTimeMinute int) bool {
	if startFreeTimeHour > endFreeTimeHour {
		return false
	} else if startFreeTimeHour == endFreeTimeHour {
		if startFreeTimeMinute > endFreeTimeMinute {
			return false
		} else if startFreeTimeMinute == endFreeTimeMinute {
			return false
		}
	}

	return true
}

/* free-timeを作成するために送られた時間情報が既に存在しているfree-timeと被っていないかをチェックする */
func IsCreateFreeTime(startFreeTimeHour int, startFreeTimeMinute int, endFreeTimeHour int, endFreeTimeMinute int, dateFreeTime *entity.DateFreeTime) bool {
	var result bool

	for _, ft := range dateFreeTime.FreeTimes {
		// true
		// 開始時刻が既存の終了時刻以上
		if startFreeTimeHour >= ft.EndHour {
			// 開始時刻（分）が既存の開始時刻（分）以上の場合
			if startFreeTimeMinute >= ft.EndMinute {
				result = true
				continue
			}
		}
		// 終了時刻が既存の開始時刻以下
		if endFreeTimeHour <= ft.StartHour {
			// 終了時刻（分）が既存の終了時刻（分）以下の場合
			if endFreeTimeMinute <= ft.StartMinute {
				result = true
				continue
			}
		}

		result = false
		break
	}

	return result
}
