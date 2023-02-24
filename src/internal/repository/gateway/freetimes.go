package gateway

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"src/internal/entity"

	"github.com/jmoiron/sqlx"
)

/* user_idとdateの条件に合う情報を、date_free_timesから1件取得する */
func GetDateFreeTime(ctx context.Context, db *sqlx.DB, userID int, year int, month int, day int) (*entity.DateFreeTime, error) {
	var dateFreeTime entity.DateFreeTime

	err := db.GetContext(ctx, &dateFreeTime, `
		SELECT
			id, user_id, year, month, day, created_at, updated_at
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

	id, err := result.LastInsertId()
	if err != nil {
		return nil, entity.ErrSQLLastInsertIdFailed
	}

	dateFreeTime.ID = int(id)

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

	id, err := result.LastInsertId()
	if err != nil {
		return nil, entity.ErrSQLLastInsertIdFailed
	}

	freeTime.ID = int(id)

	return freeTime, err
}

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
