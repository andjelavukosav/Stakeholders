package main

import (
	"context"
	"database-example/controller"
	"database-example/repo"
	"database-example/service"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Logger
	logger := log.New(os.Stdout, "[user-server] ", log.LstdFlags)

	// Kreiramo UserRepository
	userRepo, err := repo.NewUserRepository(logger)
	if err != nil {
		logger.Fatal("Failed to create UserRepository:", err)
	}
	defer func() {
		if err := userRepo.DriverClose(context.Background()); err != nil {
			logger.Println("Error closing driver:", err)
		}
	}()

	// Kreiramo UserService
	userService := &service.UserService{
		UserRepo: userRepo,
	}

	// Kreiramo handler
	userHandler := controller.NewUserHandler(userService)
	adminHandler := controller.NewAdminHandler(userService)

	// Kreiramo router i rute
	router := mux.NewRouter()
	router.HandleFunc("/users", userHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/login", userHandler.Login).Methods("POST", "OPTIONS")
	//router.HandleFunc("/users/{username}", userHandler.Get).Methods("GET", "OPTIONS")
	/*router.HandleFunc("/users/all", userHandler.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}/block", userHandler.BlockUser).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/users/{id}/unblock", userHandler.UnblockUser).Methods("PATCH", "OPTIONS")*/

	//Admin handler
	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(controller.JWTMiddleware) // prvo parsiranje JWT
	adminRoutes.Use(controller.AdminOnly)     // onda provera role
	adminRoutes.HandleFunc("/users/all", adminHandler.GetAllUsers).Methods("GET", "OPTIONS")
	adminRoutes.HandleFunc("/users/block/{id}", adminHandler.BlockUser).Methods("PATCH", "OPTIONS")
	adminRoutes.HandleFunc("/users/unblock/{id}", adminHandler.UnblockUser).Methods("PATCH", "OPTIONS")
	
	// CORS middleware
	corsAllowed := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Start server
	//go func() {
		logger.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", corsAllowed(router)); err != nil {
			logger.Fatal(err)
		}
	//}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	sig := <-sigCh
	logger.Println("Received terminate signal, shutting down:", sig)
}
