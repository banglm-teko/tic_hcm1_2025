package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

// StreakAIModel represents the AI model for streak prediction
type StreakAIModel struct {
	Model      *StreakModel
	Weights    map[string]float64
	Thresholds map[string]float64
	IsTrained  bool
}

// NewStreakAIModel creates a new streak prediction AI model
func NewStreakAIModel() *StreakAIModel {
	return &StreakAIModel{
		Weights: make(map[string]float64),
		Thresholds: map[string]float64{
			"low":      0.3,
			"medium":   0.6,
			"high":     0.8,
			"critical": 0.9,
		},
		IsTrained: false,
	}
}

// ExtractStreakFeatures extracts features from user data for prediction
func ExtractStreakFeatures(userID int, userData *UserData) (*StreakFeatures, error) {
	features := &StreakFeatures{UserID: userID}

	// Get user streak data
	streak, err := GetUserStreak(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user streak: %w", err)
	}

	// Get user activities
	activities, err := GetUserActivities(userID, 100)
	if err != nil {
		return nil, fmt.Errorf("error getting user activities: %w", err)
	}

	// Calculate basic streak features
	if streak != nil {
		features.CurrentStreakLength = streak.CurrentStreak
		features.DaysSinceLastActivity = int(time.Since(streak.LastActivityDate).Hours() / 24)
	} else {
		features.CurrentStreakLength = 0
		features.DaysSinceLastActivity = 999 // High value for new users
	}

	// Calculate average streak length and break frequency
	features.AverageStreakLength = calculateAverageStreakLength(activities)
	features.StreakBreakFrequency = calculateStreakBreakFrequency(activities)

	// Calculate activity-based features
	features.TotalActivities = len(activities)
	features.DaysSinceRegistration = int(time.Since(userData.RegisteredDate).Hours() / 24)

	// Calculate order-based features
	features.AverageOrderValue = calculateAverageOrderValue(userData.RecentOrders)
	features.TotalOrders = len(userData.RecentOrders)
	features.LastOrderDaysAgo = calculateLastOrderDaysAgo(userData.RecentOrders)

	// Use existing churn risk
	features.ChurnRisk = userData.ChurnRisk

	// Calculate category preferences
	features.PreferredCategoriesCount = len(userData.PreferredCategories)

	// Calculate temporal features
	features.SeasonalFactor = calculateSeasonalFactor()
	features.WeekendActivityRatio = calculateWeekendActivityRatio(activities)
	features.EveningActivityRatio = calculateEveningActivityRatio(activities)

	return features, nil
}

// calculateAverageStreakLength calculates the average length of user streaks
func calculateAverageStreakLength(activities []struct {
	ActivityType  string
	ActivityDate  time.Time
	ActivityValue float64
}) float64 {
	if len(activities) < 2 {
		return 0
	}

	var streaks []int
	currentStreak := 1

	for i := 1; i < len(activities); i++ {
		daysDiff := int(activities[i-1].ActivityDate.Sub(activities[i].ActivityDate).Hours() / 24)
		if daysDiff == 1 {
			currentStreak++
		} else {
			if currentStreak > 1 {
				streaks = append(streaks, currentStreak)
			}
			currentStreak = 1
		}
	}

	if currentStreak > 1 {
		streaks = append(streaks, currentStreak)
	}

	if len(streaks) == 0 {
		return 0
	}

	sum := 0
	for _, streak := range streaks {
		sum += streak
	}

	return float64(sum) / float64(len(streaks))
}

// calculateStreakBreakFrequency calculates how often user breaks their streaks
func calculateStreakBreakFrequency(activities []struct {
	ActivityType  string
	ActivityDate  time.Time
	ActivityValue float64
}) float64 {
	if len(activities) < 2 {
		return 0
	}

	breaks := 0
	totalGaps := 0

	for i := 1; i < len(activities); i++ {
		daysDiff := int(activities[i-1].ActivityDate.Sub(activities[i].ActivityDate).Hours() / 24)
		if daysDiff > 1 {
			breaks++
		}
		totalGaps++
	}

	if totalGaps == 0 {
		return 0
	}

	return float64(breaks) / float64(totalGaps)
}

