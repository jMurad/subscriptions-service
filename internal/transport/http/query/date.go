package query

import (
	"time"

	apperrors "subscriptions-service/internal/errors"
)

func ParseOptionalDate(value string, field string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return nil,
			apperrors.WrapMessage(
				err,
				apperrors.ErrValidation,
				"invalid "+field,
			)
	}

	return &parsed, nil
}
