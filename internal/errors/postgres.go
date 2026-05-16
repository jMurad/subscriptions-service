package errors

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgUniqueViolation  = "23505"
	pgCheckViolation   = "23514"
	pgNotNullViolation = "23502"
)

func MapPostgresError(err error) error {
	if err == nil {
		return nil
	}

	// context errors
	if errors.Is(err, context.Canceled) {
		return Wrap(err, ErrCanceled)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return Wrap(err, ErrTimeout)
	}

	// not found
	if errors.Is(err, pgx.ErrNoRows) {
		return Wrap(err, ErrNotFound)
	}

	// postgres errors
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return Wrap(err, ErrInternal)
	}

	switch pgErr.Code {

	case pgUniqueViolation:
		return mapUniqueViolation(err, pgErr)

	case pgCheckViolation:
		return mapCheckViolation(err, pgErr)

	case pgNotNullViolation:
		return mapNotNullViolation(err, pgErr)

	default:
		return Wrap(err, ErrInternal)
	}
}

func mapUniqueViolation(err error, pgErr *pgconn.PgError) error {
	switch pgErr.ConstraintName {

	case "uq_user_service_active":
		return WrapMessage(
			err,
			ErrConflict,
			"subscription already exists",
		)

	default:
		return WrapMessage(
			err,
			ErrConflict,
			"resource already exists",
		)
	}
}

func mapCheckViolation(err error, pgErr *pgconn.PgError) error {
	switch pgErr.ConstraintName {

	case "chk_dates":
		return WrapMessage(
			err,
			ErrValidation,
			"end_date must be greater than or equal to start_date",
		)

	case "chk_service_name":
		return WrapMessage(
			err,
			ErrValidation,
			"service_name must contain between 1 and 100 characters",
		)

	default:
		return WrapMessage(
			err,
			ErrValidation,
			"check constraint violation",
		)
	}
}

func mapNotNullViolation(err error, pgErr *pgconn.PgError) error {
	return WrapMessage(
		err,
		ErrValidation,
		pgErr.ColumnName+" is required",
	)
}
