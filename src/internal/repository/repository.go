package repository

import (
	"context"
	"fmt"
	"src/internal/entity"
	"src/internal/infra/dbutils"
	"src/internal/repository/gateway"

	"github.com/glassonion1/logz"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// --------------------------------
// 		Get
// --------------------------------

/* 現在の日付からもっとも最新の日付を取得する */
func GetNearestDateFreeTime(ctx context.Context, db *sqlx.DB, userID string) (*entity.DateFreeTime, error) {
	// date-free-timeを全件取得する
	dateFreeTimes, err := gateway.ListDateFreeTime(ctx, db, userID)
	if err != nil {

		return nil, entity.ErrNoDateFreeTimeFound
	}

	latestDateFreeTime := dateFreeTimes[0]

	return latestDateFreeTime, err
}

/* 指定した日付のfree-timeを全て格納して返す */
func GetDateFreeTimeByID(ctx context.Context, db *sqlx.DB, dateFreeTimeID int) (*entity.DateFreeTime, error) {
	// ユーザの指定した日付の情報を取得する
	dateFreeTime, err := gateway.GetDateFreeTimeByID(ctx, db, dateFreeTimeID)
	if err != nil {

		return nil, entity.ErrNoDateFreeTimeFound
	}

	// 指定した日付の全てのfree-timeを取得する
	freeTimes, err := gateway.ListFreeTime(ctx, db, dateFreeTime.ID)
	if err != nil {

		return nil, entity.ErrNoFreeTimeFound
	}

	// dateFreeTimeにfree-timeを格納する
	dateFreeTime.FreeTimes = append(dateFreeTime.FreeTimes, freeTimes...)

	return dateFreeTime, nil
}

/* 指定した日付のfree-timeを全て格納して返す */
func GetDateFreeTimeByUserIDAndDate(ctx context.Context, db *sqlx.DB, userID, year, month, day string) (*entity.DateFreeTime, error) {
	// ユーザの指定した日付の情報を取得する
	dateFreeTime, err := gateway.GetDateFreeTimeByUserIDAndDate(ctx, db, userID, year, month, day)
	if err != nil {

		return nil, entity.ErrNoDateFreeTimeFound
	}

	// 指定した日付の全てのfree-timeを取得する
	freeTimes, err := gateway.ListFreeTime(ctx, db, dateFreeTime.ID)
	if err != nil {

		return nil, entity.ErrNoFreeTimeFound
	}
	fmt.Println(freeTimes)

	dateFreeTime.FreeTimes = append(dateFreeTime.FreeTimes, freeTimes...)

	return dateFreeTime, nil
}

func GetUserByUserID(ctx context.Context, db *sqlx.DB, userID string) (*entity.User, error) {
	user, err := gateway.GetUserByUserID(ctx, db, userID)
	if err != nil {

		return nil, entity.ErrNoUserFound
	}

	return user, nil
}

func ListUserIDSharedUserID(ctx context.Context, db *sqlx.DB, userID string) ([]*entity.SharedUser, error) {
	sharedUsers, err := gateway.ListUserIDSharedUserID(ctx, db, userID)
	if err != nil {

		return nil, entity.ErrNoUserFound
	}

	return sharedUsers, nil
}

func ListSharedUser(ctx context.Context, db *sqlx.DB, userIDSharedUserID []*entity.SharedUser) ([]*entity.User, error) {
	var users []*entity.User

	for _, us := range userIDSharedUserID {
		user, err := gateway.GetUserByUserID(ctx, db, us.UserID)
		if err != nil {

			return nil, entity.ErrNoUserFound
		}

		users = append(users, user)
	}

	return users, nil
}

func ListDateFreeTimes(ctx context.Context, db *sqlx.DB, users []*entity.User, year, month, day string) ([]*entity.DateFreeTime, error) {
	var dateFreeTimes []*entity.DateFreeTime

	for _, u := range users {
		dateFreeTime, err := gateway.GetDateFreeTimeByUserIDAndDate(ctx, db, u.ID, year, month, day)
		if err != nil {

			return nil, err
		}

		dateFreeTimes = append(dateFreeTimes, dateFreeTime)
	}

	return dateFreeTimes, nil
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
