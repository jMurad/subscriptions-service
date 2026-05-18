package subscriptions

import (
	service "subscriptions-service/internal/service/subscriptions"
)

type Handler struct {
	svc service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		svc: service,
	}
}
