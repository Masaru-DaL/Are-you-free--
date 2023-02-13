package auth

import (
	"src/internal/config"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

/* セッションの設定 */
func InitSession() *sessions.CookieStore {
	cookieStore := sessions.NewCookieStore([]byte(config.Config.AUTH.SessionKey))

	return cookieStore
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
