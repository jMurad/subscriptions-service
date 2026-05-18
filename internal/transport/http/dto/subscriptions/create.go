package subscriptions

// CreateRequest represents create subscription payload
type CreateRequest struct {

	// Subscription service name
	//
	// required: true
	// example: Netflix
	ServiceName string `json:"service_name"`

	// Subscription monthly price
	//
	// required: true
	// example: 999
	Price int `json:"price"`

	// User UUID
	//
	// required: true
	// example: 550e8400-e29b-41d4-a716-446655440000
	UserID string `json:"user_id"`

	// Subscription start date
	//
	// required: true
	// example: 2025-01-01
	StartDate string `json:"start_date"`

	// Subscription end date
	//
	// example: 2025-12-31
	EndDate *string `json:"end_date"`
}

// CreateResponse represents create subscription response
type CreateResponse struct {

	// Created subscription ID
	//
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id"`
}
