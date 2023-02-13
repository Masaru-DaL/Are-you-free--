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

		/* フォームの入力チェック */
		// 空入力の場合
		if signupName == "" || signupPassword == "" || signupPasswordConfirmation == "" || signupEmail == "" {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Input_Empty,
			})
			// パスワードと再確認用パスワードが一致しなかった場合
		} else if signupPassword != signupPasswordConfirmation {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Value_Mismatched,
			})
		}
		// 入力された値が文字数制限を超えている場合
		if len(signupName) > entity.LimitCharCountOfUsername || len(signupPassword) > entity.LimitCharCountOfPassword {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Limit_Over,
			})
		}
		// 入力されたデータに空白が混入されていないかチェック
		checkResultSignupName := strings.CheckWhitespaceInString(signupName)
		checkResultSignupPassword := strings.CheckWhitespaceInString(signupPassword)
		if checkResultSignupName || checkResultSignupPassword {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Input_Whitespace,
			})
		}
		//

		// 入力されたパスワードのハッシュ化
		encryptedSignupPassword, err := auth.PasswordEncrypt(signupPassword)
		if err != nil {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Unexpected,
			})
		}

		// サインアップ情報の登録
		err = users.CreateUser(db, signupName, encryptedSignupPassword, signupEmail)
		if err != nil {
			return c.Render(http.StatusAccepted, "signup", echo.Map{
				"message": entity.Err_Unexpected,
			})
		}

		// // 入力された情報がDBに存在するかを確認する
		// storedUser, err := users.UserWhereName(db, signUpName)
		// if err != sql.ErrNoRows {
		// 	// 入力された名前がDBに名前が存在した場合、入力されたパスワードをチェックする
		// 	err := auth.CompareHashAndPlaintext(storedUser.Password, signUpPassword)
		// 	if err != nil {
		// 		return c.Render(http.StatusAccepted, "login", echo.Map{
		// 			"message": entity.Err_Already_Exist,
		// 		})
		// 	}
		// 	return c.Render(http.StatusAccepted, "signup", echo.Map{
		// 		"message": entity.Err_Already_Exist,
		// 	})
		// }

		// // 入力されたパスワードのハッシュ化
		// encryptSignupPassword, err := auth.PasswordEncrypt(signUpPassword)
		// if err != nil {
		// 	return c.Render(http.StatusAccepted, "signup", echo.Map{
		// 		"message": entity.Err_Unexpected_Error,
		// 	})
		// }

		// ユーザの新規登録
		// err = users.CreateUser(db, signUpName, encryptSignupPassword)
		// if err != nil {
		// 	return c.Render(http.StatusAccepted, "signup", echo.Map{
		// 		"message": entity.Err_Unexpected_Error,
		// 	})
		// }

		// signUpUser, err := users.UserWhereName(db, signUpName)
		// if err != nil {
		// 	return c.Render(http.StatusAccepted, "signup", echo.Map{
		// 		"message": entity.Err_Unexpected_Error,
		// 	})
		// }

		// // サインアップしたユーザに対して指定回数分のギフトコードを付与する
		// signUpUserID := signUpUser.ID
		// for i := 0; i < entity.NumberOfTimesChargeable; i++ {
		// 	randomNumber := nummanip.CreateRandomNumber()
		// 	createChargeErr := charges.CreateCharge(db, signUpUserID, randomNumber)
		// 	if createChargeErr != nil {
		// 		return c.Render(http.StatusAccepted, "signup", echo.Map{
		// 			"message": entity.Err_Unexpected_Error,
		// 		})
		// 	}
		// }

		return c.Redirect(http.StatusSeeOther, "/login")
	}
}
