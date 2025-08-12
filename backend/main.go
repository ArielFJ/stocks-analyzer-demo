package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"stock-api/internal/api"
	"stock-api/internal/config"
	"stock-api/internal/database"
	"stock-api/internal/middleware"
	"stock-api/internal/services"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	stockService := services.NewStockService(db, cfg.KarenAIToken)
	router := mux.NewRouter()

	api.SetupRoutes(router, stockService)

	// Rate limiting: 100 requests per minute per IP
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(60*time.Second/100), 10)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(middleware.RateLimitMiddleware(rateLimiter)(middleware.LoggingMiddleware(router)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
