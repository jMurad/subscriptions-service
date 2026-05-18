package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteHTTP(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	response := ErrorResponse{
		Error: ErrorBody{
			Code:    "internal_error",
			Message: "internal server error",
		},
	}

	switch {
	case errors.Is(err, ErrValidation):
		status = http.StatusBadRequest
		response.Error.Code = ErrValidation.Code
		response.Error.Message = err.Error()

	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
		response.Error.Code = ErrNotFound.Code
		response.Error.Message = err.Error()

	case errors.Is(err, ErrConflict):
		status = http.StatusConflict
		response.Error.Code = ErrConflict.Code
		response.Error.Message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(response)
}
