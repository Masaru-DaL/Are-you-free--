package admin

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

/* user情報の全件取得 */
func Users(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sqlStatement := "SELECT * FROM users"

		stmt, err := db.Prepare(sqlStatement)
		if err != nil {
			log.Printf("Failed to Prepare for all retrieval operation in the users: %v", err)
			return err
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			log.Printf("Failed to Query for all retrieval operation in the users: %v", err)
			return err
		}
		defer rows.Close()

		users := entity.Users{}
		for rows.Next() {
			user := entity.User{}
			err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.Password,
				&user.IsAdmin,
				&user.Deposit,
				&user.CreatedAt,
				&user.UpdatedAt)

			if err != nil {
				log.Printf("Failed to Scan for all retrieval operation in the users: %v", err)
				return err
			}
			users.Users = append(users.Users, user)
		}

		return c.JSON(http.StatusOK, users)
	}
}
