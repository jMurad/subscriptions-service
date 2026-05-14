package dto

// UpdateSubscriptionRequest swagger model
// @Description Запрос на обновление подписки
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}
