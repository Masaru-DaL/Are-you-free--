package freetimes

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"src/internal/entity"

	"github.com/jmoiron/sqlx"
)

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

/*
free-timeの1件取得
指定した日付のfree-time（start&end timeは複数の場合もある）を取得
*/
func GetFreeTimeByDate(ctx context.Context, db *sqlx.DB) ([]*entity.DateFreeTime, error) {
	var freeTime []*entity.DateFreeTime

	err := db.SelectContext(ctx, &freeTime, `
		SELECT
			dft.id, dft.user_id, dft.year, dft.month, dft.day, dft.created_at, dft.updated_at,
			ft.id AS "free_time.id",

			ft.start_hour AS "free_time.start_hour",
			ft.start_minute AS "free_time.start_minute",
			ft.end_hour AS "free_time.end_hour",
			ft.end_minute AS "free_time.end_minute",
			ft.created_at AS "free_time.created_at",
			ft.updated_at AS "free_time.updated_at"
		FROM
			date_free_times AS dft
		LEFT JOIN
			free_times AS ft
		ON
			dft.id = ft.date_free_time_id
	`)

	if err != nil {
		return nil, entity.ErrSQLGetFailed
	}

	return freeTime, nil
}

/*
free-timeの1件取得
指定した日付のfree-time（start&end timeは複数の場合もある）を取得
*/
// func GetFreeTimeByDate(ctx context.Context, db *sqlx.DB, dateFreeTime *entity.DateFreeTime) ([]*entity.FreeTime, error) {
// 	var freeTimes []*entity.FreeTime

// 	err := db.SelectContext(ctx, &freeTimes, `
// 		SELECT
// 			dft.id, dft.user_id, dft.year, dft.month, dft.day,
// 			ft.start_hour AS "start_hour", ft.start_hour AS "start_minute", ft.end_hour AS "end_hour", ft.end_minute AS "end_minute"
// 		FROM
// 			date_free_times AS dft
// 		LEFT JOIN
// 			free_times AS ft
// 		ON
// 			dft.id = ft.date_free_time_id
// 		WHERE
// 			dft.user_id = ? AND dft.year = ? AND dft.month = ? AND dft.day = ?
// 	`, dateFreeTime.UserID, dateFreeTime.Year, dateFreeTime.Month, dateFreeTime.Day)

// 	if err != nil {
// 		return nil, entity.ErrSQLGetFailed
// 	}

// 	return freeTimes, nil
// }

// func GetDateFree(ctx context.Context, db *sqlx.DB, userID int) ([]*entity.DateFreeTime, error) {
// 	var dateFreeTimeArray []*entity.DateFreeTime
// 	err := db.SelectContext(ctx, &dateFreeTimeArray, `
// 		SELECT
// 			dft.id, dft.user_id, dft.year, dft.month, dft.day
// 		FROM
// 			date_free_times AS dft
// 		WHERE
// 			dft.user_id = ?
// 	`, userID)

// 	if err != nil {
// 		return nil, entity.ErrSQLGetFailed
// 	}

// 	return dateFreeTimeArray, nil
// }
// func GetFree(ctx context.Context, db *sqlx.DB, dateFreeTime *entity.DateFreeTime) ([]*entity.FreeTime, error) {
// 	var freeTimes []*entity.FreeTime

// 	err := db.SelectContext(ctx, &freeTimes, `
// 		SELECT
// 			ft.start_hour, ft.start_hour, ft.end_hour, ft.end_minute
// 		FROM
// 			free_times AS ft
// 		WHERE
// 			dft.date_free_time_id = ?
// 	`, dateFreeTime.ID)

// 	if err != nil {
// 		return nil, entity.ErrSQLGetFailed
// 	}

// 	return freeTimes, nil
// }

func CreateDateFreeTime(ctx context.Context, db *sqlx.DB, dateFreeTime *entity.DateFreeTime) (*entity.DateFreeTime, error) {
	stmt, err := db.PrepareNamedContext(ctx, `
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

func CreateFreeTime(ctx context.Context, db *sqlx.DB, dateFreeTimeID int, freeTime *entity.FreeTime) (*entity.FreeTime, error) {
	stmt, err := db.PrepareNamedContext(ctx, `
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
