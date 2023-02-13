package account

import (
	"database/sql"
	"net/http"
	"src/internal/pkg/auth"

	"github.com/labstack/echo/v4"
)

/* サインアップ */
func Signup(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		signUpName := c.FormValue("username")
		signUpPassword := c.FormValue("password")
		signUpPasswordConfirmation := c.FormValue("password-confirmation")

		// フォームの入力チェック
		if signUpName == "" || signUpPassword == "" || signUpPasswordConfirmation == "" {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_No_Input,
			})
		} else if signUpPassword != signUpPasswordConfirmation {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Input_Mismatch,
			})
		}

		// 入力された文字数が文字数制限を超えているかどうか
		if len(signUpName) > entity.NumberOfCharacterLimitsOnTheForm || len(signUpPassword) > entity.NumberOfCharacterLimitsOnTheForm {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Character_Limit_Over,
			})
		}

		// 入力されたデータに空白が混入されていないかチェック
		signUpNameCheckResult := strmanip.CheckWhitespaceInString(signUpName)
		signUpPasswordCheckResult := strmanip.CheckWhitespaceInString(signUpPassword)
		if signUpNameCheckResult || signUpPasswordCheckResult {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Contains_Whitespace,
			})
		}

		// 入力された情報がDBに存在するかを確認する
		storedUser, err := users.UserWhereName(db, signUpName)
		if err != sql.ErrNoRows {
			// 入力された名前がDBに名前が存在した場合、入力されたパスワードをチェックする
			err := auth.CompareHashAndPlaintext(storedUser.Password, signUpPassword)
			if err != nil {
				return c.Render(http.StatusAccepted, "login", echo.Map{
					"message": entity.Err_Already_Exist,
				})
			}
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Already_Exist,
			})
		}

		// 入力されたパスワードのハッシュ化
		encryptSignupPassword, err := auth.PasswordEncrypt(signUpPassword)
		if err != nil {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Unexpected_Error,
			})
		}

		// ユーザの新規登録
		err = users.CreateUser(db, signUpName, encryptSignupPassword)
		if err != nil {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Unexpected_Error,
			})
		}

		signUpUser, err := users.UserWhereName(db, signUpName)
		if err != nil {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Unexpected_Error,
			})
		}

		// サインアップしたユーザに対して指定回数分のギフトコードを付与する
		signUpUserID := signUpUser.ID
		for i := 0; i < entity.NumberOfTimesChargeable; i++ {
			randomNumber := nummanip.CreateRandomNumber()
			createChargeErr := charges.CreateCharge(db, signUpUserID, randomNumber)
			if createChargeErr != nil {
				return c.Render(http.StatusAccepted, "signup", echo.Map{
					"message": entity.Err_Unexpected_Error,
				})
			}
		}

		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}
}
