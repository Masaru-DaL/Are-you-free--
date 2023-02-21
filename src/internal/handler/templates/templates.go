package templates

import (
	"context"
	"fmt"
	"net/http"
	"src/internal/config"
	"src/internal/pkg/strings"
	"src/internal/pkg/time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

/* サインアップページ */
func SignupPage(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", echo.Map{
		"error_message": nil,
	})
}

/* ログインページ */
func LoginPage(c echo.Context) error {
	sess, _ := session.Get(config.Config.Session.Name, c)
	fmt.Println(sess)
	fmt.Println(sess.Values)

	return c.Render(http.StatusOK, "login", echo.Map{
		"error_message": nil,
	})
}

/* トップページ */
func TopPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		jpWeekday := time.GetWeekdayByDate(2023, 2, 20)
		fmt.Println(jpWeekday)

		return c.Render(http.StatusOK, "index", "")
	}
}

// func TopPage(c echo.Context) error {
// 	return c.Render(http.StatusOK, "index", "")
// }

/* スケジュールページ */
func FreeTimePage(c echo.Context) error {
	return c.Render(http.StatusOK, "free-time", "")
}

/* スケジュールページ */
func FreeTimesPage(c echo.Context) error {
	return c.Render(http.StatusOK, "free-times", "")
}

/* スケジュール作成ページ */
func CreateFreeTimePage(c echo.Context) error {
	dateString := c.QueryParam("date")
	if dateString == "" {
		return c.Render(http.StatusOK, "create-free-time", echo.Map{
			"year":          nil,
			"month":         nil,
			"day":           nil,
			"weekday":       nil,
			"error_message": nil,
		})
	}

	year, month, day := strings.SplitDateByHyphen(dateString)
	jpWeekday := time.GetWeekdayByDate(2023, 2, 20)

	return c.Render(http.StatusOK, "create-free-time", echo.Map{
		"year":          year,
		"month":         month,
		"day":           day,
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
