package handler

import (
	"errors"
	"strings"

	"subscriptions-service/internal/handler/dto"
)

func ValidateCreateRequest(req dto.CreateSubscriptionRequest) error {
	if strings.TrimSpace(req.ServiceName) == "" {
		return errors.New("service_name is required")
	}

	if req.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if req.UserID == "" {
		return errors.New("user_id is required")
	}

	if req.StartDate == "" {
		return errors.New("start_date is required")
	}

	return nil
}
