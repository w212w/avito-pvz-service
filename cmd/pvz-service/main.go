package main

import (
	"avito-pvz-service/config"
	pb "avito-pvz-service/internal/grpc/pvz_v1"
	"avito-pvz-service/internal/handlers"
	"avito-pvz-service/internal/middleware"
	"avito-pvz-service/internal/repository"
	"avito-pvz-service/internal/services"
	"avito-pvz-service/internal/storage"
	logger "avito-pvz-service/pkg"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.LogLevel)

	db := storage.ConnectDB(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	pvzRepo := repository.NewPVZRepo(db)
	pvzService := services.NewPVZService(pvzRepo)
	pvzHandler := handlers.NewPVZHandler(pvzService)

	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/dummyLogin", authHandler.DummyLogin).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	router.Handle("/pvz", middleware.AuthMiddleware(cfg.JWTSecret, services.RoleModerator)(http.HandlerFunc(pvzHandler.CreatePVZ))).Methods("POST")
	router.Handle("/pvz", middleware.AuthMiddleware(cfg.JWTSecret, services.RoleModerator, services.RoleEmployee)(http.HandlerFunc(pvzHandler.GetPVZList))).Methods("GET")
	router.Handle("/receptions", middleware.AuthMiddleware(cfg.JWTSecret, services.RoleEmployee)(http.HandlerFunc(pvzHandler.CreateReception))).Methods("POST")
	router.Handle("/products", middleware.AuthMiddleware(cfg.JWTSecret, services.RoleEmployee)(http.HandlerFunc(pvzHandler.AddProduct))).Methods("POST")
	router.Handle("/pvz/{pvzId}/close_last_reception", middleware.AuthMiddleware(cfg.JWTSecret, services.RoleEmployee)(http.HandlerFunc(pvzHandler.CloseLastReception))).Methods("POST")
	router.Handle("/pvz/{pvzId}/delete_last_product", middleware.AuthMiddleware(cfg.JWTSecret, services.RoleEmployee)(http.HandlerFunc(pvzHandler.DeleteLastProduct))).Methods("POST")

	go func() {
		logger.Log.Info("HTTP Server is starting on :8080")
		if err := http.ListenAndServe(":8080", router); err != nil {
			logger.Log.WithError(err).Fatal("HTTP Server failed")
		}
	}()

	go func() {
		
		lis, err := net.Listen("tcp", ":3000")
		if err != nil {
			logger.Log.WithError(err).Fatal("Listening on :3000 failed")
		}

		grpcServer := grpc.NewServer()

		pb.RegisterPVZServiceServer(grpcServer, pb.NewGRPCServer(pvzService))

		logger.Log.Info("gRPC server listening on :3000")

		if err := grpcServer.Serve(lis); err != nil {
			logger.Log.WithError(err).Fatal("GRPC Server failed")
		}
	}()

	select {}

}
