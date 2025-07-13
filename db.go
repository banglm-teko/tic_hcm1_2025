package main

import (
	"database/sql"
	"encoding/json" // Added for JSON handling
	"fmt"
	"log"
	"math/rand" // Added for random activity generation
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var db *sql.DB

// InitDB initializes the MySQL database connection and creates tables if they don't exist
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Successfully connected to MySQL database!")

	// Create tables
	createTableSQL := []string{
		`CREATE TABLE IF NOT EXISTS users (
            user_id INT PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            last_login DATETIME,
            registered_date DATETIME
        );`,
		`CREATE TABLE IF NOT EXISTS products (
            product_id INT PRIMARY KEY,
            product_name VARCHAR(255) NOT NULL,
            category VARCHAR(255) NOT NULL,
            price DECIMAL(10, 2)
        );`,
		`CREATE TABLE IF NOT EXISTS orders (
            order_id INT PRIMARY KEY AUTO_INCREMENT,
            user_id INT,
            product_id INT,
            order_date DATETIME,
            quantity INT,
            total_price DECIMAL(10, 2),
            FOREIGN KEY (user_id) REFERENCES users(user_id),
            FOREIGN KEY (product_id) REFERENCES products(product_id)
        );`,
		`CREATE TABLE IF NOT EXISTS user_preferences (
            user_id INT PRIMARY KEY,
            preferred_categories TEXT,
            churn_risk DECIMAL(3, 2),
            FOREIGN KEY (user_id) REFERENCES users(user_id)
        );`,
		`CREATE TABLE IF NOT EXISTS offers (
            offer_id INT PRIMARY KEY AUTO_INCREMENT,
            user_id INT,
            offer_type VARCHAR(50),
            offer_value VARCHAR(100),
            target_category VARCHAR(255),
            generated_message TEXT,
            sent_date DATETIME,
            is_used BOOLEAN DEFAULT FALSE,
            FOREIGN KEY (user_id) REFERENCES users(user_id)
        );`,
		`CREATE TABLE IF NOT EXISTS user_streaks (
            user_id INT PRIMARY KEY,
            current_streak INT DEFAULT 0,
            longest_streak INT DEFAULT 0,
            last_activity_date DATETIME,
            streak_type VARCHAR(50) DEFAULT 'engagement',
            is_active BOOLEAN DEFAULT TRUE,
            FOREIGN KEY (user_id) REFERENCES users(user_id)
        );`,
		`CREATE TABLE IF NOT EXISTS user_activities (
            activity_id INT PRIMARY KEY AUTO_INCREMENT,
            user_id INT,
            activity_type VARCHAR(50),
            activity_date DATETIME,
            activity_value FLOAT DEFAULT 0,
            FOREIGN KEY (user_id) REFERENCES users(user_id)
        );`,
		`CREATE TABLE IF NOT EXISTS streak_predictions (
            prediction_id INT PRIMARY KEY AUTO_INCREMENT,
            user_id INT,
            prediction_date DATETIME,
            probability_of_streak_drop DECIMAL(5,4),
            predicted_days_to_streak_drop INT,
            risk_level VARCHAR(20),
            confidence DECIMAL(5,4),
            actual_streak_dropped BOOLEAN DEFAULT NULL,
            FOREIGN KEY (user_id) REFERENCES users(user_id)
        );`,
		`CREATE TABLE IF NOT EXISTS streak_models (
            model_id INT PRIMARY KEY AUTO_INCREMENT,
            model_type VARCHAR(100),
            version VARCHAR(50),
            training_date DATETIME,
            accuracy DECIMAL(5,4),
            parameters JSON,
            feature_names JSON,
            is_active BOOLEAN DEFAULT TRUE
        );`,
	}

	for _, sqlStmt := range createTableSQL {
		_, err := db.Exec(sqlStmt)
		if err != nil {
			log.Fatalf("Error creating table: %v\nSQL: %s", err, sqlStmt)
		}
	}
	fmt.Println("Database tables checked/created successfully.")
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
		fmt.Println("Database connection closed.")
	}
}

