package main

import (
	"fmt"
	"log"
	"net/http"
)

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

func main() {
	http.HandleFunc("/hello", helloHandler)

	log.Println("Servidor listening in port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

