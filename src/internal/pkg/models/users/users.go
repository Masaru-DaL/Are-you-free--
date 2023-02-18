package users

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"src/internal/entity"

	"github.com/jmoiron/sqlx"
)

/*
ユーザ情報の1件取得
指定したユーザIDのユーザの情報を取得する
*/
func GetUserByUserID(ctx context.Context, db *sqlx.DB, userID int) (*entity.User, error) {
	var user entity.User

	err := db.GetContext(ctx, &user, `
		SELECT
			*
		FROM
			users
		WHERE
			id = ?
	`, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrNoUserFound
	}

	if err != nil {
		return nil, entity.ErrSQLGetFailed
	}

	return &user, err
}

/*
ユーザ情報の1件取得
指定した名前のユーザの情報を取得する
*/
func GetUserByUsername(ctx context.Context, db *sqlx.DB, username string) (*entity.User, error) {
	var user entity.User

	err := db.GetContext(ctx, &user, `
		SELECT
			*
		FROM
			users
		WHERE
			name = ?
	`, username)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrNoUserFound
	}

	if err != nil {
		return nil, entity.ErrSQLGetFailed
	}

	return &user, err
}

/* ユーザの新規作成 */
func CreateUser(ctx context.Context, db *sqlx.DB, user *entity.User) (*entity.User, error) {
	stmt, err := db.PrepareNamedContext(ctx, `
		INSERT INTO users
		(
			name,
			password,
			email
		)
		VALUES
		(
			:name,
			:password,
			:email
		)
	`)

	if err != nil {
		log.Println(err)
		return nil, entity.ErrSQLCreateStmt
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	result, err := stmt.Exec(user)
	if err != nil {
		log.Println(err)
		return nil, entity.ErrSQLExecFailed
	}

	cnt, err := result.RowsAffected()
	if err != nil || cnt != 1 {
		return nil, entity.ErrSQLResultNotDesired
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, entity.ErrSQLLastInsertIdFailed
	}

	user.ID = int(id)

	return user, err
}
