package pgerr

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolation = "23505"

// IsUniqueViolation reports whether err comes from a write that breaks the named UNIQUE
// constraint or index.
func IsUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolation && pgErr.ConstraintName == constraint
}
