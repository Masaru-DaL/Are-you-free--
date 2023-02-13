package account

import (
	"database/sql"
	"fmt"
	"net/http"
	"src/internal/config"
	"src/internal/entity"
	"src/internal/pkg/auth"
	"src/internal/pkg/models/users"
	"src/internal/pkg/strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

/* サインアップ機能 */
func Signup(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// フォームから送信された値
		signupName := c.FormValue("username")
		signupPassword := c.FormValue("password")
		signupPasswordConfirmation := c.FormValue("password-confirmation")
		signupEmail := c.FormValue("e-mail")

		/* Checking input data from forms Start here. */
		// 空入力の場合
		if signupName == "" || signupPassword == "" || signupPasswordConfirmation == "" || signupEmail == "" {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_INPUT_EMPTY,
			})
			// パスワードと再確認用パスワードが一致しなかった場合
		} else if signupPassword != signupPasswordConfirmation {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_INPUT_MISSING,
			})
		}
		// 入力された値が文字数制限を超えている場合
		if len(signupName) > entity.LimitCharCountOfUsername || len(signupPassword) > entity.LimitCharCountOfPassword {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_INPUT_LIMIT_OVER,
			})
		}
		// 入力されたデータに空白が混入されている場合
		checkResultSignupName := strings.CheckWhitespaceInString(signupName)
		checkResultSignupPassword := strings.CheckWhitespaceInString(signupPassword)
		if checkResultSignupName || checkResultSignupPassword {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_NO_WHITESPACE,
			})
		}
		// 入力されたデータにクオーテーションが混入されている場合
		checkResultSignupName = strings.CheckQuotationInString(signupName)
		checkResultSignupPassword = strings.CheckQuotationInString(signupPassword)
		if checkResultSignupName || checkResultSignupPassword {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_NO_QUOTATION,
			})
		}
		// メールアドレスの形式が正しくない場合
		checkResultSignupEmail := strings.CheckEmailFormat(signupEmail)
		if !checkResultSignupEmail {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_FAILED_EMAIL_FORMAT,
			})
		}
		// 入力された名前でDBから情報を取得する
		user, err := users.UserReqUsername(db, signupName)
		// ユーザ情報が既に存在している場合
		if err != sql.ErrNoRows {
			// 入力されたパスワードをユーザ情報のパスワードと比較する
			err = auth.CompareHashAndPlaintext(user.Password, signupPassword)
			// err == nil => DBに既に存在している
			if err == nil {
				return c.Render(http.StatusOK, "signup", echo.Map{
					"error_message": entity.ERR_ALREADY_USER_EXISTS,
				})
			}
		}
		/* Checking input data from form Here. */

		// 入力されたパスワードのハッシュ化
		encryptedSignupPassword, err := auth.PasswordEncrypt(signupPassword)
		if err != nil {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
			})
		}
		// サインアップ情報の登録
		err = users.CreateUser(db, signupName, encryptedSignupPassword, signupEmail)
		if err != nil {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
			})
		}

		return c.Render(http.StatusCreated, "login", echo.Map{
			"error_message": nil,
		})
	}
}

/* ログイン機能 */
func Login(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// フォームから送信された値
		loginName := c.FormValue("username")
		loginPassword := c.FormValue("password")

		/* Checking input data from forms Start here. */
		// 空入力の場合
		if loginName == "" || loginPassword == "" {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_INPUT_EMPTY,
			})
		}
		// ユーザ名をチェックする
		user, err := users.UserReqUsername(db, loginName)
		// ユーザ情報が取得できなかった場合
		if err == sql.ErrNoRows {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.ERR_NO_USER,
			})
			// パスワードをチェックする
		} else {
			// 入力されたパスワードをユーザ情報のパスワードと比較する
			err = auth.CompareHashAndPlaintext(user.Password, loginPassword)
			// err != nil => DBに存在していない
			if err != nil {
				return c.Render(http.StatusOK, "signup", echo.Map{
					"error_message": entity.ERR_ALREADY_USER_EXISTS,
				})
			}
		}
		/* Checking input data from form Here. */

		/* Setting up a user's session From here. */
		session, _ := session.Get(config.Config.AUTH.Session, c)
		session.Options = &sessions.Options{
			Path:     config.Config.AUTH.SessionPath,
			MaxAge:   config.Config.AUTH.SessionExpirationSec * config.Config.AUTH.SessionExpirationDay,
			HttpOnly: true,
		}
		session.Values["UserID"] = user.ID
		session.Values["UserName"] = user.Name
		session.Save(c.Request(), c.Response())
		/* Setting up a user's session here. */

		fmt.Println(session)
		fmt.Println(session.Values)

		return c.Redirect(http.StatusSeeOther, "/index")
	}
}
