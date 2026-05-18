package subscriptions

// UpdateRequest represents update subscription payload
type UpdateRequest struct {

	// Subscription service name
	//
	// example: Netflix
	ServiceName *string `json:"service_name"`

	// Subscription monthly price
	//
	// example: 999
	Price *int `json:"price"`

	// Subscription end date
	//
	// example: 2025-12-31
	EndDate *string `json:"end_date"`
}

// UpdateResponse represents update subscription response
type UpdateResponse struct {

	// Response message
	//
	// example: subscription updated
	Message string `json:"message"`
}
