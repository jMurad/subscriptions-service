package subscriptions

import (
	"net/http"

	apperrors "subscriptions-service/internal/errors"
	mapper "subscriptions-service/internal/transport/http/mapper/subscriptions"
	"subscriptions-service/internal/transport/http/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Get subscription by ID
//
// @Summary Get subscription by ID
// @Description Returns subscription details by subscription ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} subscriptions.SubscriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w,
			apperrors.WrapMessage(
				err,
				apperrors.ErrValidation,
				"invalid id",
			),
		)

		return
	}

	sub, err := h.svc.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w,
		http.StatusOK,
		mapper.ToSubscriptionResponse(sub),
	)
}
