package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv" // For loading .env file
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env. Assuming environment variables are set.")
	}

	// Get MySQL DSN from environment variable
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		log.Fatal("MYSQL_DSN environment variable is not set. Please set it in .env or system.")
	}

	// Initialize DB
	InitDB(mysqlDSN)
	defer CloseDB()

	// Insert sample data
	InsertSampleData()

	// Simulate scanning users to find those at churn risk
	usersToCheck := []int{101} // Just checking user 101 for this example

	for _, userID := range usersToCheck {
		fmt.Printf("\nProcessing User ID: %d\n", userID)
		userData, err := GetUserData(userID)
		if err != nil {
			log.Printf("Error getting user data for ID %d: %v", userID, err)
			continue
		}
		if userData == nil {
			fmt.Printf("User with ID %d not found.\n", userID)
			continue
		}

		fmt.Printf("  Username: %s\n", userData.Username)
		if userData.LastLogin != nil {
			fmt.Printf("  Last Login: %s\n", userData.LastLogin.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("  Last Login: N/A")
		}
		fmt.Printf("  Recent Orders: %v\n", userData.RecentOrders)
		fmt.Printf("  Preferred Categories: %v\n", userData.PreferredCategories)
		fmt.Printf("  Churn Risk: %.2f\n", userData.ChurnRisk)

		shouldOffer, preferredCategory := AssessUserForOffer(userData)

		if shouldOffer {
			offerValue := "25% giảm giá"
			offerType := "Discount"
			// In a real scenario, the offer value might be determined by an AI model
			// based on CLTV and churn risk level.

			fmt.Printf("--> Preparing to generate offer: %s for %s <--\n", offerValue, preferredCategory)

			personalizedMessage, err := GeneratePersonalizedMessageWithLLM(
				userData.Username,
				offerValue,
				preferredCategory,
			)
			if err != nil {
				log.Printf("Error generating personalized message: %v", err)
				personalizedMessage = "Chúng tôi có một ưu đãi đặc biệt dành cho bạn!" // Fallback message
			}

			fmt.Println("\n--- Generated Offer ---")
			fmt.Printf("To: %s\n", userData.Username)
			fmt.Printf("Content: %s\n", personalizedMessage)
			fmt.Printf("Offer Type: %s\n", offerType)
			fmt.Printf("Value: %s\n", offerValue)
			fmt.Printf("Category: %s\n", preferredCategory)
			fmt.Println("-----------------------\n")

			// Save the offer to DB
			offerToSave := Offer{
				UserID:           userID,
				OfferType:        offerType,
				OfferValue:       offerValue,
				TargetCategory:   preferredCategory,
				GeneratedMessage: personalizedMessage,
				SentDate:         time.Now(),
				IsUsed:           false, // Default to false
			}
			err = SaveOffer(offerToSave)
			if err != nil {
				log.Printf("Error saving offer: %v", err)
			}
		} else {
			fmt.Printf("No offer generated for %s at this time.\n", userData.Username)
		}
	}

	// --- Example: Query saved offers ---
	fmt.Println("\n--- Querying saved offers from DB ---")
	savedOffers, err := GetSavedOffers(101)
	if err != nil {
		log.Printf("Error fetching saved offers: %v", err)
	} else if len(savedOffers) > 0 {
		for _, offer := range savedOffers {
			fmt.Printf("Offer ID: %d, User ID: %d, Type: %s, Value: %s, Category: %s, Message: %s, Sent: %s, Used: %t\n",
				offer.OfferID, offer.UserID, offer.OfferType, offer.OfferValue, offer.TargetCategory,
				offer.GeneratedMessage, offer.SentDate.Format("2006-01-02 15:04:05"), offer.IsUsed)
		}
	} else {
		fmt.Println("No offers found for User ID 101.")
	}

	err = godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env. Assuming environment variables are set.")
	}

	mysqlDSN = os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		log.Fatal("MYSQL_DSN environment variable is not set. Please set it in .env or system.")
	}

	InitDB(mysqlDSN)
	defer CloseDB()

	// Chèn dữ liệu mẫu (sẽ tự động thêm sản phẩm mới)
	InsertSampleData()

	// --- BẮT ĐẦU PHẦN API ---
	// Đăng ký handler từ file product_api.go
	http.HandleFunc("/api/products", GetProductsHandler)
	fmt.Println("API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
