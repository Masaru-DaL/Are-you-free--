package freetimes

import (
	"context"
	"database/sql"
	"errors"
	"src/internal/entity"

	"github.com/jmoiron/sqlx"
)

/*
free-timeの1件取得
指定した日付のfree-time（start&end timeは複数の場合もある）を取得
*/
func GetFreeTimeByDate(ctx context.Context, db *sqlx.DB, year int, month int, day int) (*entity.FreeTime, error) {
	var freeTime entity.FreeTime

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
