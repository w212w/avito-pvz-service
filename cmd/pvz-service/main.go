package main

import (
	"avito-pvz-service/config"
	"avito-pvz-service/internal/handlers"
	"avito-pvz-service/internal/repository"
	"avito-pvz-service/internal/services"
	"avito-pvz-service/internal/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db := storage.ConnectDB(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	router := mux.NewRouter()

	router.HandleFunc("/dummyLogin", authHandler.DummyLogin).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	log.Println("Server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