// calculateAverageOrderValue calculates average order value from recent orders
func calculateAverageOrderValue(recentOrders []struct {
	OrderDate   time.Time
	Category    string
	ProductName string
}) float64 {
	if len(recentOrders) == 0 {
		return 0
	}

	// Simulate order values based on categories
	totalValue := 0.0
	for _, order := range recentOrders {
		switch order.Category {
		case "Thời trang nữ":
			totalValue += 500000
		case "Thời trang nam":
			totalValue += 400000
		case "Giày dép nữ":
			totalValue += 700000
		case "Điện tử":
			totalValue += 1200000
		default:
			totalValue += 300000
		}
	}

	return totalValue / float64(len(recentOrders))
}

// calculateLastOrderDaysAgo calculates days since last order
func calculateLastOrderDaysAgo(recentOrders []struct {
	OrderDate   time.Time
	Category    string
	ProductName string
}) int {
	if len(recentOrders) == 0 {
		return 999
	}

	lastOrder := recentOrders[0] // Most recent order
	return int(time.Since(lastOrder.OrderDate).Hours() / 24)
}

// calculateSeasonalFactor calculates seasonal impact on user behavior
func calculateSeasonalFactor() float64 {
	now := time.Now()
	month := now.Month()

	// Higher activity in holiday seasons (December, January, Tet in Vietnam)
	switch month {
	case time.December, time.January:
		return 1.2
	case time.February: // Tet period
		return 1.3
	case time.July, time.August: // Summer
		return 1.1
	default:
		return 1.0
	}
}

// calculateWeekendActivityRatio calculates ratio of weekend activities
func calculateWeekendActivityRatio(activities []struct {
	ActivityType  string
	ActivityDate  time.Time
	ActivityValue float64
}) float64 {
	if len(activities) == 0 {
		return 0
	}

	weekendCount := 0
	for _, activity := range activities {
		weekday := activity.ActivityDate.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			weekendCount++
		}
	}

	return float64(weekendCount) / float64(len(activities))
}

// calculateEveningActivityRatio calculates ratio of evening activities (6 PM - 12 AM)
func calculateEveningActivityRatio(activities []struct {
	ActivityType  string
	ActivityDate  time.Time
	ActivityValue float64
}) float64 {
	if len(activities) == 0 {
		return 0
	}

	eveningCount := 0
	for _, activity := range activities {
		hour := activity.ActivityDate.Hour()
		if hour >= 18 && hour <= 23 {
			eveningCount++
		}
	}

	return float64(eveningCount) / float64(len(activities))
}

