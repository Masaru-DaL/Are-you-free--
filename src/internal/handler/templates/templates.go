package templates

import (
	"net/http"

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
	return c.Render(http.StatusOK, "login", echo.Map{
		"error_message": nil,
	})
}

/* トップページ */
func TopPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "")
}

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
	return c.Render(http.StatusOK, "create-free-time", "")
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
