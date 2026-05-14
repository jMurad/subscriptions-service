package dto

type TotalRequest struct {
	From        string `json:"from"`
	To          string `json:"to"`
	UserID      string `json:"user_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
}

type TotalResponse struct {
	Total int `json:"total"`
}
