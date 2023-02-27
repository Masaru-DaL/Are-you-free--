package freetimes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"src/internal/entity"
	"src/internal/infra/config"
	"src/internal/repository"
	"src/internal/test/time"
	"src/internal/utils/strings"
	"src/internal/utils/times"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CreateFreeTime(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)

		var year int
		var month int
		var day int
		var yearStr string
		var monthStr string
		var dayStr string

		/* 入力されたデータの処理 */
		yearStr, monthStr, dayStr = c.FormValue("year"), c.FormValue("month"), c.FormValue("day")
		// 日付が事前に入力されていた場合
		if yearStr != "" || monthStr != "" || dayStr != "" {
			// 別々に送られてきた日付文字列をチェックする
			isYearStr := strings.IsYearString(yearStr)
			isMonthStr := strings.IsMonthDayString(monthStr)
			isDayStr := strings.IsMonthDayString(dayStr)
			if !isYearStr || !isMonthStr || !isDayStr {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}

			year, _ = strconv.Atoi(yearStr)
			month, _ = strconv.Atoi(monthStr)
			day, _ = strconv.Atoi(dayStr)

			// 日付が事前に入力されていなかった場合
		} else {
			dateStr := c.FormValue("date")
			// 日付文字列が空だった場合
			if dateStr == "" {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_NO_CHOICE,
				})
			}
			// 日付文字列をチェックする
			isDateString := strings.IsDateString(dateStr)
			if !isDateString {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}
			// 入力された日付情報が現在の日付より後かチェックする
			isAfterCurrentTime := times.IsAfterCurrentTime(dateStr)
			if !isAfterCurrentTime {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_CHOICE_DATE,
				})
			}
			yearStr, monthStr, dayStr = strings.SplitDateByHyphen(dateStr)
			year, _ = strconv.Atoi(yearStr)
			month, _ = strconv.Atoi(monthStr)
			day, _ = strconv.Atoi(dayStr)
		}

		// free-time情報を受け取る
		startFreeTimeHourStr, startFreeTimeMinuteStr := c.FormValue("start-free-time-hour"), c.FormValue("start-free-time-minute")
		endFreeTimeHourStr, endFreeTimeMinuteStr := c.FormValue("end-free-time-hour"), c.FormValue("end-free-time-minute")
		// 全ての時間情報を入力されたかどうか
		if startFreeTimeHourStr == "" || startFreeTimeMinuteStr == "" || endFreeTimeHourStr == "" || endFreeTimeMinuteStr == "" {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_NO_CHOICE,
			})
		}
		// 入力された時間文字列をチェックする
		isStartFreeTimeHourStr := strings.IsTimeString(startFreeTimeHourStr)
		isStartFreeTimeMinuteStr := strings.IsTimeString(startFreeTimeMinuteStr)
		isEndFreeTimeHourStr := strings.IsTimeString(endFreeTimeHourStr)
		isEndFreeTimeMinuteStr := strings.IsTimeString(endFreeTimeMinuteStr)
		if !isStartFreeTimeHourStr || !isStartFreeTimeMinuteStr || !isEndFreeTimeHourStr || !isEndFreeTimeMinuteStr {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
			})
		}

		startFreeTimeHour, _ := strconv.Atoi(startFreeTimeHourStr)
		startFreeTimeMinute, _ := strconv.Atoi(startFreeTimeMinuteStr)
		endFreeTimeHour, _ := strconv.Atoi(endFreeTimeHourStr)
		endFreeTimeMinute, _ := strconv.Atoi(endFreeTimeMinuteStr)
		// 入力されたfree-timeが正常に送られたかどうかをチェックする
		checkResult := time.CheckInputTime(startFreeTimeHour, startFreeTimeMinute, endFreeTimeHour, endFreeTimeMinute)
		if !checkResult {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_CHOICE_TIME,
			})
		}

		// 入力された日付のdate-free-timeを取得する
		dateFreeTime, err := repository.GetDateFreeTime(ctx, db, userID, year, month, day)
		// 指定した日付のdate-free-timeが存在しなかった場合
		if errors.Is(err, entity.ErrNoDateFreeTimeFound) {
			// DateFreeTime構造体で値を設定する
			dateFreeTime = &entity.DateFreeTime{
				UserID: userID,
				Year:   year,
				Month:  month,
				Day:    day,
			}
			// DateFreeTimeの作成
			dateFreeTime, err = repository.CreateDateFreeTime(ctx, db, dateFreeTime)
			if err != nil {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}
		} else if err != nil {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
			})
		}

		// 指定した日付のfree-timeが存在した場合
		if dateFreeTime.FreeTimes != nil {
			// free-timeが作成できるかどうか
			checkResult = time.IsCreateFreeTime(startFreeTimeHour, startFreeTimeMinute, endFreeTimeHour, endFreeTimeMinute, dateFreeTime)
			if !checkResult {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_ALREADY_FREE_TIME_EXISTS,
				})
			}
		}

		// free-timeを作成するための値を構造体に格納する
		freeTime := &entity.FreeTime{
			DateFreeTimeID: dateFreeTime.ID,
			StartHour:      startFreeTimeHour,
			StartMinute:    startFreeTimeMinute,
			EndHour:        endFreeTimeHour,
			EndMinute:      endFreeTimeMinute,
		}
		// free-timeの作成
		_, err = repository.CreateFreeTime(ctx, db, freeTime)
		if err != nil {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
			})
		}

		jpWeekday := time.GetWeekdayByDate(year, month, day)

		successMessage := fmt.Sprintf("%s/%s/%s（%s）%s:%s〜%s:%sでfree-timeを作成しました。", yearStr, monthStr, dayStr, jpWeekday, startFreeTimeHourStr, startFreeTimeMinuteStr, endFreeTimeHourStr, endFreeTimeMinuteStr)

		return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
			"year":            nil,
			"month":           nil,
			"day":             nil,
			"weekday":         nil,
			"error_message":   nil,
			"success_message": successMessage,
		})
	}
}
