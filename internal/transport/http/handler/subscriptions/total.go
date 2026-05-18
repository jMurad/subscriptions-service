package subscriptions

import (
	"net/http"

	dto "subscriptions-service/internal/transport/http/dto/subscriptions"
	"subscriptions-service/internal/transport/http/query"
	"subscriptions-service/internal/transport/http/response"
)

// Calculate subscriptions total
//
// @Summary Calculate subscriptions total cost
// @Description Returns total subscriptions cost with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} subscriptions.TotalResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/total [get]
func (h *Handler) Total(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	userID, err := query.ParseOptionalUUID(
		queryValues.Get("user_id"),
		"user_id",
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	serviceName := query.ParseOptionalString(
		queryValues.Get("service_name"),
	)

	from, err := query.ParseOptionalDate(
		queryValues.Get("from"),
		"from",
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	to, err := query.ParseOptionalDate(
		queryValues.Get("to"),
		"to",
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	total, err := h.svc.Total(
		r.Context(),
		userID,
		serviceName,
		from,
		to,
	)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		dto.TotalResponse{
			Total: total,
		},
	)
}
