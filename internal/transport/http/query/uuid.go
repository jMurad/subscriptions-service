package query

import (
	apperrors "subscriptions-service/internal/errors"

	"github.com/google/uuid"
)

func ParseOptionalUUID(value string, field string) (*uuid.UUID, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(value)
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
