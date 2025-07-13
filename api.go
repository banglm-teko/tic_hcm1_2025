package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Hardcoded credentials for demo purposes
var hardcodedUsers = map[string]string{
	"admin":    "admin123",
	"user1":    "password1",
	"demo":     "demo123",
	"testuser": "testpass",
	"john_doe": "john123",
}

// userIDs maps usernames to their corresponding user IDs
var userIDs = map[string]int{
	"admin":    1,
	"user1":    101,
	"demo":     102,
	"testuser": 103,
	"john_doe": 104,
}

// StartAPIServer starts the HTTP server with login endpoints
func StartAPIServer(port string) {
	// Initialize database connection if MYSQL_DSN is available
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN != "" {
		fmt.Println("Initializing database connection...")
		InitDB(mysqlDSN)
		defer CloseDB()
		fmt.Println("Database connection established.")
	} else {
		fmt.Println("Warning: MYSQL_DSN not set. Database features will be disabled.")
	}

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/login", handleLogin).Methods("POST")
	api.HandleFunc("/health", handleHealth).Methods("GET")

	// CORS middleware
	r.Use(corsMiddleware)

	fmt.Printf("API Server starting on port %s...\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /api/login - Login with username/password")
	fmt.Println("  GET  /api/health - Health check")
	fmt.Println("\nHardcoded users:")
	for username := range hardcodedUsers {
		fmt.Printf("  Username: %s, Password: %s\n", username, hardcodedUsers[username])
	}
	fmt.Println()

	log.Fatal(http.ListenAndServe(":"+port, r))
}

// handleLogin processes login requests
func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, `{"success": false, "message": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if loginReq.Username == "" || loginReq.Password == "" {
		http.Error(w, `{"success": false, "message": "Username and password are required"}`, http.StatusBadRequest)
		return
	}

	// Check hardcoded credentials
	if password, exists := hardcodedUsers[loginReq.Username]; exists && password == loginReq.Password {
		// Login successful
		userID := userIDs[loginReq.Username]

		// Update last login time in database if user exists
		if userID > 0 {
			UpdateUserLastLogin(userID)
		}

		response := LoginResponse{
			Success: true,
			Message: "Login successful",
			UserID:  userID,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		log.Printf("Successful login for user: %s (ID: %d)", loginReq.Username, userID)
	} else {
		// Login failed
		response := LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)

		log.Printf("Failed login attempt for user: %s", loginReq.Username)
	}
}

// handleHealth provides a simple health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "TIC HCM1 2025 API",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// corsMiddleware adds CORS headers to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UpdateUserLastLogin updates the last login time for a user in the database
func UpdateUserLastLogin(userID int) {
	// Check if database connection is available
	if db == nil {
		log.Printf("Database not available, skipping last login update for user %d", userID)
		return
	}

	query := "UPDATE users SET last_login = NOW() WHERE user_id = ?"

	_, err := db.Exec(query, userID)
	if err != nil {
		log.Printf("Error updating last login for user %d: %v", userID, err)
	} else {
		log.Printf("Updated last login time for user %d", userID)
	}
}
