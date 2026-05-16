package response

import (
	"net/http"

	apperrors "subscriptions-service/internal/errors"
)

func mapStatusCode(code string) int {
	switch code {

	case apperrors.ErrValidation.Code:
		return http.StatusBadRequest

	case apperrors.ErrNotFound.Code:
		return http.StatusNotFound

	case apperrors.ErrConflict.Code:
		return http.StatusConflict

	case apperrors.ErrTimeout.Code:
		return http.StatusRequestTimeout

	case apperrors.ErrCanceled.Code:
		return http.StatusRequestTimeout

	default:
		return http.StatusInternalServerError
	}
}
