package postgres

import (
	"context"
	"database/sql"
)

func (c *CoreStorage) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}
