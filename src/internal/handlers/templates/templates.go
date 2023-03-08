package templates

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusOK, "index", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* スケジュールページ */
func FreeTimePage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// セッションからユーザ情報を取得する
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)
		user, _ := repository.GetUserByUserID(ctx, db, userID)

		switch c.Request().Method {
		case "GET":
			// 指定された直近のdateFreeTimeを取得する
			dateFreeTimeIDStr := c.Param("id")
			dateFreeTimeID, _ := strconv.Atoi(dateFreeTimeIDStr)

			dateFreeTime, err := repository.GetDateFreeTimeByID(ctx, db, dateFreeTimeID)
			if errors.Is(err, entity.ErrNoDateFreeTimeFound) {
				log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

				return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
					"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
				})
			} else if err != nil {
				log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

				return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
					"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
				})
			}

			// 曜日を取得する
			weekday := times.GetWeekdayByDate(dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)

			// 自身が共有している人の中間テーブルを取得する
			userIDsharedUserID, err := repository.ListUserIDSharedUserID(ctx, db, userID)
			// 共有している人がいた場合
			if userIDsharedUserID != nil {
				if err != nil {
					log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

					return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
						"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
					})
				}

				// ユーザが共有しているユーザ情報を全て取得する
				sharedUsers, err := repository.ListSharedUser(ctx, db, userIDsharedUserID)
				if err != nil {
					log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

					return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
						"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
					})
				}

				// ユーザが共有しているユーザのdate-free-timeを全て取得する
				sharedDateFreeTimes, err := repository.ListDateFreeTimes(ctx, db, sharedUsers, dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)
				if err != nil {
					log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

					return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
						"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
					})
				}

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

				// 共有している人がいなかった場合
			} else {

				return c.Render(http.StatusOK, "free-time", map[string]interface{}{
					"year":                   dateFreeTime.Year,
					"month":                  dateFreeTime.Month,
					"day":                    dateFreeTime.Day,
					"weekday":                weekday,
					"user":                   user,
					"date_free_time":         dateFreeTime,
					"shared_users":           nil,
					"shared_date_free_times": nil,
				})
			}

		case "POST":
			// POSTされた日付データからdateFreeTimeを取得する
			date := c.FormValue("date")
			year, month, day := strings.SplitDateByHyphen(date)

			dateFreeTime, err := repository.GetDateFreeTimeByUserIDAndDate(ctx, db, userID, year, month, day)
			if errors.Is(err, entity.ErrNoDateFreeTimeFound) {
				log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

				return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
					"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
				})
			} else if err != nil {
				log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

				return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
					"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
				})
			}

			// 曜日を取得する
			weekday := times.GetWeekdayByDate(dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)

			// 自身が共有している人の中間テーブルを取得する
			userIDsharedUserID, err := repository.ListUserIDSharedUserID(ctx, db, userID)
			// 共有している人がいた場合
			if userIDsharedUserID != nil {
				if err != nil {
					log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

					return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
						"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
					})
				}

				// ユーザが共有しているユーザ情報を全て取得する
				sharedUsers, err := repository.ListSharedUser(ctx, db, userIDsharedUserID)
				if err != nil {
					log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

					return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
						"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
					})
				}

				// ユーザが共有しているユーザのdate-free-timeを全て取得する
				sharedDateFreeTimes, err := repository.ListDateFreeTimes(ctx, db, sharedUsers, dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)
				if err != nil {
					log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

					return c.Render(http.StatusInternalServerError, "free-time", map[string]interface{}{
						"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
					})
				}

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

				// 共有している人がいなかった場合
			} else {

				return c.Render(http.StatusOK, "free-time", map[string]interface{}{
					"year":                   dateFreeTime.Year,
					"month":                  dateFreeTime.Month,
					"day":                    dateFreeTime.Day,
					"weekday":                weekday,
					"user":                   user,
					"date_free_time":         dateFreeTime,
					"shared_users":           nil,
					"shared_date_free_times": nil,
				})
			}

		default:
			// サポートされていないHTTPメソッドの場合は、TOPページへリダイレクトする
			return c.Redirect(http.StatusSeeOther, "/index")
		}
	}
}

/* スケジュールページ */
func FreeTimesPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

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
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusOK, "update-free-time", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* 自身のスケジュールを共有するページ */
func ShareWithSomeonePage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusOK, "share-with-someone", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* 相手のスケジュールを共有するページ */
func ShareWithMePage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusOK, "share-with-me", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* アカウントページ */
func AccountPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusOK, "account", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* アカウント編集ページ */
func AccountEditPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusOK, "account-edit", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* パスワードリセットページ */
func PasswordResetPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusOK, "password-reset", map[string]interface{}{
			"nearest_date_free_time_id": nearestDateFreeTime.ID,
		})
	}
}

/* パスワード再登録ページ */
func PasswordReRegistrationPage(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sessionを取得し、userIDを取得
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(string)

		// userの現在の日付から最も近いfree-timeを取得する
		nearestDateFreeTime, err := gateway.GetNearestDateFreeTime(ctx, db, userID)
		if errors.Is(err, entity.ErrNoFreeTimeFound) {
			log.Printf(entity.ERR_NO_DATE_FREE_TIME_FOUND+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_NO_FREE_TIME_FOUND,
			})
		} else if err != nil {
			log.Printf(entity.ERR_INTERNAL_SERVER_ERROR+": %v", err)

			return c.Render(http.StatusOK, "index", map[string]interface{}{
				"error_message": entity.MESSAGE_INTERNAL_SERVER_ERROR,
			})
		}

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
