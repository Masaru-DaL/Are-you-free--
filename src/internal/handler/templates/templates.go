package templates

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

/* サインアップページ */
func SignupPage(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", "")
}

/* ログインページ */
func LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login", "")
}

/* トップページ */
func TopPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "")
}

/* スケジュールページ */
func SchedulePage(c echo.Context) error {
	return c.Render(http.StatusOK, "schedule", "")
}

/* スケジュールページ */
func SchedulesPage(c echo.Context) error {
	return c.Render(http.StatusOK, "schedule", "")
}

/* スケジュール作成ページ */
func CreateSchedulePage(c echo.Context) error {
	return c.Render(http.StatusOK, "create-schedule", "")
}

/* スケジュール更新ページ */
func UpdateSchedulePage(c echo.Context) error {
	return c.Render(http.StatusOK, "update-schedule", "")
}

/* スケジュール共有ページ */
func SharingPage(c echo.Context) error {
	return c.Render(http.StatusOK, "sharing", "")
}

/* アカウントページ */
func AccountPage(c echo.Context) error {
	return c.Render(http.StatusOK, "account", "")
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
