package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// getProductsHandler handles requests to /api/products
func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	products, err := GetProducts() // Gọi hàm GetProducts từ db.go
	if err != nil {
		log.Printf("Error getting products: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products) // Mã hóa danh sách sản phẩm thành JSON và gửi về client
}
