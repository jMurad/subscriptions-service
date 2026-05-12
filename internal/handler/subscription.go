package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	dto "subscriptions-service/internal/dto/subscription"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date", http.StatusBadRequest)
		return
	}

	var endDate *time.Time

	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			http.Error(w, "invalid end_date", http.StatusBadRequest)
			return
		}

		endDate = &parsedEndDate
	}

	sub := model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	id, err := h.service.Create(r.Context(), sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CreateSubscriptionResponse{
		ID: id.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.SubscriptionResponse{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   subscription.StartDate.Format("01-2006"),
	}

	if subscription.EndDate != nil {
		formatted := subscription.EndDate.Format("01-2006")
		response.EndDate = &formatted
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil {
			limit = parsedLimit
		}
	}

	if offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err == nil {
			offset = parsedOffset
		}
	}

	subscriptions, err := h.service.List(
		r.Context(),
		limit,
		offset,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	response := make([]dto.SubscriptionResponse, 0, len(subscriptions))

	for _, subscription := range subscriptions {

		item := dto.SubscriptionResponse{
			ID:          subscription.ID,
			ServiceName: subscription.ServiceName,
			Price:       subscription.Price,
			UserID:      subscription.UserID,
			StartDate:   subscription.StartDate.Format("01-2006"),
		}

		if subscription.EndDate != nil {
			formatted := subscription.EndDate.Format("01-2006")
			item.EndDate = &formatted
		}

		response = append(response, item)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