// InsertSampleData inserts sample data into the database
func InsertSampleData() {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return
	}
	defer tx.Rollback() // Rollback on error

	// Insert user B
	_, err = tx.Exec("INSERT IGNORE INTO users (user_id, username, email, last_login, registered_date) VALUES (?, ?, ?, ?, ?)",
		101, "userB", "userb@example.com", time.Now().Add(-100*24*time.Hour), time.Now().Add(-365*24*time.Hour))
	if err != nil {
		log.Printf("Error inserting user B: %v", err)
		return
	}

	// Insert products
	products := []Product{
		{ProductID: 1, ProductName: "Áo Khoác Nữ Denim", Category: "Thời trang nữ", Price: 450000},
		{ProductID: 2, ProductName: "Váy Hoa Công Sở", Category: "Thời trang nữ", Price: 600000},
		{ProductID: 3, ProductName: "Giày Cao Gót Đen", Category: "Giày dép nữ", Price: 700000},
		{ProductID: 4, ProductName: "Quần Jeans Nam Slim Fit", Category: "Thời trang nam", Price: 500000},
		{ProductID: 5, ProductName: "Tai Nghe Bluetooth", Category: "Điện tử", Price: 1200000},
	}
	for _, p := range products {
		_, err := tx.Exec("INSERT IGNORE INTO products (product_id, product_name, category, price) VALUES (?, ?, ?, ?)",
			p.ProductID, p.ProductName, p.Category, p.Price)
		if err != nil {
			log.Printf("Error inserting product %s: %v", p.ProductName, err)
			return
		}
	}

	// Insert orders for user B
	// Check for existence before inserting to avoid duplicate primary key errors if run multiple times
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = ? AND product_id = ? AND order_date = ?",
		101, 1, time.Now().Add(-120*24*time.Hour)).Scan(&count)
	if err != nil || count == 0 {
		_, err = tx.Exec("INSERT INTO orders (user_id, product_id, order_date, quantity, total_price) VALUES (?, ?, ?, ?, ?)",
			101, 1, time.Now().Add(-120*24*time.Hour), 1, 450000)
		if err != nil {
			log.Printf("Error inserting order 1: %v", err)
			return
		}
	}

	err = tx.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = ? AND product_id = ? AND order_date = ?",
		101, 2, time.Now().Add(-110*24*time.Hour)).Scan(&count)
	if err != nil || count == 0 {
		_, err = tx.Exec("INSERT INTO orders (user_id, product_id, order_date, quantity, total_price) VALUES (?, ?, ?, ?, ?)",
			101, 2, time.Now().Add(-110*24*time.Hour), 1, 600000)
		if err != nil {
			log.Printf("Error inserting order 2: %v", err)
			return
		}
	}

	// Insert user preferences for user B (simulated AI output)
	_, err = tx.Exec("INSERT IGNORE INTO user_preferences (user_id, preferred_categories, churn_risk) VALUES (?, ?, ?)",
		101, "Thời trang nữ", 0.85)
	if err != nil {
		log.Printf("Error inserting user preferences: %v", err)
		return
	}

	// Insert sample user activities for user B
	activityTypes := []string{"login", "browse", "add_to_cart", "purchase", "review"}
	activityValues := []float64{1.0, 2.0, 5.0, 10.0, 3.0}

	// Generate activities over the last 30 days with some gaps to simulate streak breaks
	baseDate := time.Now().Add(-30 * 24 * time.Hour)
	for i := 0; i < 30; i++ {
		// Skip some days to create realistic streak patterns
		if i == 5 || i == 6 || i == 12 || i == 13 || i == 20 || i == 21 || i == 22 {
			continue // Skip weekends and some random days
		}

		activityDate := baseDate.Add(time.Duration(i) * 24 * time.Hour)
		activityType := activityTypes[rand.Intn(len(activityTypes))]
		activityValue := activityValues[rand.Intn(len(activityValues))]

		_, err = tx.Exec("INSERT IGNORE INTO user_activities (user_id, activity_type, activity_date, activity_value) VALUES (?, ?, ?, ?)",
			101, activityType, activityDate, activityValue)
		if err != nil {
			log.Printf("Error inserting activity for day %d: %v", i, err)
			return
		}
	}

	// Insert user streak data for user B
	_, err = tx.Exec("INSERT IGNORE INTO user_streaks (user_id, current_streak, longest_streak, last_activity_date, streak_type, is_active) VALUES (?, ?, ?, ?, ?, ?)",
		101, 3, 15, time.Now().Add(-3*24*time.Hour), "engagement", true)
	if err != nil {
		log.Printf("Error inserting user streak: %v", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}
	fmt.Println("Sample data inserted successfully.")
}

