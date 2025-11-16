package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)


func helloHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/hello" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "test in hello endpoint")
}


func createOrderHandler(repo *DynamoDBRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(CreateOrderResponse{
				Success: false,
				Message: "Method not allowed. Use POST",
			})
			return
		}

		
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CreateOrderResponse{
				Success: false,
				Message: "Invalid request body",
			})
			return
		}

		
		if strings.TrimSpace(req.OrderName) == "" || strings.TrimSpace(req.UserName) == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CreateOrderResponse{
				Success: false,
				Message: "orderName and userName are required",
			})
			return
		}

		
		order, err := repo.CreateOrder(req.OrderName, req.UserName)
		if err != nil {
			log.Printf("Error creating order: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(CreateOrderResponse{
				Success: false,
				Message: "Error creating order",
			})
			return
		}

		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateOrderResponse{
			Success: true,
			Message: "Order created successfully",
			OrderID: order.OrderID,
		})
	}
}

