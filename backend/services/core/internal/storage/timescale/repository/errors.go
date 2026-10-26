package repository

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// repository errors
var (
	ErrNotFound       = errors.New("repository error not found")
	ErrNoAffectedRows = errors.New("repository error no affected rows")
	ErrConflict       = errors.New("repository error conflict")
)

func IsConflict(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return true
	}
	return false
}

func IsInvalidForeignKey(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
		return true
	}
	return false
}
