package main

import (
	"fmt"
	"log"
	"time"
)

// DemoStreakAI demonstrates the AI streak prediction system
func DemoStreakAI() {
	fmt.Println("=== AI Streak Drop Prediction System Demo ===\n")

	// Initialize the AI model
	fmt.Println("1. Initializing AI Model...")
	streakModel := NewStreakAIModel()
	fmt.Println("   ✓ Model initialized successfully")

	// Generate training data
	fmt.Println("\n2. Generating Training Data...")
	trainingData := GenerateTrainingData(500) // Smaller dataset for demo
	fmt.Printf("   ✓ Generated %d training samples\n", len(trainingData))

	// Train the model
	fmt.Println("\n3. Training Model...")
	err := streakModel.TrainModel(trainingData)
	if err != nil {
		log.Printf("Training failed: %v", err)
		return
	}
	fmt.Printf("   ✓ Model trained successfully (Accuracy: %.1f%%)\n", streakModel.Model.Accuracy*100)

	// Create sample user data for prediction
	fmt.Println("\n4. Creating Sample User Data...")
	sampleUserData := &UserData{
		User: User{
			UserID:         101,
			Username:       "demo_user",
			Email:          "demo@example.com",
			RegisteredDate: time.Now().Add(-180 * 24 * time.Hour), // 6 months ago
		},
		RecentOrders: []struct {
			OrderDate   time.Time
			Category    string
			ProductName string
		}{
			{OrderDate: time.Now().Add(-5 * 24 * time.Hour), Category: "Thời trang nữ", ProductName: "Áo Khoác"},
			{OrderDate: time.Now().Add(-15 * 24 * time.Hour), Category: "Thời trang nữ", ProductName: "Váy Hoa"},
			{OrderDate: time.Now().Add(-30 * 24 * time.Hour), Category: "Giày dép nữ", ProductName: "Giày Cao Gót"},
		},
		PreferredCategories: []string{"Thời trang nữ", "Giày dép nữ"},
		ChurnRisk:           0.75, // High churn risk
	}
	fmt.Println("   ✓ Sample user data created")

	// Extract features
	fmt.Println("\n5. Extracting Features...")
	features := &StreakFeatures{
		UserID:                   101,
		DaysSinceLastActivity:    8, // User hasn't been active for 8 days
		CurrentStreakLength:      3, // Current streak is 3 days
		AverageStreakLength:      7.5,
		StreakBreakFrequency:     0.3, // 30% of the time they break streaks
		TotalActivities:          25,
		DaysSinceRegistration:    180,
		AverageOrderValue:        550000,
		TotalOrders:              3,
		ChurnRisk:                0.75,
		PreferredCategoriesCount: 2,
		LastOrderDaysAgo:         5,
		SeasonalFactor:           1.0,
		WeekendActivityRatio:     0.4,
		EveningActivityRatio:     0.6,
	}
	fmt.Println("   ✓ Features extracted")

	// Make prediction
	fmt.Println("\n6. Making Prediction...")
	probability := streakModel.predictProbability(*features)
	predictedDays := streakModel.calculatePredictedDaysToStreakDrop(*features, probability)
	riskLevel := streakModel.determineRiskLevel(probability)
	confidence := streakModel.calculateConfidence(*features)
	recommendedActions := streakModel.generateRecommendedActions(*features, riskLevel)

	// Display results
	fmt.Println("\n=== PREDICTION RESULTS ===")
	fmt.Printf("User ID: %d\n", features.UserID)
	fmt.Printf("Username: %s\n", sampleUserData.Username)
	fmt.Printf("Streak Drop Probability: %.1f%%\n", probability*100)
	fmt.Printf("Predicted Days to Streak Drop: %d\n", predictedDays)
	fmt.Printf("Risk Level: %s\n", riskLevel)
	fmt.Printf("Confidence: %.1f%%\n", confidence*100)

	fmt.Println("\nKey Features:")
	fmt.Printf("  • Days since last activity: %d\n", features.DaysSinceLastActivity)
	fmt.Printf("  • Current streak length: %d\n", features.CurrentStreakLength)
	fmt.Printf("  • Average streak length: %.1f\n", features.AverageStreakLength)
	fmt.Printf("  • Streak break frequency: %.1f%%\n", features.StreakBreakFrequency*100)
	fmt.Printf("  • Last order days ago: %d\n", features.LastOrderDaysAgo)
	fmt.Printf("  • Churn risk: %.1f%%\n", features.ChurnRisk*100)

	fmt.Println("\nRecommended Actions:")
	for i, action := range recommendedActions {
		fmt.Printf("  %d. %s\n", i+1, action)
	}

	// Feature importance analysis
	fmt.Println("\n=== FEATURE IMPORTANCE ===")
	fmt.Println("Model weights (higher = more important):")
	for feature, weight := range streakModel.Weights {
		fmt.Printf("  • %s: %.3f\n", feature, weight)
	}

	// Model information
	fmt.Println("\n=== MODEL INFORMATION ===")
	fmt.Printf("Model Type: %s\n", streakModel.Model.ModelType)
	fmt.Printf("Version: %s\n", streakModel.Model.Version)
	fmt.Printf("Training Date: %s\n", streakModel.Model.TrainingDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Accuracy: %.1f%%\n", streakModel.Model.Accuracy*100)
	fmt.Printf("Training Samples: %d\n", len(trainingData))

	fmt.Println("\n=== DEMO COMPLETED ===")
	fmt.Println("This demonstrates a complete AI system for predicting user streak drops.")
	fmt.Println("The model analyzes user behavior patterns and provides actionable insights.")
}
