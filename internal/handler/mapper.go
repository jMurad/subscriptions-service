package handler

import (
	"errors"
	"time"

	"subscriptions-service/internal/handler/dto"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
)

const subscriptionDateLayout = "01-2006"

// Request DTO to domain model
func toSubscriptionModel(req dto.CreateSubscriptionRequest) (model.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return model.Subscription{}, errors.New("invalid user_id")
	}

	startDate, err := time.Parse(
		subscriptionDateLayout,
		req.StartDate,
	)
	if err != nil {
		return model.Subscription{}, errors.New("invalid start_date")
	}

	var endDate *time.Time

	if req.EndDate != "" {

		parsedEndDate, err := time.Parse(
			subscriptionDateLayout,
			req.EndDate,
		)
		if err != nil {
			return model.Subscription{}, errors.New("invalid end_date")
		}

		endDate = &parsedEndDate
	}

	return model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

// Domain model to response DTO
func toSubscriptionResponse(sub model.Subscription) dto.SubscriptionResponse {
	response := dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format(subscriptionDateLayout),
	}

	if sub.EndDate != nil {

		formattedEndDate := sub.EndDate.Format(
			subscriptionDateLayout,
		)

		response.EndDate = &formattedEndDate
	}

	return response
}

// Domain models to response DTO list
func toSubscriptionResponseList(subscriptions []model.Subscription) []dto.SubscriptionResponse {
	response := make(
		[]dto.SubscriptionResponse,
		0,
		len(subscriptions),
	)

	for _, sub := range subscriptions {

		response = append(
			response,
			toSubscriptionResponse(sub),
		)
	}

	return response
}

// Id to response DTO
func toCreateSubscriptionResponse(id uuid.UUID) dto.CreateSubscriptionResponse {
	return dto.CreateSubscriptionResponse{
		ID: id.String(),
	}
}

// Request update DTO to domain update model
func toSubscriptionUpdateModel(req dto.UpdateSubscriptionRequest) (model.SubscriptionUpdate, error) {
	update := model.SubscriptionUpdate{
		ServiceName: req.ServiceName,
		Price:       req.Price,
	}

	if req.EndDate != nil {

		parsedEndDate, err := time.Parse(subscriptionDateLayout, *req.EndDate)
		if err != nil {
			return model.SubscriptionUpdate{}, errors.New("error invalid EndDate")
		}

		update.EndDate = &parsedEndDate
	}

	return update, nil
}
