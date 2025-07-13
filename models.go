package main

import "time"

// LoginRequest represents the login API request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login API response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int    `json:"user_id,omitempty"`
}

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

// UserStreak represents user engagement streak data
type UserStreak struct {
	UserID           int       `json:"user_id"`
	CurrentStreak    int       `json:"current_streak"`
	LongestStreak    int       `json:"longest_streak"`
	LastActivityDate time.Time `json:"last_activity_date"`
	StreakType       string    `json:"streak_type"` // "login", "purchase", "engagement"
	IsActive         bool      `json:"is_active"`
}

// StreakFeatures contains features for ML model prediction
type StreakFeatures struct {
	UserID                   int     `json:"user_id"`
	DaysSinceLastActivity    int     `json:"days_since_last_activity"`
	CurrentStreakLength      int     `json:"current_streak_length"`
	AverageStreakLength      float64 `json:"average_streak_length"`
	StreakBreakFrequency     float64 `json:"streak_break_frequency"`
	TotalActivities          int     `json:"total_activities"`
	DaysSinceRegistration    int     `json:"days_since_registration"`
	AverageOrderValue        float64 `json:"average_order_value"`
	TotalOrders              int     `json:"total_orders"`
	ChurnRisk                float64 `json:"churn_risk"`
	PreferredCategoriesCount int     `json:"preferred_categories_count"`
	LastOrderDaysAgo         int     `json:"last_order_days_ago"`
	SeasonalFactor           float64 `json:"seasonal_factor"`
	WeekendActivityRatio     float64 `json:"weekend_activity_ratio"`
	EveningActivityRatio     float64 `json:"evening_activity_ratio"`
}

// StreakPrediction contains the model's prediction output
type StreakPrediction struct {
	UserID                    int            `json:"user_id"`
	ProbabilityOfStreakDrop   float64        `json:"probability_of_streak_drop"`
	PredictedDaysToStreakDrop int            `json:"predicted_days_to_streak_drop"`
	RiskLevel                 string         `json:"risk_level"` // "low", "medium", "high", "critical"
	Confidence                float64        `json:"confidence"`
	RecommendedActions        []string       `json:"recommended_actions"`
	Features                  StreakFeatures `json:"features"`
}

// StreakModel represents the trained ML model
type StreakModel struct {
	ModelType    string                 `json:"model_type"`
	Version      string                 `json:"version"`
	TrainingDate time.Time              `json:"training_date"`
	Accuracy     float64                `json:"accuracy"`
	Parameters   map[string]interface{} `json:"parameters"`
	FeatureNames []string               `json:"feature_names"`
}

// StreakTrainingData represents training data for the model
type StreakTrainingData struct {
	Features StreakFeatures `json:"features"`
	Label    bool           `json:"label"` // true if streak was broken within prediction window
}
