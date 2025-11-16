package main

// Order representa un pedido
type Order struct {
	OrderID      string    `json:"orderId" dynamodbav:"orderId"`
	OrderName    string    `json:"orderName" dynamodbav:"orderName"`
	UserName     string    `json:"userName" dynamodbav:"userName"`
	CreatedAt    string    `json:"createdAt" dynamodbav:"createdAt"`
}

// CreateOrderRequest representa la petición para crear un pedido
type CreateOrderRequest struct {
	OrderName string `json:"orderName"`
	UserName  string `json:"userName"`
}

// CreateOrderResponse representa la respuesta al crear un pedido
type CreateOrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	OrderID string `json:"orderId,omitempty"`
}

