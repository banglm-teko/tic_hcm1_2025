package main

import (
	"database/sql"
	"fmt"
	"log"
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
