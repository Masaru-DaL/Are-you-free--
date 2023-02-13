package users

import (
	"database/sql"
	"log"
)

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
