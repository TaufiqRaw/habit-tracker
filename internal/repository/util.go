package repository

import (
	"context"
	"database/sql"
)

type execable interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}