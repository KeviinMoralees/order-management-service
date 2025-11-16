package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// helloHandler maneja el endpoint /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Solo aceptamos GET a /hello
	if r.URL.Path != "/hello" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Remember, you are the best v1")
}

// createOrderHandler maneja la creación de pedidos
func createOrderHandler(repo *DynamoDBRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Solo aceptamos POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(CreateOrderResponse{
				Success: false,
				Message: "Method not allowed. Use POST",
			})
			return
		}

		// Decodificar el body
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CreateOrderResponse{
				Success: false,
				Message: "Invalid request body",
			})
			return
		}

		// Validar campos requeridos
		if strings.TrimSpace(req.OrderName) == "" || strings.TrimSpace(req.UserName) == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CreateOrderResponse{
				Success: false,
				Message: "orderName and userName are required",
			})
			return
		}

		// Crear el pedido en DynamoDB
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

		// Respuesta exitosa
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateOrderResponse{
			Success: true,
			Message: "Order created successfully",
			OrderID: order.OrderID,
		})
	}
}

