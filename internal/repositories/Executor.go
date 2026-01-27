package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type QueryExecutor interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

var _ QueryExecutor = (*sqlx.DB)(nil)
var _ QueryExecutor = (*sqlx.Tx)(nil)