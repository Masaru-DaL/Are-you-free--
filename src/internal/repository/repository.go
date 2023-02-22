package repository

import (
	"context"
	"src/internal/entity"
	"src/internal/repository/gateway"

	"github.com/jmoiron/sqlx"
)

// free-time: 指定した日付のfree-timeを全て格納して返す
func GetDateFreeTime(ctx context.Context, db *sqlx.DB, userID int, year int, month int, day int) (*entity.DateFreeTime, error) {
	// ユーザの指定した日付の情報を取得する
	dateFreeTime, err := gateway.GetDateFreeTime(ctx, db, userID, year, month, day)
	if err != nil {
		return nil, entity.ErrNoDateFreeTimeFound
	}

	// 指定した日付の全てのfree-timeを取得する
	freeTimes, err := gateway.ListFreeTime(ctx, db, dateFreeTime.ID)
	if err != nil {
		return nil, entity.ErrNoFreeTimeFound
	}

	dateFreeTime.FreeTimes = append(dateFreeTime.FreeTimes, freeTimes...)

	return dateFreeTime, nil
}
