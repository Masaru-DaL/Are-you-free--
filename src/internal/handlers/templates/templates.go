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

	"github.com/glassonion1/logz"
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
		userID := sess.Values[config.Config.Session.KeyName].(string)
		fmt.Println(userID)
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"nearest_date_free_time_id": nil,
			})
		} else if err != nil {

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
		// セッションからユーザ情報を取得する
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)
		user, _ := repository.GetUserByUserID(ctx, db, userID)

		switch c.Request().Method {
		case "GET":
			// date_free_time_idを指定し、dateFreeTimeを取得する
			dateFreeTimeIDStr := c.Param("id")
			dateFreeTimeID, _ := strconv.Atoi(dateFreeTimeIDStr)

			dateFreeTime, err := repository.GetDateFreeTimeByID(ctx, db, dateFreeTimeID)
			if err != nil {
				logz.Errorf(ctx, entity.ERR_INTERNAL_SERVER_ERROR)

				return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
					"error_message": "エラーが発生しました。",
				})
			}

			// 自身が共有している人の中間テーブルを取得する
			userIDsharedUserID, err := repository.ListUserIDSharedUserID(ctx, db, userID)
			if err != nil {
				logz.Errorf(ctx, entity.ERR_INTERNAL_SERVER_ERROR)

				return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
					"error_message": "エラーが発生しました。",
				})
			}

			// ユーザが共有しているユーザ情報を全て取得する
			sharedUsers, err := repository.ListSharedUser(ctx, db, userIDsharedUserID)
			if err != nil {
				logz.Errorf(ctx, entity.ERR_INTERNAL_SERVER_ERROR)

				return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
					"error_message": "エラーが発生しました。",
				})
			}

			// ユーザが共有しているユーザのdate-free-timeを全て取得する
			var sharedDateFreeTimes []*entity.DateFreeTime
			for _, us := range userIDsharedUserID {
				dateFreeTime, err := repository.GetDateFreeTimeByUserIDAndDate(ctx, db, us.SharedUserID, dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)
				if err != nil {
					logz.Errorf(ctx, entity.ERR_INTERNAL_SERVER_ERROR)

					return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
						"error_message": "エラーが発生しました。",
					})
				}

				sharedDateFreeTimes = append(sharedDateFreeTimes, dateFreeTime)
			}

			weekday := times.GetWeekdayByDate(dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)

			return c.Render(http.StatusOK, "free-time", map[string]interface{}{
				"year":                   dateFreeTime.Year,
				"month":                  dateFreeTime.Month,
				"day":                    dateFreeTime.Day,
				"weekday":                weekday,
				"user":                   user,
				"date_free_time":         dateFreeTime,
				"shared_users":           sharedUsers,
				"shared_date_free_times": sharedDateFreeTimes,
			})
		case "POST":
			date := c.FormValue("date")
			// POSTリクエストの処理
			return c.String(http.StatusOK, fmt.Sprintf("Date: %s", date))
		default:
			// サポートされていないHTTPメソッドの場合は、405 Method Not Allowedを返す
			return c.NoContent(http.StatusMethodNotAllowed)

			// dateStr := c.QueryParam("date")
			// yearStr, monthStr, dayStr := strings.SplitDateByHyphen(dateStr)

		}
	}
}

/* スケジュールページ */
func FreeTimesPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)
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
	jpWeekday := times.GetWeekdayByDate(yearStr, monthStr, dayStr)

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
		userID := sess.Values[config.Config.Session.KeyName].(string)
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
		userID := sess.Values[config.Config.Session.KeyName].(string)
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
		userID := sess.Values[config.Config.Session.KeyName].(string)
		nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "share-with-me", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* アカウントページ */
func AccountPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sess, _ := session.Get(config.Config.Session.Name, c)
		// userID := sess.Values[config.Config.Session.KeyName].(string)
		// nearestDateFreeTime, _ := gateway.GetNearestDateFreeTime(ctx, db, userID)

		return c.Render(http.StatusOK, "account", map[string]interface{}{
			// "nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* アカウント編集ページ */
func AccountEditPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)
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
		userID := sess.Values[config.Config.Session.KeyName].(string)
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
		userID := sess.Values[config.Config.Session.KeyName].(string)
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
