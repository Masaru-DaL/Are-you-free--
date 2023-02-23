package dbutils

import (
	"context"
	"fmt"

	"github.com/glassonion1/logz"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TXHandler is handler for working with transaction.
// This is wrapper function for commit and rollback.
func TXHandler(ctx context.Context, db *sqlx.DB, f func(*sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return errors.Wrap(err, "start transaction failed")
	}

	defer func() {
		if p := recover(); p != nil || err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != nil {
				logz.Criticalf(ctx, fmt.Sprintf("rollback failed: %v, err: %v", rollBackErr, err))
			}

			logz.Debugf(ctx, "Rollback operation")

			if p != nil {
				err = errors.Wrap(err, fmt.Sprintf("recovered: %v", p))
			} else {
				err = entity.ErrSQLTransactionError
			}
		}
	}()

	if err := f(tx); err != nil {
		return err
	}

	return nil
}
