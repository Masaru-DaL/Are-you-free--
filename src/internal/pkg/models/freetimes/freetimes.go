package freetimes

import (
	"context"

	"github.com/jmoiron/sqlx"
)

/*
free-timeの1件取得
指定した日付のfree-time（start&end timeは複数の場合もある）を取得
*/
func GetFreeTimeByDate(ctx context.Context, db *sqlx.DB, year int, month int, day int) (*entity.FreeTime, error) {

}
