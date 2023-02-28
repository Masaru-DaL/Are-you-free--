package templates

import (
	"context"
	"fmt"
	"net/http"
	"src/internal/entity"
	"src/internal/infra/config"
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
		sess, _ := session.Get(config.Config.Session.Name, c)
		// fmt.Println("----------1111111111----------")
		// fmt.Println(sess)
		// fmt.Println(sess.Values)
		// fmt.Println(sess.ID)
		userID := sess.Values[config.Config.Session.KeyName].(int)
		fmt.Println(userID)

		// jpWeekday := times.GetWeekdayByDate(2023, 2, 20)
		// fmt.Println(jpWeekday)

		dateFreeTimes, _ := gateway.ListDateFreeTime(ctx, db, userID)
		fmt.Println(dateFreeTimes)
		fmt.Println(dateFreeTimes[0])

		return c.Render(http.StatusOK, "index", "")
	}
}

// func TopPage(c echo.Context) error {
// 	return c.Render(http.StatusOK, "index", "")
// }

/* スケジュールページ */
func FreeTimePage(c echo.Context) error {
	dateStr := c.QueryParam("date")
	yearStr, monthStr, dayStr := strings.SplitDateByHyphen(dateStr)
	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)
	day, _ := strconv.Atoi(dayStr)
	weekday := times.GetWeekdayByDate(year, month, day)

	return c.Render(http.StatusOK, "free-time", map[string]interface{}{
		"year":      year,
		"month":     month,
		"day":       day,
		"year_str":  yearStr,
		"month_str": monthStr,
		"day_str":   dayStr,
		"weekday":   weekday,
	})
}

/* スケジュールページ */
func FreeTimesPage(c echo.Context) error {
	return c.Render(http.StatusOK, "free-times", "")
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
func UpdateFreeTimePage(c echo.Context) error {
	return c.Render(http.StatusOK, "update-free-time", "")
}

/* 自身のスケジュールを共有するページ */
func ShareWithSomeonePage(c echo.Context) error {
	return c.Render(http.StatusOK, "share-with-someone", "")
}

/* 相手のスケジュールを共有するページ */
func ShareWithMePage(c echo.Context) error {
	return c.Render(http.StatusOK, "share-with-me", "")
}

/* アカウントページ */
func AccountPage(c echo.Context) error {
	return c.Render(http.StatusOK, "account", "")
}

/* アカウント編集ページ */
func AccountEditPage(c echo.Context) error {
	return c.Render(http.StatusOK, "account-edit", "")
}

/* パスワードリセットページ */
func PasswordResetPage(c echo.Context) error {
	return c.Render(http.StatusOK, "password-reset", "")
}

/* パスワード再登録ページ */
func PasswordReRegistrationPage(c echo.Context) error {
	return c.Render(http.StatusOK, "password-re-registration", "")
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
