package gateway

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"src/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// --------------------------------
// 		Get
// --------------------------------

func GetDateFreeTimeByID(ctx context.Context, db *sqlx.DB, dateFreeTimeID int) (*entity.DateFreeTime, error) {
	var dateFreeTime entity.DateFreeTime

	err := db.GetContext(ctx, &dateFreeTime, `
		SELECT
			id,
			user_id,
			year,
			month,
			day,
			created_at,
			updated_at
		FROM
			date_free_times
		WHERE
			id = ?
	`, dateFreeTimeID)

	if errors.Is(err, sql.ErrNoRows) {

		return nil, entity.ErrNoFreeTimeFound
	}

	if err != nil {

		return nil, entity.ErrSQLGetFailed
	}

	return &dateFreeTime, nil
}

func GetNearestDateFreeTime(ctx context.Context, db *sqlx.DB, userID string) (*entity.DateFreeTime, error) {
	var dateFreeTime entity.DateFreeTime

	err := db.GetContext(ctx, &dateFreeTime, `
		SELECT
			id,
			user_id,
			year,
			month,
			day,
			created_at,
			updated_at
		FROM
			date_free_times
		WHERE
			user_id = ?
		ORDER BY
			ABS((julianday(year || '-' || month || '-' || day) - julianday(date('now', 'localtime', 'start of day')))) ASC
			LIMIT 1;
	`, userID)

	if errors.Is(err, sql.ErrNoRows) {

		return nil, entity.ErrNoFreeTimeFound
	}

	if err != nil {

		return nil, entity.ErrSQLGetFailed
	}

	return &dateFreeTime, nil
}

/* user_idとdateの条件に合う情報を、date_free_timesから1件取得する */
func GetDateFreeTimeByUserIDAndDate(ctx context.Context, db *sqlx.DB, userID string, year string, month string, day string) (*entity.DateFreeTime, error) {
	var dateFreeTime entity.DateFreeTime

	err := db.GetContext(ctx, &dateFreeTime, `
		SELECT
			id,
			user_id,
			year,
			month,
			day,
			created_at,
			updated_at
		FROM
			date_free_times
		WHERE
			user_id = ? AND year = ? AND month = ? AND day = ?
	`, userID, year, month, day)

	if errors.Is(err, sql.ErrNoRows) {

		return nil, entity.ErrNoFreeTimeFound
	}

	if err != nil {

		return nil, entity.ErrSQLGetFailed
	}

	return &dateFreeTime, nil
}

/* ユーザと共有している人の中間テーブル情報を全て取得する */
func ListUserIDSharedUserID(ctx context.Context, db *sqlx.DB, userID string) ([]*entity.SharedUser, error) {
	var sharedUsers []*entity.SharedUser

	err := db.SelectContext(ctx, &sharedUsers, `
		SELECT
			user_id,
			shared_user_id,
			created_at,
			updated_at
		FROM
			shares
		WHERE
			user_id = ?
	`, userID)

	if err != nil {

		return nil, entity.ErrSQLGetFailed
	}

	return sharedUsers, nil
}

/* ユーザの登録した全てのfree-timeを取得する */
func ListDateFreeTime(ctx context.Context, db *sqlx.DB, userID string) ([]*entity.DateFreeTime, error) {
	var dateFreeTimes []*entity.DateFreeTime

	err := db.SelectContext(ctx, &dateFreeTimes, `
		SELECT
			id,
			user_id,
			year,
			month,
			day,
			created_at,
			updated_at
		FROM
			date_free_times
		WHERE
			user_id = ?
		ORDER BY
			year, month, day
	`, userID)

	if err != nil {

		return nil, entity.ErrSQLGetFailed
	}

	return dateFreeTimes, nil
}

