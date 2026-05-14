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

// Create godoc
// @Summary Создать подписку
// @Description Создание новой подписки пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionRequest true "Subscription payload"
// @Success 200 {object} dto.CreateSubscriptionResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/ [post]
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

// GetByID godoc
// @Summary Получить подписку по ID
// @Description Возвращает подписку по UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /subscriptions/{id} [get]
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

// List godoc
// @Summary Получить список подписок
// @Description Возвращает список подписок с пагинацией
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/ [get]
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

// Update godoc
// @Summary Обновить подписку
// @Description Частичное обновление подписки
// @Tags subscriptions
// @Accept json
// @Param id path string true "Subscription ID"
// @Param request body dto.UpdateSubscriptionRequest true "Update payload"
// @Success 204
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [patch]
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

// Delete godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по UUID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [delete]
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

	log.Info("subscription deleted successfully",
		zap.String("subscription_id", id.String()),
	)
}

// Total godoc
// @Summary Получить суммарную стоимость подписок
// @Description Считает общую стоимость подписок за период
// @Tags subscriptions
// @Produce json
// @Param from query string true "From month"
// @Param to query string true "To month"
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Success 200 {object} dto.TotalResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) Total(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	req := dto.TotalRequest{
		From:        r.URL.Query().Get("from"),
		To:          r.URL.Query().Get("to"),
		UserID:      r.URL.Query().Get("user_id"),
		ServiceName: r.URL.Query().Get("service_name"),
	}

	if err := validateTotalRequest(req); err != nil {
		log.Warn("total request validation failed",
			zap.Error(err),
			zap.Any("request", req),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filters, err := toTotalFilters(req)
	if err != nil {
		log.Warn("failed to map total request",
			zap.Error(err),
			zap.Any("request", req),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	total, err := h.service.Total(
		r.Context(),
		filters.UserID,
		filters.ServiceName,
		filters.From,
		filters.To,
	)
	if err != nil {
		log.Error("failed to calculate total",
			zap.Error(err),
			zap.Any("user_id", filters.UserID),
			zap.Any("service_name", filters.ServiceName),
			zap.Time("from", filters.From),
			zap.Time("to", filters.To),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	response := toTotalResponse(total)

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("failed to encode total response",
			zap.Error(err),
		)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)

		return
	}

	log.Info("total subscriptions request completed successfully",
		zap.Any("user_id", filters.UserID),
		zap.Any("service_name", filters.ServiceName),
		zap.Time("from", filters.From),
		zap.Time("to", filters.To),
		zap.Int("total", total),
	)
}
