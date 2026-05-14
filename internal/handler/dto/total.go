package dto

// TotalRequest swagger model
// @Description Запрос суммарной стоимости подписок
type TotalRequest struct {
	From        string `json:"from"`
	To          string `json:"to"`
	UserID      string `json:"user_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
}

// TotalResponse swagger model
// @Description Ответ на запрос суммарной стоимости подписок
type TotalResponse struct {
	Total int `json:"total"`
}
