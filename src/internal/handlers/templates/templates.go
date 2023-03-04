package templates

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"src/internal/entity"
	"src/internal/infra/config"
	"src/internal/repository"
	"src/internal/repository/gateway"
	"src/internal/utils/strings"
	"src/internal/utils/times"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

/* サインアップページ */
func SignupPage(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", map[string]interface{}{
		"error_message": nil,
	})
}

/* ログインページ */
func LoginPage(c echo.Context) error {
	sess, _ := session.Get(config.Config.Session.Name, c)
	fmt.Println(sess)
	fmt.Println(sess.Values)

	return c.Render(http.StatusOK, "login", map[string]interface{}{
		"error_message": nil,
	})
}

/* トップページ */
func TopPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("session")
		sess, _ := session.Get(config.Config.Session.Name, c)
		fmt.Println("sessID: ", sess.ID)
		fmt.Println("sessValues: ", sess.Values)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		fmt.Println(userID)
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"nearest_date_free_time_id": nil,
			})
		}
		fmt.Println(nearestDateFreeTime)
		fmt.Println(nearestDateFreeTime.ID)

		// return c.Render(http.StatusOK, "index", "")
		return c.Render(http.StatusOK, "index", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

// func TopPage(c echo.Context) error {
// 	return c.Render(http.StatusOK, "index", "")
// }

/* スケジュールページ */
func FreeTimePage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)

		user, _ := repository.GetUserByUserID(ctx, db, userID)

		dateFreeTimeIDStr := c.Param("id")
		dateFreeTimeID, _ := strconv.Atoi(dateFreeTimeIDStr)

		dateFreeTime, _ := gateway.GetDateFreeTimeByID(ctx, db, dateFreeTimeID)
		year, _ := strconv.Atoi(dateFreeTime.Year)
		month, _ := strconv.Atoi(dateFreeTime.Month)
		day, _ := strconv.Atoi(dateFreeTime.Day)
		dateFreeTime, _ = repository.GetDateFreeTime(ctx, db, userID, dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)

		// dateStr := c.QueryParam("date")
		// yearStr, monthStr, dayStr := strings.SplitDateByHyphen(dateStr)
		weekday := times.GetWeekdayByDate(year, month, day)
		fmt.Println("----------1111111111----------")
		fmt.Println(dateFreeTime)

		fmt.Println(userID)
		fmt.Println(year)
		fmt.Println(month)
		fmt.Println(day)
		return c.Render(http.StatusOK, "free-time", map[string]interface{}{
			"year":           year,
			"month":          month,
			"day":            day,
			"year_str":       dateFreeTime.Year,
			"month_str":      dateFreeTime.Month,
			"day_str":        dateFreeTime.Day,
			"weekday":        weekday,
			"user":           user,
			"date_free_time": dateFreeTime,
		})
	}
}

/* スケジュールページ */
func FreeTimesPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "free-times", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* スケジュール作成ページ */
func CreateFreeTimePage(c echo.Context) error {
	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
			"year":          nil,
			"month":         nil,
			"day":           nil,
			"weekday":       nil,
			"error_message": nil,
		})
	}
	// 入力された日付情報が現在の日付より後かチェックする
	isAfterCurrentTime := times.IsAfterCurrentTime(dateStr)
	if !isAfterCurrentTime {
		return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
			"error_message": entity.ERR_CHOICE_DATE,
		})
	}

	yearStr, monthStr, dayStr := strings.SplitDateByHyphen(dateStr)
	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)
	day, _ := strconv.Atoi(dayStr)
	jpWeekday := times.GetWeekdayByDate(year, month, day)

	return c.Render(http.StatusOK, "create-free-time", map[string]interface{}{
		"year_str":      yearStr,
		"month_str":     monthStr,
		"day_str":       dayStr,
		"weekday":       jpWeekday,
		"error_message": nil,
	})
}

/* スケジュール更新ページ */
func UpdateFreeTimePage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "update-free-time", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* 自身のスケジュールを共有するページ */
func ShareWithSomeonePage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "share-with-someone", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* 相手のスケジュールを共有するページ */
func ShareWithMePage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "share-with-me", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* アカウントページ */
func AccountPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "account", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* アカウント編集ページ */
func AccountEditPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "account-edit", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* パスワードリセットページ */
func PasswordResetPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "password-reset", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* パスワード再登録ページ */
func PasswordReRegistrationPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "password-re-registration", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* htmlのformにもPUTやDELETEにも対応させるメソッド */
func MethodOverride(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == "POST" {
			method := c.Request().PostFormValue("_method")
			if method == "PUT" || method == "PATCH" || method == "DELETE" {
				c.Request().Method = method
			}
		}
		return next(c)
	}
}
