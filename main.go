package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Obtener el nombre de la tabla de DynamoDB desde variable de entorno
	// Por defecto usa "Orders" si no está configurada
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "Orders"
	}

	// Inicializar el repository de DynamoDB
	repo, err := NewDynamoDBRepository(tableName)
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

