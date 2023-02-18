package freetimes

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"src/internal/entity"

	"github.com/jmoiron/sqlx"
)

/*
free-timeの1件取得
指定した日付のfree-time（start&end timeは複数の場合もある）を取得
*/
func GetFreeTimeByDate(ctx context.Context, db *sqlx.DB, year int, month int, day int) (*entity.DateFreeTime, error) {
	var freeTime entity.DateFreeTime

	err := db.GetContext(ctx, &freeTime, `
		SELECT
			ft.id, ft.user_id, ft.year, ft.month, ft.day, ft.created_at, ft.updated_at,
			sft.hour, sft.minute,
			eft.hour, eft.minute
		FROM
			free_times AS ft
		LEFT JOIN
			start_free_times AS sft
		ON
			ft.start_time_id = sft.id
		LEFT JOIN
			end_free_times AS eft
		ON
			ft.end_time_id = eft.id
		WHERE
			ft.year = ? AND ft.month = ? AND ft.day = ?
	`, year, month, day)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrNoFreeTimeFound
	}

	if err != nil {
		return nil, entity.ErrSQLGetFailed
	}

	return &freeTime, nil
}

func CreateDateFreeTime(ctx context.Context, db *sqlx.DB, freeTime *entity.DateFreeTime) (*entity.DateFreeTime, error) {
	stmt, err := db.PrepareNamedContext(ctx, `
		INSERT INTO free_times
		(
			year,
			month,
			day
		)
		VALUES
		(
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

func CreateStartFreeTime(ctx context.Context, db *sqlx.DB, startFreeTime *entity.StartFreeTime) (*entity.StartFreeTime, error) {
	stmt, err := db.PrepareNamedContext(ctx, `
		INSERT INTO free_times
		(
			hour,
			minute
		)
		VALUES
		(
			:hour,
			:minute
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

	result, err := stmt.Exec(startFreeTime)
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

	startFreeTime.ID = int(id)

	return startFreeTime, err
}

func CreateEndFreeTime(ctx context.Context, db *sqlx.DB, endFreeTime *entity.EndFreeTime) (*entity.EndFreeTime, error) {
	stmt, err := db.PrepareNamedContext(ctx, `
		INSERT INTO free_times
		(
			hour,
			minute
		)
		VALUES
		(
			:hour,
			:minute
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

	result, err := stmt.Exec(endFreeTime)
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

	endFreeTime.ID = int(id)

	return endFreeTime, err
}
