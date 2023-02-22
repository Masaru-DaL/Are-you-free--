package freetimes

import (
	"context"
	"fmt"
	"net/http"
	"src/internal/config"
	"src/internal/entity"
	"src/internal/repository/gateway"
	"src/internal/utils/num"
	"src/internal/utils/strings"
	"src/internal/utils/time"
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
			year, _ = strconv.Atoi(yearStr)
			month, _ = strconv.Atoi(monthStr)
			day, _ = strconv.Atoi(dayStr)

			// 日付が事前に入力されていなかった場合
		} else {
			dateStr := c.FormValue("date")
			if dateStr == "" {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_NO_CHOICE,
				})
			}
			yearStr, monthStr, dayStr = strings.SplitDateByHyphen(dateStr)
			year, _ = strconv.Atoi(yearStr)
			month, _ = strconv.Atoi(monthStr)
			day, _ = strconv.Atoi(dayStr)
		}

		startFreeTimeHourStr, startFreeTimeMinuteStr := c.FormValue("start-free-time-hour"), c.FormValue("start-free-time-minute")
		endFreeTimeHourStr, endFreeTimeMinuteStr := c.FormValue("end-free-time-hour"), c.FormValue("end-free-time-minute")
		if startFreeTimeHourStr == "" || startFreeTimeMinuteStr == "" || endFreeTimeHourStr == "" || endFreeTimeMinuteStr == "" {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_NO_CHOICE,
			})
		}
		startFreeTimeHour, _ := strconv.Atoi(startFreeTimeHourStr)
		startFreeTimeMinute, _ := strconv.Atoi(startFreeTimeMinuteStr)
		endFreeTimeHour, _ := strconv.Atoi(endFreeTimeHourStr)
		endFreeTimeMinute, _ := strconv.Atoi(endFreeTimeMinuteStr)
		if startFreeTimeHour > endFreeTimeHour {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_CHOICE_TIME,
			})
		} else if startFreeTimeHour == endFreeTimeHour && startFreeTimeMinute > endFreeTimeMinute {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_CHOICE_TIME,
			})
		} else if startFreeTimeHour == endFreeTimeHour && startFreeTimeMinute == endFreeTimeMinute {
			return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
				"error_message": entity.ERR_CHOICE_SAME_TIME,
			})
		}

		dateFreeTime, err := gateway.GetDateFreeTimeByUserIDAndDate(ctx, db, userID, year, month, day)
		// 存在しなかった場合はDateFreeTimeを作成する
		if err != nil {
			// DateFreeTime構造体で値を設定する
			dateFreeTime = &entity.DateFreeTime{
				UserID: userID,
				Year:   year,
				Month:  month,
				Day:    day,
			}
			// DateFreeTimeの作成
			dateFreeTime, err := gateway.CreateDateFreeTime(ctx, db, dateFreeTime)
			if err != nil {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}

			// FreeTime構造体へ値を設定する
			freeTime := &entity.FreeTime{
				DateFreeTimeID: dateFreeTime.ID,
				StartHour:      startFreeTimeHour,
				StartMinute:    startFreeTimeMinute,
				EndHour:        endFreeTimeHour,
				EndMinute:      endFreeTimeMinute,
			}
			// FreeTimeの作成
			_, err = gateway.CreateFreeTime(ctx, db, freeTime)
			if err != nil {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}

			// 既に存在していた場合はCreateDateFreeTimeは作成しない
		} else {
			// FreeTime構造体へ値を設定する
			freeTime := &entity.FreeTime{
				DateFreeTimeID: dateFreeTime.ID,
				StartHour:      startFreeTimeHour,
				StartMinute:    startFreeTimeMinute,
				EndHour:        endFreeTimeHour,
				EndMinute:      endFreeTimeMinute,
			}
			// FreeTimeの作成
			_, err = gateway.CreateFreeTime(ctx, db, freeTime)
			if err != nil {
				return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}
		}

		jpWeekday := time.GetWeekdayByDate(year, month, day)
		startFreeTimeHourStr = num.NumToFormattedString(startFreeTimeHour)
		startFreeTimeMinuteStr = num.NumToFormattedString(startFreeTimeMinute)
		endFreeTimeHourStr = num.NumToFormattedString(endFreeTimeHour)
		endFreeTimeMinuteStr = num.NumToFormattedString(endFreeTimeMinute)

		successMessage := fmt.Sprintf("%s/%s/%s（%s）%s:%s〜%s:%sでfree-timeを作成しました。", yearStr, monthStr, dayStr, jpWeekday, startFreeTimeHourStr, startFreeTimeMinuteStr, endFreeTimeHourStr, endFreeTimeMinuteStr)

		fmt.Println(successMessage)

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
