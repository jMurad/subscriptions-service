package subscriptions

import (
	"net/http"

	mapper "subscriptions-service/internal/transport/http/mapper/subscriptions"
	"subscriptions-service/internal/transport/http/query"
	"subscriptions-service/internal/transport/http/response"
)

// List subscriptions
//
// @Summary List subscriptions
// @Description Returns paginated list of subscriptions
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Pagination limit"
// @Param offset query int false "Pagination offset"
// @Success 200 {array} subscriptions.SubscriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := query.ParsePagination(r)
	if err != nil {
		response.Error(w, err)
		return
	}

	subs, err := h.svc.List(
		r.Context(),
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
