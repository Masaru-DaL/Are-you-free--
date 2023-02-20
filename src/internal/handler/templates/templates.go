package templates

import (
	"fmt"
	"net/http"
	"src/internal/config"
	"src/internal/pkg/strings"

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
// func TopPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		dateFreeTime, _ := freetimes.GetDateFreeTime(ctx, db, 3, 2023, 02, 20)
// 		fmt.Println(dateFreeTime)
// 		fmt.Println(dateFreeTime.ID)

// 		freeTime, err := freetimes.GetFreeTimeByDate(ctx, db)
// 		if err != nil {
// 			fmt.Println(err)
// 			fmt.Println("----------1111111111----------")
// 			fmt.Println(freeTime)
// 		} else {
// 			fmt.Println("----------2222222222----------")
// 			fmt.Println(freeTime)
// 			fmt.Println(freeTime[0])
// 		}

// 		return c.Render(http.StatusOK, "index", "")
// 	}
// }

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
	dateString := c.QueryParam("today")
	if dateString == "" {
		return c.Render(http.StatusOK, "create-free-time", echo.Map{
			"year":  nil,
			"month": nil,
			"day":   nil,
		})
	}

	year, month, day := strings.SplitDateByHyphen(dateString)

	return c.Render(http.StatusOK, "create-free-time", echo.Map{
		"year":  year,
		"month": month,
		"day":   day,
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
