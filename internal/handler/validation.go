package handler

import (
	"errors"
	"strconv"
	"strings"

	"subscriptions-service/internal/handler/dto"

	"github.com/google/uuid"
)

func validateCreateRequest(req dto.CreateSubscriptionRequest) error {
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

func validateListParams(limitParam string, offsetParam string) (int, int, error) {
	limit := 10
	offset := 0

	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil {
			return 0, 0, errors.New("invalid limit")
		}

		if parsedLimit <= 0 {
			return 0, 0, errors.New("limit must be greater than 0")
		}

		limit = parsedLimit
	}

	if offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err != nil {
			return 0, 0, errors.New("invalid offset")
		}

		if parsedOffset < 0 {
			return 0, 0, errors.New("offset must be non-negative")
		}

		offset = parsedOffset
	}

	return limit, offset, nil
}

func validateUpdateSubscriptionRequest(req dto.UpdateSubscriptionRequest) error {
	if req.Price != nil && *req.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	return nil
}

func validateTotalRequest(req dto.TotalRequest) error {
	if strings.TrimSpace(req.From) == "" {
		return errors.New("from is required")
	}

	if strings.TrimSpace(req.To) == "" {
		return errors.New("to is required")
	}

	return nil
}
