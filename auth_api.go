package main

import (
	"encoding/json"

	"log"
	"net/http"
	"time"
)

// LoginHandler handles user login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Error retrieving user by email: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// --- Password Verification (In real app, use bcrypt.CompareHashAndPassword) ---
	if user.PasswordHash != req.Password { // Simple comparison for demonstration
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// --- Login Successful ---
	response := LoginResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Token:    "dummy_jwt_token", // In a real app, generate a JWT
		Message:  "Login successful",
	}

	// --- Streak Check & Offer Notification Logic ---
	// Check if this user is "user_b@example.com"
	if user.Email == "user_b@example.com" {
		// Simulate check if user B has broken streak (e.g., last login was long ago)
		// For simplicity, we assume user_b@example.com ALWAYS receives this notification
		// upon successful login if their last_login was before a certain threshold or if they
		// are marked as high churn risk in user_preferences.

		// Let's use the churn risk logic directly here for user_b
		userData, err := GetUserData(user.UserID)
		if err != nil {
			log.Printf("Error getting user data for streak check: %v", err)
			// Proceed without offer if data retrieval fails
		} else {
			shouldOffer, targetCategory := AssessUserForOffer(userData) // Reuse the AI assessment logic
			if shouldOffer {
				offerValue := "25% giảm giá"
				// Generate personalized message for the offer
				personalizedMessage, err := GeneratePersonalizedMessageWithLLM(
					user.Username,
					offerValue,
					targetCategory,
				)
				if err != nil {
					log.Printf("Error generating LLM message for user %s: %v", user.Email, err)
					personalizedMessage = "Bạn có một ưu đãi đặc biệt đang chờ!"
				}

				// Save the offer to DB
				offerToSave := Offer{
					UserID:           user.UserID,
					OfferType:        "Discount",
					OfferValue:       offerValue,
					TargetCategory:   targetCategory,
					GeneratedMessage: personalizedMessage,
					SentDate:         time.Now(),
					IsUsed:           false,
				}
				saveErr := SaveOffer(offerToSave)
				if saveErr != nil {
					log.Printf("Error saving offer for user %s: %v", user.Email, saveErr)
				} else {
					// Prepare push notification for frontend
					// In a real system, this would trigger an actual push notification service (e.g., Firebase Cloud Messaging)
					response.OfferNotification = &OfferNotification{
						Title:   "Ưu đãi đặc biệt dành cho bạn! 🎉",
						Message: personalizedMessage,
						OfferID: 0, // In real app, get actual offer ID after saving
					}
					log.Printf("Push notification prepared for %s: %s", user.Email, personalizedMessage)
				}
			}
		}
	}

	// Update last login time
	err = UpdateUserLastLogin(user.UserID, time.Now())
	if err != nil {
		log.Printf("Error updating last login for user %d: %v", user.UserID, err)
		// This is not a critical error, so we proceed with login response
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
