package auth

import (
	"net/http"
	"src/internal/config"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

/* 非認証ルートにアクセスした場合に処理を行うミドルウェア */
func UnAuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// セッション情報の取得
		session, err := session.Get(config.Config.AUTH.Session, c)
		if err != nil {
			return handleUnAuthError(c)
		}
		// セッションIDが空ではなかった場合
		if session.Values["UserID"] != nil {
			return handleUnAuthError(c)
		}

		return next(c)
	}
}

/* 認証ルートにアクセスした場合にアクセスした場合に処理を行うミドルウェア */
func AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// セッション情報の取得
		session, err := session.Get(config.Config.AUTH.Session, c)
		if err != nil {
			return handleAuthError(c)
		}
		// セッションIDが空だった場合
		if session.Values["UserID"] == nil {
			return handleAuthError(c)
		}

		return next(c)
	}
}

/* 認証エラー処理 */
func handleUnAuthError(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/index")
}

/* 認証エラー処理 */
func handleAuthError(c echo.Context) error {
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
