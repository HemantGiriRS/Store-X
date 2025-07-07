package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"

	"storex/database"
	"storex/server"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found or failed to load")
	}

	// Connect to Database and migrate
	if dbError := database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		database.SSLModeDisable); dbError != nil {
		logrus.Fatalf("failed to initialize and migrate database: %v", dbError)
	}
	defer func() {
		if dbError := database.ShutdownDB(); dbError != nil {
			logrus.Panicf("db not properly shut down: %v", dbError)
		}
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := server.SetupRoutes()
	logrus.Infof("Server running on http://localhost:%s", port)
	if serverError := http.ListenAndServe(":"+port, r); serverError != nil {
		logrus.Panicf("failed to start server: %v", serverError)
	}
}
