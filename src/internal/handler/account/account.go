package account

import (
	"database/sql"
	"net/http"
	"src/internal/entity"
	"src/internal/pkg/auth"
	"src/internal/pkg/models/users"
	"src/internal/pkg/strings"

	"github.com/labstack/echo/v4"
)

/* サインアップ */
func Signup(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// フォームから送られた値
		signupName := c.FormValue("username")
		signupPassword := c.FormValue("password")
		signupPasswordConfirmation := c.FormValue("password-confirmation")
		signupEmail := c.FormValue("e-mail")

		/* Checking input data from forms Start here. */
		// 空入力の場合
		if signupName == "" || signupPassword == "" || signupPasswordConfirmation == "" || signupEmail == "" {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Input_Empty,
			})
			// パスワードと再確認用パスワードが一致しなかった場合
		} else if signupPassword != signupPasswordConfirmation {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Value_Mismatched,
			})
		}
		// 入力された値が文字数制限を超えている場合
		if len(signupName) > entity.LimitCharCountOfUsername || len(signupPassword) > entity.LimitCharCountOfPassword {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Limit_Over,
			})
		}
		// 入力されたデータに空白が混入されている場合
		checkResultSignupName := strings.CheckWhitespaceInString(signupName)
		checkResultSignupPassword := strings.CheckWhitespaceInString(signupPassword)
		if checkResultSignupName || checkResultSignupPassword {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Input_Whitespace,
			})
		}
		// 入力されたデータにクオーテーションが混入されている場合
		checkResultSignupName = strings.CheckQuotationInString(signupName)
		checkResultSignupPassword = strings.CheckQuotationInString(signupPassword)
		if checkResultSignupName || checkResultSignupPassword {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Input_Quotation,
			})
		}
		// メールアドレスの形式が正しくない場合
		checkResultSignupEmail := strings.CheckEmailFormat(signupEmail)
		if !checkResultSignupEmail {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Email_Format,
			})
		}
		// 入力された名前でDBから情報を取得する
		user, err := users.UserReqUsername(db, signupName)
		// ユーザ情報が既に存在している場合
		if err != sql.ErrNoRows {
			// 入力されたパスワードをユーザ情報のパスワードと比較する
			err = auth.CompareHashAndPlaintext(user.Password, signupPassword)
			if err == nil {
				return c.Render(http.StatusOK, "signup", echo.Map{
					"error_message": entity.Err_Info_Already,
				})
			}
		}
		/* Checking input data from form Here. */

		// 入力されたパスワードのハッシュ化
		encryptedSignupPassword, err := auth.PasswordEncrypt(signupPassword)
		if err != nil {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Unexpected,
			})
		}
		// サインアップ情報の登録
		err = users.CreateUser(db, signupName, encryptedSignupPassword, signupEmail)
		if err != nil {
			return c.Render(http.StatusOK, "signup", echo.Map{
				"error_message": entity.Err_Unexpected,
			})
		}

		return c.Render(http.StatusCreated, "login", echo.Map{
			"error_message": nil,
		})
	}
}
