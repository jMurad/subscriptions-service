package handler

import (
	"encoding/json"
	"net/http"
	"subscriptions-service/internal/handler/dto"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	log := logger.FromContext(r.Context())

	log.Info("create subscription request received")

	var req dto.CreateSubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode create subscription request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateCreateRequest(req)
	if err != nil {
		log.Warn("create subscription validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscription, err := toSubscriptionModel(req)
	if err != nil {
		log.Warn("failed to map create subscription request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), subscription)
	if err != nil {
		log.Error("failed to create subscription", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := toCreateSubscriptionResponse(id)

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("failed to encode create subscription response", zap.Error(err))
		return
	}

	log.Info("subscription created successfully", zap.String("subscription_id", id.String()))
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	idParam := chi.URLParam(r, "id")

	id, err := validateUUID(idParam)
	if err != nil {
		log.Warn("invalid subscription id", zap.String("subscription_id", idParam), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Warn("subscription not found", zap.String("subscription_id", id.String()), zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := toSubscriptionResponse(*subscription)

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("failed to encode get subscription response", zap.Error(err))
		return
	}

	log.Info("subscription returned successfully", zap.String("subscription_id", id.String()))
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	limit, offset, err := validateListParams(limitParam, offsetParam)
	if err != nil {
		log.Warn("invalid list query params",
			zap.String("limit", limitParam),
			zap.String("offset", offsetParam),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscriptions, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		log.Error("failed to list subscriptions", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := toSubscriptionResponseList(subscriptions)

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("failed to encode subscriptions list response", zap.Error(err))
		return
	}

	log.Info("subscriptions listed successfully",
		zap.Int("count", len(subscriptions)),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Warn("invalid subscription id", zap.String("subscription_id", idParam), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req dto.UpdateSubscriptionRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode update subscription request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateUpdateSubscriptionRequest(req)
	if err != nil {
		log.Warn("update subscription validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update, err := toSubscriptionUpdateModel(req)
	if err != nil {
		log.Warn("failed to map update subscription request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Update(r.Context(), id, update)
	if err != nil {
		log.Error("failed to update subscription",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	log.Info("subscription updated successfully", zap.String("subscription_id", id.String()))
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Warn("invalid subscription id", zap.String("subscription_id", idParam), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		log.Error("failed to delete subscription",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	log.Info("subscription deleted successfully", zap.String("subscription_id", id.String()))
}
