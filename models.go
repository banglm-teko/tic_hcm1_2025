package main

import "time"

// User struct represents a row in the users table
type User struct {
	UserID         int        `json:"user_id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	LastLogin      *time.Time `json:"last_login"` // Use pointer for nullable fields
	RegisteredDate time.Time  `json:"registered_date"`
}

// Order struct represents a row in the orders table
type Order struct {
	OrderID    int       `json:"order_id"`
	UserID     int       `json:"user_id"`
	ProductID  int       `json:"product_id"`
	OrderDate  time.Time `json:"order_date"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
}

// Product struct represents a row in the products table
type Product struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
}

// UserPreference struct represents a row in the user_preferences table
type UserPreference struct {
	UserID              int     `json:"user_id"`
	PreferredCategories string  `json:"preferred_categories"` // Comma-separated string
	ChurnRisk           float64 `json:"churn_risk"`
}

// Offer struct represents a row in the offers table
type Offer struct {
	OfferID          int       `json:"offer_id"`
	UserID           int       `json:"user_id"`
	OfferType        string    `json:"offer_type"`
	OfferValue       string    `json:"offer_value"`
	TargetCategory   string    `json:"target_category"`
	GeneratedMessage string    `json:"generated_message"`
	SentDate         time.Time `json:"sent_date"`
	IsUsed           bool      `json:"is_used"`
}

// UserData combines various user-related information for processing
type UserData struct {
	User
	RecentOrders []struct {
		OrderDate   time.Time
		Category    string
		ProductName string
	}
	PreferredCategories []string
	ChurnRisk           float64
}

// OpenAI structures for API request/response
type OpenAIRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
