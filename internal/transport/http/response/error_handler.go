package response

import (
	"errors"
	"net/http"

	apperrors "subscriptions-service/internal/errors"
)

func Error(w http.ResponseWriter, err error) {
	var appErr *apperrors.Error

	if !errors.As(err, &appErr) {
		appErr = apperrors.ErrInternal
	}
	status := mapStatusCode(
		appErr.Code,
	)

	JSON(
		w,
		status,
		ErrorResponse{
			Error: ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		},
	)
}
