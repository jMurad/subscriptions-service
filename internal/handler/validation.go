package handler

import (
	"errors"
	"strings"

	"subscriptions-service/internal/handler/dto"

	"github.com/google/uuid"
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

func validateUUID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid")
	}

	return parsedID, nil
}
