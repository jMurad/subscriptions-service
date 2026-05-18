package subscriptions

import (
	"subscriptions-service/internal/model"
	dto "subscriptions-service/internal/transport/http/dto/subscriptions"
	"time"
)

func ToSubscriptionResponse(sub *model.Subscription) dto.SubscriptionResponse {
	var endDate *string

	if sub.EndDate != nil {
		formatted := sub.EndDate.Format(time.DateOnly)
		endDate = &formatted
	}

	return dto.SubscriptionResponse{
		ID:          sub.ID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format(time.DateOnly),
		EndDate:     endDate,
	}
}

func ToSubscriptionResponses(subs []model.Subscription) []dto.SubscriptionResponse {
	result := make([]dto.SubscriptionResponse, 0, len(subs))

	for _, sub := range subs {
		copySub := sub
		result = append(result, ToSubscriptionResponse(&copySub))
	}

	return result
}
