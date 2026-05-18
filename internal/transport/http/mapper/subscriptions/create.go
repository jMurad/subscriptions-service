package subscriptions

import (
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	dto "subscriptions-service/internal/transport/http/dto/subscriptions"

	"github.com/google/uuid"
)

func ToSubscription(req dto.CreateRequest) (model.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return model.Subscription{},
			apperrors.WrapMessage(
				err,
				apperrors.ErrValidation,
				"invalid user_id",
			)
	}

	startDate, err := time.Parse(time.DateOnly, req.StartDate)
	if err != nil {
		return model.Subscription{},
			apperrors.WrapMessage(
				err,
				apperrors.ErrValidation,
				"invalid start_date",
			)
	}

	var endDate *time.Time

	if req.EndDate != nil {
		parsed, err := time.Parse(time.DateOnly, *req.EndDate)
		if err != nil {
			return model.Subscription{},
				apperrors.WrapMessage(
					err,
					apperrors.ErrValidation,
					"invalid end_date",
				)
		}

		endDate = &parsed
	}

	return model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
