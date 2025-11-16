package main

import (
	"log"
	"net/http"
)

func main() {
	// Configuración hardcodeada
	tableName := "Orders"
	region := "us-east-1"

	log.Printf("Connecting to DynamoDB table '%s' in region '%s'", tableName, region)

	// Inicializar el repository de DynamoDB
	repo, err := NewDynamoDBRepository(tableName, region)
	if err != nil {
		log.Printf("Warning: Could not initialize DynamoDB: %v", err)
		log.Println("The /orders endpoint will not work properly")
	}

	// Configurar los handlers
	http.HandleFunc("/hello", helloHandler)
	
	if repo != nil {
		http.HandleFunc("/orders", createOrderHandler(repo))
		log.Println("Endpoint /orders configured successfully")
	}

	log.Println("Servidor listening in port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

