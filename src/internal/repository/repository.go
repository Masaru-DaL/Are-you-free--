package repository

import (
	"context"
	"src/internal/entity"
	"src/internal/infra/dbutils"
	"src/internal/repository/gateway"

	"github.com/glassonion1/logz"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

/* 現在の日付からもっとも最新の日付を取得する */

/* 指定した日付のfree-timeを全て格納して返す */
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

// date-free-timeを作成（トランザクション対応）
func CreateDateFreeTime(ctx context.Context, db *sqlx.DB, dateFreeTime *entity.DateFreeTime) (*entity.DateFreeTime, error) {
	if err := dbutils.TXHandler(ctx, db, func(tx *sqlx.Tx) (err error) {
		dateFreeTime, err = gateway.CreateDateFreeTime(ctx, tx, dateFreeTime)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logz.Debugf(ctx, err.Error())
		return nil, errors.Wrap(err, "failed to create date free time")
	}

	return dateFreeTime, nil
}

// free-timeを作成(トランザクション対応)
func CreateFreeTime(ctx context.Context, db *sqlx.DB, freeTime *entity.FreeTime) (*entity.FreeTime, error) {
	if err := dbutils.TXHandler(ctx, db, func(tx *sqlx.Tx) (err error) {
		freeTime, err = gateway.CreateFreeTime(ctx, tx, freeTime)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logz.Debugf(ctx, err.Error())
		return nil, errors.Wrap(err, "failed to create free time")
	}

	return freeTime, nil
}
