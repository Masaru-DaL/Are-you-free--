package users

import (
	"database/sql"
	"log"
	"src/internal/entity"
)

/*
	ユーザ情報の1件取得

サインアップ時の存在チェック: 名前とパスワードの両方が同じユーザが存在するかどうか
*/
func UserReqUsernameAndPassword(db *sql.DB, username string, password string) (entity.User, error) {
	sqlStatement := "SELECT * FROM users WHERE name = ? AND password = ?"

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Printf("Failed to Prepare for a single retrieval operation in the users: %v", err)
		return entity.User{}, err
	}
	defer stmt.Close()

	user := entity.User{}
	err = stmt.QueryRow(username, password).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		log.Printf("Failed to QueryRow for a single retrieval operation in the users: %v", err)
		return entity.User{}, err
	}

	return user, err
}

/* ユーザの新規作成 */
func CreateUser(db *sql.DB, userName string, encryptedPassword string, email string) error {
	sqlStatement := "INSERT INTO users(name, password, email) VALUES(?, ?, ?)"

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Printf("Failed to Prepare for create retrieval operation in the users: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userName, encryptedPassword, email)
	if err != nil {
		log.Printf("Failed to Exec for create retrieval operation in the users: %v", err)
		return err
	}

	return err
}