/*
ユーザ情報の1件取得
指定したユーザIDのユーザの情報を取得する
*/
func GetUserByUserID(ctx context.Context, db *sqlx.DB, userID string) (*entity.User, error) {
	var user entity.User

	err := db.GetContext(ctx, &user, `
		SELECT
			id,
			name,
			password,
			email,
			is_admin,
			created_at,
			updated_at
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
			id,
			name,
			password,
			email,
			is_admin,
			created_at,
			updated_at
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

// --------------------------------
// 		List
// --------------------------------

/* date_free_time_idの条件に合う情報を、free-timesから全て取得する */
func ListFreeTime(ctx context.Context, db *sqlx.DB, dateFreeTimeID int) ([]*entity.FreeTime, error) {
	var freeTimes []*entity.FreeTime

	err := db.SelectContext(ctx, &freeTimes, `
		SELECT
			id, date_free_time_id, start_hour, start_minute, end_hour, end_minute, created_at, updated_at
		FROM
			free_times
		WHERE
			date_free_time_id = ?
	`, dateFreeTimeID)

	if err != nil {

		return nil, entity.ErrSQLGetFailed
	}

	return freeTimes, nil
}

// --------------------------------
// 		Create
// --------------------------------

/* date_free_timeを作成する */
func CreateDateFreeTime(ctx context.Context, tx *sqlx.Tx, dateFreeTime *entity.DateFreeTime) (*entity.DateFreeTime, error) {
	stmt, err := tx.PrepareNamedContext(ctx, `
		INSERT INTO date_free_times
		(
			user_id,
			year,
			month,
			day
		)
		VALUES
		(
			:user_id,
			:year,
			:month,
			:day
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

	result, err := stmt.Exec(dateFreeTime)
	if err != nil {
		log.Println(err)

		return nil, entity.ErrSQLExecFailed
	}

	cnt, err := result.RowsAffected()
	if err != nil || cnt != 1 {

		return nil, entity.ErrSQLResultNotDesired
	}

	return dateFreeTime, err
}

/* free_timeを作成する */
func CreateFreeTime(ctx context.Context, tx *sqlx.Tx, freeTime *entity.FreeTime) (*entity.FreeTime, error) {
	stmt, err := tx.PrepareNamedContext(ctx, `
		INSERT INTO free_times
		(
			date_free_time_id,
			start_hour,
			start_minute,
			end_hour,
			end_minute
		)
		VALUES
		(
			:date_free_time_id,
			:start_hour,
			:start_minute,
			:end_hour,
			:end_minute
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

	result, err := stmt.Exec(freeTime)
	if err != nil {
		log.Println(err)

		return nil, entity.ErrSQLExecFailed
	}

	cnt, err := result.RowsAffected()
	if err != nil || cnt != 1 {

		return nil, entity.ErrSQLResultNotDesired
	}

	return freeTime, err
}

/* ユーザの新規作成 */
func CreateUser(ctx context.Context, db *sqlx.DB, user *entity.User) (*entity.User, error) {
	stmt, err := db.PrepareNamedContext(ctx, `
		INSERT INTO users
		(
			id,
			name,
			password,
			email
		)
		VALUES
		(
			:id,
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

	uuid := uuid.New().String()
	user.ID = uuid
	result, err := stmt.Exec(user)
	if err != nil {
		log.Println(err)

		return nil, entity.ErrSQLExecFailed
	}

	cnt, err := result.RowsAffected()
	if err != nil || cnt != 1 {

		return nil, entity.ErrSQLResultNotDesired
	}

	return user, err
}

// --------------------------------
// 		Update
// --------------------------------

/* free-timeの更新 */
func UpdateFreeTime(ctx context.Context, tx *sqlx.Tx, freeTime *entity.FreeTime) (*entity.FreeTime, error) {
	stmt, err := tx.PrepareNamedContext(ctx, `
	UPDATE
		free_times
	SET
		start_hour = :start_hour,
		start_minute = :start_minute,
		end_hour = :end_hour,
		end_minute = :end_minute
	WHERE
		date_time_id = :date_time_id
	`)

	if err != nil {

		return nil, entity.ErrSQLCreateStmt
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	result, err := stmt.Exec(freeTime)
	if err != nil {

		return nil, entity.ErrSQLExecFailed
	}

	cnt, err := result.RowsAffected()
	if err != nil || cnt > 1 {

		return nil, entity.ErrSQLResultNotDesired
	}

	return freeTime, nil
}

// --------------------------------
// 		Delete
// --------------------------------
