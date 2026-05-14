package dto

// CreateSubscriptionRequest swagger model
// @Description Запрос на создание подписки
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

// CreateSubscriptionResponse swagger model
// @Description Ответ на создание подписки
type CreateSubscriptionResponse struct {
	ID string `json:"id"`
}
