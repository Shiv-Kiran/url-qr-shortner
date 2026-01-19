package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/yourusername/url-qr-shortner/internal/config"
	"github.com/yourusername/url-qr-shortner/internal/database"
	"github.com/yourusername/url-qr-shortner/internal/handlers"
	"github.com/yourusername/url-qr-shortner/internal/services"
)

func main() {
	cfg := config.Load()

	// Initialize database
	if err := database.InitDB(cfg.DBDSN); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize in-memory cache
	services.InitCache()

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/urls", handlers.HandleCreateURL).Methods("POST")
	api.HandleFunc("/urls/{shortCode}", handlers.HandleGetURL).Methods("GET")

	// Redirect route
	router.HandleFunc("/{shortCode}", handlers.HandleRedirect).Methods("GET")
	router.HandleFunc("/", handlers.HandleHome).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		ExposedHeaders:   []string{"X-Total-Count"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", address)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
