package freetime

import (
	"context"
	"fmt"
	"net/http"
	"src/internal/config"
	"src/internal/entity"
	"src/internal/pkg/models/freetimes"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CreateFreeTime(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		// user, _ := users.GetUserByUserID(ctx, db, userID)

		yearStr, monthStr, dayStr := c.FormValue("year"), c.FormValue("month"), c.FormValue("day")
		year, _ := strconv.Atoi(yearStr)
		month, _ := strconv.Atoi(monthStr)
		day, _ := strconv.Atoi(dayStr)

		startFreeTimeHourStr, startFreeTimeMinuteStr := c.FormValue("start-free-time-hour"), c.FormValue("start-free-time-minute")
		endFreeTimeHourStr, endFreeTimeMinuteStr := c.FormValue("end-free-time-hour"), c.FormValue("end-free-time-minute")
		startFreeTimeHour, _ := strconv.Atoi(startFreeTimeHourStr)
		startFreeTimeMinute, _ := strconv.Atoi(startFreeTimeMinuteStr)
		endFreeTimeHour, _ := strconv.Atoi(endFreeTimeHourStr)
		endFreeTimeMinute, _ := strconv.Atoi(endFreeTimeMinuteStr)

		fmt.Println(startFreeTimeHour)
		fmt.Println(startFreeTimeMinute)
		fmt.Println(endFreeTimeHour)
		fmt.Println(endFreeTimeMinute)

		var dateFreeTime *entity.DateFreeTime
		var freeTime *entity.FreeTime

		dateFreeTime, err := freetimes.GetDateFreeTime(ctx, db, userID, year, month, day)

		// 存在しなかった場合はDateFreeTimeを作成する
		if err != nil {
			// DateFreeTime構造体で値を設定する
			dateFreeTime := &entity.DateFreeTime{
				UserID: userID,
				Year:   year,
				Month:  month,
				Day:    day,
			}
			// DateFreeTimeの作成
			dateFreeTime, err := freetimes.CreateDateFreeTime(ctx, db, dateFreeTime)
			if err != nil {
				return c.Render(http.StatusOK, "create-free-time", echo.Map{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}

			// FreeTime構造体へ値を設定する
			freeTime = &entity.FreeTime{
				DateFreeTimeID: dateFreeTime.ID,
				StartHour:      startFreeTimeHour,
				StartMinute:    startFreeTimeMinute,
				EndHour:        endFreeTimeHour,
				EndMinute:      endFreeTimeMinute,
			}
			// FreeTimeの作成
			freeTime, err = freetimes.CreateFreeTime(ctx, db, freeTime)
			if err != nil {
				return c.Render(http.StatusOK, "create-free-time", echo.Map{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}

			// 既に存在していた場合はCreateDateFreeTimeは作成しない
		} else {
			// FreeTime構造体へ値を設定する
			freeTime = &entity.FreeTime{
				DateFreeTimeID: dateFreeTime.ID,
				StartHour:      startFreeTimeHour,
				StartMinute:    startFreeTimeMinute,
				EndHour:        endFreeTimeHour,
				EndMinute:      endFreeTimeMinute,
			}
			// FreeTimeの作成
			freeTime, err = freetimes.CreateFreeTime(ctx, db, freeTime)
			if err != nil {
				return c.Render(http.StatusOK, "create-free-time", echo.Map{
					"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
				})
			}
		}

		return c.JSON(http.StatusOK, "success")
	}
}