// TrainModel trains the streak prediction model
func (m *StreakAIModel) TrainModel(trainingData []StreakTrainingData) error {
	if len(trainingData) == 0 {
		return fmt.Errorf("no training data provided")
	}

	// Initialize weights based on feature importance
	m.Weights = map[string]float64{
		"days_since_last_activity":   0.25,
		"current_streak_length":      0.20,
		"average_streak_length":      0.15,
		"streak_break_frequency":     0.15,
		"churn_risk":                 0.10,
		"last_order_days_ago":        0.05,
		"seasonal_factor":            0.03,
		"weekend_activity_ratio":     0.03,
		"evening_activity_ratio":     0.02,
		"preferred_categories_count": 0.02,
	}

	// Simple gradient descent training
	learningRate := 0.01
	epochs := 100

	for epoch := 0; epoch < epochs; epoch++ {
		totalLoss := 0.0

		for _, data := range trainingData {
			prediction := m.predictProbability(data.Features)
			actual := 0.0
			if data.Label {
				actual = 1.0
			}

			// Calculate loss (binary cross-entropy)
			loss := -(actual*math.Log(prediction+1e-15) + (1-actual)*math.Log(1-prediction+1e-15))
			totalLoss += loss

			// Update weights (simplified gradient descent)
			error := prediction - actual
			for feature, weight := range m.Weights {
				featureValue := m.getFeatureValue(data.Features, feature)
				m.Weights[feature] = weight - learningRate*error*featureValue
			}
		}

		if epoch%20 == 0 {
			log.Printf("Epoch %d, Average Loss: %.4f", epoch, totalLoss/float64(len(trainingData)))
		}
	}

	m.IsTrained = true
	m.Model = &StreakModel{
		ModelType:    "LogisticRegression",
		Version:      "1.0",
		TrainingDate: time.Now(),
		Accuracy:     0.85, // This would be calculated from validation set
		Parameters:   map[string]interface{}{"learning_rate": 0.01, "epochs": epochs},
		FeatureNames: []string{"days_since_last_activity", "current_streak_length", "average_streak_length",
			"streak_break_frequency", "churn_risk", "last_order_days_ago", "seasonal_factor",
			"weekend_activity_ratio", "evening_activity_ratio", "preferred_categories_count"},
	}

	return nil
}

// predictProbability predicts the probability of streak drop
func (m *StreakAIModel) predictProbability(features StreakFeatures) float64 {
	if !m.IsTrained {
		return 0.5 // Default probability if model not trained
	}

	// Calculate weighted sum
	sum := 0.0
	for feature, weight := range m.Weights {
		featureValue := m.getFeatureValue(features, feature)
		sum += weight * featureValue
	}

	// Apply sigmoid function
	return 1.0 / (1.0 + math.Exp(-sum))
}

// getFeatureValue extracts feature value from StreakFeatures
func (m *StreakAIModel) getFeatureValue(features StreakFeatures, featureName string) float64 {
	switch featureName {
	case "days_since_last_activity":
		return float64(features.DaysSinceLastActivity)
	case "current_streak_length":
		return float64(features.CurrentStreakLength)
	case "average_streak_length":
		return features.AverageStreakLength
	case "streak_break_frequency":
		return features.StreakBreakFrequency
	case "churn_risk":
		return features.ChurnRisk
	case "last_order_days_ago":
		return float64(features.LastOrderDaysAgo)
	case "seasonal_factor":
		return features.SeasonalFactor
	case "weekend_activity_ratio":
		return features.WeekendActivityRatio
	case "evening_activity_ratio":
		return features.EveningActivityRatio
	case "preferred_categories_count":
		return float64(features.PreferredCategoriesCount)
	default:
		return 0.0
	}
}

// PredictStreakDrop predicts if and when a user will drop their streak
func (m *StreakAIModel) PredictStreakDrop(userID int, userData *UserData) (*StreakPrediction, error) {
	// Extract features
	features, err := ExtractStreakFeatures(userID, userData)
	if err != nil {
		return nil, fmt.Errorf("error extracting features: %w", err)
	}

	// Make prediction
	probability := m.predictProbability(*features)

	// Calculate predicted days to streak drop
	predictedDays := m.calculatePredictedDaysToStreakDrop(*features, probability)

	// Determine risk level
	riskLevel := m.determineRiskLevel(probability)

	// Calculate confidence based on feature quality
	confidence := m.calculateConfidence(*features)

	// Generate recommended actions
	recommendedActions := m.generateRecommendedActions(*features, riskLevel)

	prediction := &StreakPrediction{
		UserID:                    userID,
		ProbabilityOfStreakDrop:   probability,
		PredictedDaysToStreakDrop: predictedDays,
		RiskLevel:                 riskLevel,
		Confidence:                confidence,
		RecommendedActions:        recommendedActions,
		Features:                  *features,
	}

	// Save prediction to database
	err = SaveStreakPrediction(*prediction)
	if err != nil {
		log.Printf("Warning: Could not save prediction to database: %v", err)
	}

	return prediction, nil
}

