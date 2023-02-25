package times

import (
	"src/internal/entity"
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

/* 送信された時間情報が正しいかチェックする */
func CheckInputTime(startFreeTimeHour int, startFreeTimeMinute int, endFreeTimeHour int, endFreeTimeMinute int) bool {
	if startFreeTimeHour > endFreeTimeHour {
		return false
	} else if startFreeTimeHour == endFreeTimeHour {
		if startFreeTimeMinute > endFreeTimeMinute {
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
				break
			}
		}
		// 終了時刻が既存の開始時刻以下
		if endFreeTimeHour <= ft.StartHour {
			// 終了時刻（分）が既存の終了時刻（分）以下の場合
			if endFreeTimeMinute <= ft.StartMinute {
				result = true
				break
			}
		}

		result = false
		break
	}

	return result
}

// 入力された年月日が現在の年月日以降かチェックする
func IsAfterCurrentTime(dateTime string) bool {
	var timeFormat = "2006-01-02"

	parsedTime, _ := time.Parse(timeFormat, dateTime)
	timeNow := time.Now()
	timeTokyo, _ := time.LoadLocation("Asia/Tokyo")
	currentTimeTokyo := timeNow.In(timeTokyo)

	dateData := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, time.Local)
	currentDateData := time.Date(currentTimeTokyo.Year(), currentTimeTokyo.Month(), currentTimeTokyo.Day(), 0, 0, 0, 0, time.Local)

	timeDiff := dateData.Sub(currentDateData)

	if timeDiff == 0 {
		return true
	} else if timeDiff > 0 {
		return true
	} else {
		return false
	}
}
