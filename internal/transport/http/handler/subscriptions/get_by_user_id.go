package subscriptions

import (
	"net/http"

	apperrors "subscriptions-service/internal/errors"
	mapper "subscriptions-service/internal/transport/http/mapper/subscriptions"
	"subscriptions-service/internal/transport/http/query"
	"subscriptions-service/internal/transport/http/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Get subscriptions by user ID
//
// @Summary Get subscriptions by user ID
// @Description Returns paginated list of user subscriptions
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Pagination limit"
// @Param offset query int false "Pagination offset"
// @Success 200 {array} subscriptions.SubscriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/user/{user_id} [get]
func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(
		r,
		"user_id",
	)

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		response.Error(w,
			apperrors.WrapMessage(
				err,
				apperrors.ErrValidation,
				"invalid user_id",
			),
		)

		return
	}

	limit, offset, err := query.ParsePagination(r)
	if err != nil {
		response.Error(w, err)
		return
	}

	subs, err := h.svc.GetByUserID(
		r.Context(),
		userID,
		limit,
		offset,
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w,
		http.StatusOK,
		mapper.ToSubscriptionResponses(subs),
	)
}
