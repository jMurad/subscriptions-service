package subscriptions

import (
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	dto "subscriptions-service/internal/transport/http/dto/subscriptions"
)

func ToSubscriptionUpdate(req dto.UpdateRequest) (model.SubscriptionUpdate, error) {
	var endDate *time.Time

	if req.EndDate != nil {
		parsed, err := time.Parse(time.DateOnly, *req.EndDate)
		if err != nil {
			return model.SubscriptionUpdate{},
				apperrors.WrapMessage(
					err,
					apperrors.ErrValidation,
					"invalid end_date",
				)
		}

		endDate = &parsed
	}

	return model.SubscriptionUpdate{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		EndDate:     endDate,
	}, nil
}
