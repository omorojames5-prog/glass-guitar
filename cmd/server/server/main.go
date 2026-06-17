package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/omorojames5-prog/glass-guitar/internal/config"
	"github.com/omorojames5-prog/glass-guitar/internal/routes"
	"github.com/omorojames5-prog/glass-guitar/pkg/database"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set Gin mode
	if config.AppConfig.ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	database.ConnectDB()

	// Setup router
	router := routes.SetupRouter()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
