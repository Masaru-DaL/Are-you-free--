package auth

import "golang.org/x/crypto/bcrypt"

/* パスワードの暗号化 */
func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

/* 暗号化されたパスワードと平文のパスワードの比較 */
func CompareHashAndPlaintext(hash, plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
}
