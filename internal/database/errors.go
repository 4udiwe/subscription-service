package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeUniqueViolation     = "23505"
	CodeForeignKeyViolation = "23503"
)

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == CodeUniqueViolation
	}
	return false
}

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == CodeForeignKeyViolation
	}
	return false
}
