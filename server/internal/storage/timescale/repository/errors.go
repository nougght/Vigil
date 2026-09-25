package repository

import (
	"errors"

	"github.com/jackc/pgconn"
)

// postgres error codes
const (
	ErrCodeUniqueViolation = "23505"
)

// repository errors
var (
	ErrNotFound       = errors.New("repository error not found")
	ErrNoAffectedRows = errors.New("repository error no affected rows")
	ErrConflict       = errors.New("repository error conflict")
)

func IsConflict(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == ErrCodeUniqueViolation {
		return true
	}
	return false
}
