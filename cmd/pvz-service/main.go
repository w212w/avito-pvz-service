package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	// db := repository.ConnectDB(cfg)
	// defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/auth", authHandler.Auth).Methods("POST")

	log.Println("Server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
