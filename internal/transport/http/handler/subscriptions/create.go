package subscriptions

import (
	"encoding/json"
	"net/http"
	apperrors "subscriptions-service/internal/errors"
	dto "subscriptions-service/internal/transport/http/dto/subscriptions"
	mapper "subscriptions-service/internal/transport/http/mapper/subscriptions"
	"subscriptions-service/internal/transport/http/response"
)

// Create subscription
//
// @Summary Create subscription
// @Description Create a new user subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body subscriptions.CreateRequest true "Subscription payload"
// @Success 201 {object} subscriptions.CreateResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w,
			apperrors.WrapMessage(
				err,
				apperrors.ErrValidation,
				"invalid request body",
			),
		)

		return
	}

	sub, err := mapper.ToSubscription(req)
	if err != nil {
		response.Error(w, err)
		return
	}

	id, err := h.svc.Create(
		r.Context(),
		sub,
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w,
		http.StatusCreated,
		dto.CreateResponse{
			ID: id.String(),
		},
	)
}
