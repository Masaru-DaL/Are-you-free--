package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

/* 認証エラー処理 */
func HandleUnAuthError(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/index")
}

/* 認証エラー処理 */
func HandleAuthError(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/auth/login")
}

/* パスワードの暗号化 */
func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

/* 暗号化されたパスワードと平文のパスワードの比較 */
func CompareHashAndPlaintext(hash, plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
}
