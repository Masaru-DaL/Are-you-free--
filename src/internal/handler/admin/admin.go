package admin

import (
	"context"
	"fmt"
	"net/http"
	"src/internal/entity"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

/* user情報の全件取得 */
func GetUsers(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []*entity.User

		err := db.SelectContext(ctx, &users, `
			SELECT
				id, name, password, email, is_admin,
				datetime(created_at, '+9 hours') as created_at,
				datetime(updated_at, '+9 hours') as updated_at
			FROM
				users
		`)

		if err != nil {
			fmt.Println(err)
			return entity.ErrSQLGetFailed
		}

		return c.JSON(http.StatusOK, users)
	}
}
