package query

import (
	"net/http"
	"strconv"

	apperrors "subscriptions-service/internal/errors"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

func ParsePagination(r *http.Request) (int, int, error) {
	limit := defaultLimit
	offset := 0

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, 0,
				apperrors.WrapMessage(
					err,
					apperrors.ErrValidation,
					"invalid limit",
				)
		}

		limit = parsed
	}

	if value := r.URL.Query().Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, 0,
				apperrors.WrapMessage(
					err,
					apperrors.ErrValidation,
					"invalid offset",
				)
		}

		offset = parsed
	}

	if limit <= 0 {
		limit = defaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	if offset < 0 {
		return 0, 0,
			apperrors.WrapMessage(
				nil,
				apperrors.ErrValidation,
				"offset must be greater than or equal to 0",
			)
	}

	return limit, offset, nil
}