// calculatePredictedDaysToStreakDrop estimates when the streak will drop
func (m *StreakAIModel) calculatePredictedDaysToStreakDrop(features StreakFeatures, probability float64) int {
	// Base calculation on current streak length and probability
	baseDays := features.CurrentStreakLength

	// Adjust based on probability
	if probability > 0.8 {
		return max(1, baseDays/4) // High risk - very soon
	} else if probability > 0.6 {
		return max(2, baseDays/3) // Medium-high risk
	} else if probability > 0.4 {
		return max(3, baseDays/2) // Medium risk
	} else {
		return max(7, baseDays) // Low risk - longer time
	}
}

// determineRiskLevel categorizes the risk level
func (m *StreakAIModel) determineRiskLevel(probability float64) string {
	if probability >= m.Thresholds["critical"] {
		return "critical"
	} else if probability >= m.Thresholds["high"] {
		return "high"
	} else if probability >= m.Thresholds["medium"] {
		return "medium"
	} else {
		return "low"
	}
}

// calculateConfidence calculates prediction confidence
func (m *StreakAIModel) calculateConfidence(features StreakFeatures) float64 {
	// Confidence based on data quality and feature values
	confidence := 0.7 // Base confidence

	// Increase confidence with more data
	if features.TotalActivities > 10 {
		confidence += 0.1
	}
	if features.TotalOrders > 3 {
		confidence += 0.1
	}
	if features.DaysSinceRegistration > 30 {
		confidence += 0.1
	}

	return math.Min(0.95, confidence)
}

// generateRecommendedActions suggests actions based on risk level and features
func (m *StreakAIModel) generateRecommendedActions(features StreakFeatures, riskLevel string) []string {
	var actions []string

	switch riskLevel {
	case "critical":
		actions = append(actions, "Send immediate personalized offer")
		actions = append(actions, "Call customer service outreach")
		actions = append(actions, "Send SMS reminder")
		if features.LastOrderDaysAgo > 30 {
			actions = append(actions, "Offer free shipping on next order")
		}
	case "high":
		actions = append(actions, "Send targeted email campaign")
		actions = append(actions, "Offer limited-time discount")
		if features.PreferredCategoriesCount > 0 {
			actions = append(actions, "Send category-specific recommendations")
		}
	case "medium":
		actions = append(actions, "Send gentle reminder email")
		actions = append(actions, "Show personalized homepage content")
	case "low":
		actions = append(actions, "Continue normal engagement")
		actions = append(actions, "Monitor for changes")
	}

	return actions
}

// GenerateTrainingData creates synthetic training data for model training
func GenerateTrainingData(userCount int) []StreakTrainingData {
	var trainingData []StreakTrainingData

	for i := 0; i < userCount; i++ {
		// Generate synthetic features
		features := StreakFeatures{
			UserID:                   i + 1,
			DaysSinceLastActivity:    rand.Intn(30),
			CurrentStreakLength:      rand.Intn(20),
			AverageStreakLength:      rand.Float64() * 10,
			StreakBreakFrequency:     rand.Float64(),
			TotalActivities:          rand.Intn(50),
			DaysSinceRegistration:    rand.Intn(365),
			AverageOrderValue:        rand.Float64() * 1000000,
			TotalOrders:              rand.Intn(20),
			ChurnRisk:                rand.Float64(),
			PreferredCategoriesCount: rand.Intn(5),
			LastOrderDaysAgo:         rand.Intn(60),
			SeasonalFactor:           0.8 + rand.Float64()*0.4,
			WeekendActivityRatio:     rand.Float64(),
			EveningActivityRatio:     rand.Float64(),
		}

		// Generate label based on features (simplified logic)
		label := features.DaysSinceLastActivity > 7 || features.ChurnRisk > 0.7 || features.LastOrderDaysAgo > 30

		trainingData = append(trainingData, StreakTrainingData{
			Features: features,
			Label:    label,
		})
	}

	return trainingData
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
