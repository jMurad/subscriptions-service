package subscriptions

import (
	"net/http"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/transport/http/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Delete subscription
//
// @Summary Delete subscription
// @Description Soft delete subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = h.svc.Delete(
		r.Context(),
		id,
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
