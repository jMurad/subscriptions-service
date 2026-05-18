package subscriptions

// SubscriptionResponse subscription response DTO
type SubscriptionResponse struct {

	// Subscription ID
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id"`

	// Subscription service name
	// example: Netflix
	ServiceName string `json:"service_name"`

	// Subscription price in cents
	// example: 999
	Price int `json:"price"`

	// User ID
	// example: 550e8400-e29b-41d4-a716-446655440111
	UserID string `json:"user_id"`

	// Subscription start date
	// example: 2025-01-01
	StartDate string `json:"start_date"`

	// Subscription end date
	// example: 2025-12-31
	EndDate *string `json:"end_date,omitempty"`
}
