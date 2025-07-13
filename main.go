//package main
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//	"os"
//	"strings" // Added for DemoStreakAI
//	"time"
//
//	"github.com/joho/godotenv" // For loading .env file
//)
//
//func main() {
//	// Load environment variables from .env file
//	err := godotenv.Load()
//	if err != nil {
//		log.Println("No .env file found or error loading .env. Assuming environment variables are set.")
//	}
//
//	// Check if user wants to run API server
//	if len(os.Args) > 1 && os.Args[1] == "api" {
//		// Run API server
//		port := os.Getenv("API_PORT")
//		if port == "" {
//			port = "8080" // Default port
//		}
//		StartAPIServer(port)
//		return
//	}
//
//	// Get MySQL DSN from environment variable
//	mysqlDSN := os.Getenv("MYSQL_DSN")
//	if mysqlDSN == "" {
//		log.Fatal("MYSQL_DSN environment variable is not set. Please set it in .env or system.")
//	}
//
//	// Initialize DB
//	InitDB(mysqlDSN)
//	defer CloseDB()
//
//	// Insert sample data
//	InsertSampleData()
//
//	// Initialize and train the streak prediction AI model
//	fmt.Println("\n=== Initializing Streak Prediction AI Model ===")
//	streakModel := NewStreakAIModel()
//
//	// Generate training data and train the model
//	fmt.Println("Generating training data...")
//	trainingData := GenerateTrainingData(1000) // Generate 1000 synthetic training samples
//
//	fmt.Println("Training streak prediction model...")
//	err = streakModel.TrainModel(trainingData)
//	if err != nil {
//		log.Printf("Error training streak model: %v", err)
//	} else {
//		fmt.Printf("Model trained successfully! Accuracy: %.2f%%\n", streakModel.Model.Accuracy*100)
//
//		// Save the trained model to database
//		err = SaveStreakModel(*streakModel.Model)
//		if err != nil {
//			log.Printf("Warning: Could not save model to database: %v", err)
//		}
//	}
//
//	// Simulate scanning users to find those at churn risk
//	usersToCheck := []int{101} // Just checking user 101 for this example
//
//	for _, userID := range usersToCheck {
//
//		fmt.Printf("\nProcessing User ID: %d\n", userID)
//		userData, err := GetUserData(userID)
//		if err != nil {
//			log.Printf("Error getting user data for ID %d: %v", userID, err)
//			continue
//		}
//		if userData == nil {
//			fmt.Printf("User with ID %d not found.\n", userID)
//			continue
//		}
//
//		fmt.Printf("  Username: %s\n", userData.Username)
//		if userData.LastLogin != nil {
//			fmt.Printf("  Last Login: %s\n", userData.LastLogin.Format("2006-01-02 15:04:05"))
//		} else {
//			fmt.Println("  Last Login: N/A")
//		}
//		fmt.Printf("  Recent Orders: %v\n", userData.RecentOrders)
//		fmt.Printf("  Preferred Categories: %v\n", userData.PreferredCategories)
//		fmt.Printf("  Churn Risk: %.2f\n", userData.ChurnRisk)
//
//		// --- Streak Prediction Analysis ---
//		fmt.Println("\n  --- Streak Prediction Analysis ---")
//		streakPrediction, err := streakModel.PredictStreakDrop(userID, userData)
//		if err != nil {
//			log.Printf("Error predicting streak drop for user %d: %v", userID, err)
//		} else {
//			fmt.Printf("  Streak Drop Probability: %.2f%%\n", streakPrediction.ProbabilityOfStreakDrop*100)
//			fmt.Printf("  Predicted Days to Streak Drop: %d\n", streakPrediction.PredictedDaysToStreakDrop)
//			fmt.Printf("  Risk Level: %s\n", streakPrediction.RiskLevel)
//			fmt.Printf("  Confidence: %.2f%%\n", streakPrediction.Confidence*100)
//			fmt.Printf("  Recommended Actions:\n")
//			for i, action := range streakPrediction.RecommendedActions {
//				fmt.Printf("    %d. %s\n", i+1, action)
//			}
//
//			// Display key features
//			fmt.Printf("  Key Features:\n")
//			fmt.Printf("    - Days since last activity: %d\n", streakPrediction.Features.DaysSinceLastActivity)
//			fmt.Printf("    - Current streak length: %d\n", streakPrediction.Features.CurrentStreakLength)
//			fmt.Printf("    - Average streak length: %.1f\n", streakPrediction.Features.AverageStreakLength)
//			fmt.Printf("    - Streak break frequency: %.2f\n", streakPrediction.Features.StreakBreakFrequency)
//			fmt.Printf("    - Last order days ago: %d\n", streakPrediction.Features.LastOrderDaysAgo)
//		}
//
//		shouldOffer, preferredCategory := AssessUserForOffer(userData)
//
//		if shouldOffer {
//			offerValue := "25% giảm giá"
//			offerType := "Discount"
//			// In a real scenario, the offer value might be determined by an AI model
//			// based on CLTV and churn risk level.
//
//			fmt.Printf("--> Preparing to generate offer: %s for %s <--\n", offerValue, preferredCategory)
//
//			personalizedMessage, err := GeneratePersonalizedMessageWithLLM(
//				userData.Username,
//				offerValue,
//				preferredCategory,
//			)
//			if err != nil {
//				log.Printf("Error generating personalized message: %v", err)
//				personalizedMessage = "Chúng tôi có một ưu đãi đặc biệt dành cho bạn!" // Fallback message
//			}
//
//			fmt.Println("\n--- Generated Offer ---")
//			fmt.Printf("To: %s\n", userData.Username)
//			fmt.Printf("Content: %s\n", personalizedMessage)
//			fmt.Printf("Offer Type: %s\n", offerType)
//			fmt.Printf("Value: %s\n", offerValue)
//			fmt.Printf("Category: %s\n", preferredCategory)
//			fmt.Println("-----------------------\n")
//
//			// Save the offer to DB
//			offerToSave := Offer{
//				UserID:           userID,
//				OfferType:        offerType,
//				OfferValue:       offerValue,
//				TargetCategory:   preferredCategory,
//				GeneratedMessage: personalizedMessage,
//				SentDate:         time.Now(),
//				IsUsed:           false, // Default to false
//			}
//			err = SaveOffer(offerToSave)
//			if err != nil {
//				log.Printf("Error saving offer: %v", err)
//			}
//		} else {
//			fmt.Printf("No offer generated for %s at this time.\n", userData.Username)
//		}
//	}
//
//	// --- Example: Query saved offers ---
//	fmt.Println("\n--- Querying saved offers from DB ---")
//	savedOffers, err := GetSavedOffers(101)
//	if err != nil {
//		log.Printf("Error fetching saved offers: %v", err)
//	} else if len(savedOffers) > 0 {
//		for _, offer := range savedOffers {
//			fmt.Printf("Offer ID: %d, User ID: %d, Type: %s, Value: %s, Category: %s, Message: %s, Sent: %s, Used: %t\n",
//				offer.OfferID, offer.UserID, offer.OfferType, offer.OfferValue, offer.TargetCategory,
//				offer.GeneratedMessage, offer.SentDate.Format("2006-01-02 15:04:05"), offer.IsUsed)
//		}
//	} else {
//		fmt.Println("No offers found for User ID 101.")
//	}
//
//	// --- Streak Prediction Queries ---
//	fmt.Println("\n--- Querying Streak Predictions from DB ---")
//	streakPredictions, err := GetStreakPredictions(101, 5)
//	if err != nil {
//		log.Printf("Error fetching streak predictions: %v", err)
//	} else if len(streakPredictions) > 0 {
//		for _, pred := range streakPredictions {
//			fmt.Printf("Prediction - User ID: %d, Probability: %.2f%%, Risk: %s, Confidence: %.2f%%\n",
//				pred.UserID, pred.ProbabilityOfStreakDrop*100, pred.RiskLevel, pred.Confidence*100)
//		}
//	} else {
//		fmt.Println("No streak predictions found for User ID 101.")
//	}
//
//	// --- Model Management ---
//	fmt.Println("\n--- Active Streak Model Info ---")
//	activeModel, err := GetActiveStreakModel()
//	if err != nil {
//		log.Printf("Error fetching active model: %v", err)
//	} else if activeModel != nil {
//		fmt.Printf("Model Type: %s\n", activeModel.ModelType)
//		fmt.Printf("Version: %s\n", activeModel.Version)
//		fmt.Printf("Training Date: %s\n", activeModel.TrainingDate.Format("2006-01-02 15:04:05"))
//		fmt.Printf("Accuracy: %.2f%%\n", activeModel.Accuracy*100)
//		fmt.Printf("Features: %v\n", activeModel.FeatureNames)
//	} else {
//		fmt.Println("No active streak model found.")
//	}
//
//	// Run the AI streak prediction demo
//	fmt.Println("\n" + strings.Repeat("=", 60))
//	DemoStreakAI()
//
//	err = godotenv.Load()
//	if err != nil {
//		log.Println("No .env file found or error loading .env. Assuming environment variables are set.")
//	}
//
//	mysqlDSN = os.Getenv("MYSQL_DSN")
//	if mysqlDSN == "" {
//		log.Fatal("MYSQL_DSN environment variable is not set. Please set it in .env or system.")
//	}
//
//	InitDB(mysqlDSN)
//	defer CloseDB()
//
//	// Chèn dữ liệu mẫu (sẽ tự động thêm sản phẩm mới)
//	InsertSampleData()
//
//	// --- BẮT ĐẦU PHẦN API ---
//	// Đăng ký handler từ file product_api.go
//	http.HandleFunc("/api/products", GetProductsHandler)
//	fmt.Println("API server starting on :8080")
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env. Assuming environment variables are set.")
	}

	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		log.Fatal("MYSQL_DSN environment variable is not set. Please set it in .env or system.")
	}

	InitDB(mysqlDSN)
	defer CloseDB()

	// Chèn dữ liệu mẫu (đã bao gồm user B và user A)
	InsertSampleData()

	// --- Đăng ký API Endpoints ---
	http.HandleFunc("/api/products", GetProductsHandler) // API lấy danh sách sản phẩm
	http.HandleFunc("/api/login", LoginHandler)          // API đăng nhập

	fmt.Println("API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