// GetUserData fetches all relevant data for a given user
func GetUserData(userID int) (*UserData, error) {
	userData := &UserData{
		User: User{UserID: userID},
	}

	// Get basic user info
	var lastLogin sql.NullTime // Use sql.NullTime for nullable DATETIME fields
	err := db.QueryRow("SELECT username, email, last_login FROM users WHERE user_id = ?", userID).Scan(
		&userData.Username, &userData.Email, &lastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	} else if err != nil {
		return nil, fmt.Errorf("error fetching user info: %w", err)
	}
	if lastLogin.Valid {
		userData.LastLogin = &lastLogin.Time
	}

	// Get recent orders
	rows, err := db.Query(`
        SELECT o.order_date, p.category, p.product_name
        FROM orders o
        JOIN products p ON o.product_id = p.product_id
        WHERE o.user_id = ?
        ORDER BY o.order_date DESC
        LIMIT 5
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching recent orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderDate time.Time
		var category, productName string
		if err := rows.Scan(&orderDate, &category, &productName); err != nil {
			log.Printf("Error scanning recent order: %v", err)
			continue
		}
		userData.RecentOrders = append(userData.RecentOrders, struct {
			OrderDate   time.Time
			Category    string
			ProductName string
		}{OrderDate: orderDate, Category: category, ProductName: productName})
	}

	// Get user preferences (churn risk and preferred categories from AI)
	var preferredCategoriesStr sql.NullString
	var churnRisk sql.NullFloat64
	err = db.QueryRow("SELECT preferred_categories, churn_risk FROM user_preferences WHERE user_id = ?", userID).Scan(
		&preferredCategoriesStr, &churnRisk,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error fetching user preferences: %w", err)
	}

	if preferredCategoriesStr.Valid {
		userData.PreferredCategories = strings.Split(preferredCategoriesStr.String, ",")
	}
	if churnRisk.Valid {
		userData.ChurnRisk = churnRisk.Float64
	}

	return userData, nil
}

// SaveOffer saves the generated offer to the database
func SaveOffer(offer Offer) error {
	_, err := db.Exec(
		"INSERT INTO offers (user_id, offer_type, offer_value, target_category, generated_message, sent_date, is_used) VALUES (?, ?, ?, ?, ?, ?, ?)",
		offer.UserID, offer.OfferType, offer.OfferValue, offer.TargetCategory, offer.GeneratedMessage, offer.SentDate, offer.IsUsed,
	)
	if err != nil {
		return fmt.Errorf("error saving offer: %w", err)
	}
	fmt.Printf("Offer saved for User ID: %d\n", offer.UserID)
	return nil
}

// GetSavedOffers retrieves offers saved for a specific user (for verification)
func GetSavedOffers(userID int) ([]Offer, error) {
	rows, err := db.Query("SELECT offer_id, user_id, offer_type, offer_value, target_category, generated_message, sent_date, is_used FROM offers WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching saved offers: %w", err)
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var offer Offer
		if err := rows.Scan(
			&offer.OfferID, &offer.UserID, &offer.OfferType, &offer.OfferValue,
			&offer.TargetCategory, &offer.GeneratedMessage, &offer.SentDate, &offer.IsUsed,
		); err != nil {
			log.Printf("Error scanning saved offer: %v", err)
			continue
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

// GetUserStreak retrieves streak information for a user
func GetUserStreak(userID int) (*UserStreak, error) {
	streak := &UserStreak{UserID: userID}

	err := db.QueryRow(`
		SELECT current_streak, longest_streak, last_activity_date, streak_type, is_active
		FROM user_streaks WHERE user_id = ?
	`, userID).Scan(&streak.CurrentStreak, &streak.LongestStreak, &streak.LastActivityDate, &streak.StreakType, &streak.IsActive)

	if err == sql.ErrNoRows {
		return nil, nil // User has no streak data yet
	} else if err != nil {
		return nil, fmt.Errorf("error fetching user streak: %w", err)
	}

	return streak, nil
}

// UpdateUserStreak updates or creates streak data for a user
func UpdateUserStreak(userID int, currentStreak int, lastActivityDate time.Time) error {
	_, err := db.Exec(`
		INSERT INTO user_streaks (user_id, current_streak, longest_streak, last_activity_date, streak_type, is_active)
		VALUES (?, ?, ?, ?, 'engagement', TRUE)
		ON DUPLICATE KEY UPDATE
		current_streak = VALUES(current_streak),
		longest_streak = GREATEST(longest_streak, VALUES(current_streak)),
		last_activity_date = VALUES(last_activity_date),
		is_active = TRUE
	`, userID, currentStreak, currentStreak, lastActivityDate)

	return err
}

// RecordUserActivity records a new user activity
func RecordUserActivity(userID int, activityType string, activityValue float64) error {
	_, err := db.Exec(`
		INSERT INTO user_activities (user_id, activity_type, activity_date, activity_value)
		VALUES (?, ?, NOW(), ?)
	`, userID, activityType, activityValue)

	return err
}

// GetUserActivities retrieves recent activities for a user
func GetUserActivities(userID int, limit int) ([]struct {
	ActivityType  string
	ActivityDate  time.Time
	ActivityValue float64
}, error) {
	rows, err := db.Query(`
		SELECT activity_type, activity_date, activity_value
		FROM user_activities
		WHERE user_id = ?
		ORDER BY activity_date DESC
		LIMIT ?
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []struct {
		ActivityType  string
		ActivityDate  time.Time
		ActivityValue float64
	}

	for rows.Next() {
		var activity struct {
			ActivityType  string
			ActivityDate  time.Time
			ActivityValue float64
		}
		err := rows.Scan(&activity.ActivityType, &activity.ActivityDate, &activity.ActivityValue)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// SaveStreakPrediction saves a prediction to the database
func SaveStreakPrediction(prediction StreakPrediction) error {
	_, err := db.Exec(`
		INSERT INTO streak_predictions 
		(user_id, prediction_date, probability_of_streak_drop, predicted_days_to_streak_drop, risk_level, confidence)
		VALUES (?, NOW(), ?, ?, ?, ?)
	`, prediction.UserID, prediction.ProbabilityOfStreakDrop, prediction.PredictedDaysToStreakDrop,
		prediction.RiskLevel, prediction.Confidence)

	return err
}

// GetStreakPredictions retrieves predictions for a user
func GetStreakPredictions(userID int, limit int) ([]StreakPrediction, error) {
	rows, err := db.Query(`
		SELECT user_id, probability_of_streak_drop, predicted_days_to_streak_drop, 
		       risk_level, confidence, actual_streak_dropped
		FROM streak_predictions
		WHERE user_id = ?
		ORDER BY prediction_date DESC
		LIMIT ?
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []StreakPrediction
	for rows.Next() {
		var pred StreakPrediction
		var actualDropped sql.NullBool
		err := rows.Scan(&pred.UserID, &pred.ProbabilityOfStreakDrop, &pred.PredictedDaysToStreakDrop,
			&pred.RiskLevel, &pred.Confidence, &actualDropped)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, pred)
	}

	return predictions, nil
}

// SaveStreakModel saves a trained model to the database
func SaveStreakModel(model StreakModel) error {
	paramsJSON, err := json.Marshal(model.Parameters)
	if err != nil {
		return err
	}

	featuresJSON, err := json.Marshal(model.FeatureNames)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO streak_models (model_type, version, training_date, accuracy, parameters, feature_names, is_active)
		VALUES (?, ?, ?, ?, ?, ?, TRUE)
	`, model.ModelType, model.Version, model.TrainingDate, model.Accuracy, paramsJSON, featuresJSON)

	return err
}

// GetActiveStreakModel retrieves the currently active model
func GetActiveStreakModel() (*StreakModel, error) {
	var model StreakModel
	var paramsJSON, featuresJSON []byte

	err := db.QueryRow(`
		SELECT model_type, version, training_date, accuracy, parameters, feature_names
		FROM streak_models
		WHERE is_active = TRUE
		ORDER BY training_date DESC
		LIMIT 1
	`).Scan(&model.ModelType, &model.Version, &model.TrainingDate, &model.Accuracy, &paramsJSON, &featuresJSON)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if err := json.Unmarshal(paramsJSON, &model.Parameters); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(featuresJSON, &model.FeatureNames); err != nil {
		return nil, err
	}

	return &model, nil
}
