package subscriptions

import (
	"encoding/json"
	"net/http"

	apperrors "subscriptions-service/internal/errors"
	dto "subscriptions-service/internal/transport/http/dto/subscriptions"
	mapper "subscriptions-service/internal/transport/http/mapper/subscriptions"
	"subscriptions-service/internal/transport/http/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Update subscription
//
// @Summary Update subscription
// @Description Partially update subscription fields
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body subscriptions.UpdateRequest true "Update payload"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 408 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(
			w,
			apperrors.WrapMessage(
				err,
				apperrors.ErrValidation,
				"invalid id",
			),
		)

		return
	}

	var req dto.UpdateRequest

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

	update, err := mapper.ToSubscriptionUpdate(req)
	if err != nil {
		response.Error(w, err)
		return
	}

	err = h.svc.Update(
		r.Context(),
		id,
		update,
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w,
		http.StatusOK,
		nil,
	)
}
