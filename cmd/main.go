package main

import (
	"fmt"
	"os"
	"vengeful-be/internal/database"
	"vengeful-be/internal/handlers"
	"vengeful-be/internal/middleware"
	"vengeful-be/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel) // Set debug level

	// Log loaded environment variables (safely)
	logger.WithFields(logrus.Fields{
		"MONGODB_DATABASE":        os.Getenv("MONGODB_DATABASE"),
		"MONGODB_HOST":            os.Getenv("MONGODB_HOST"),
		"MONGODB_COLLECTION_NAME": os.Getenv("MONGODB_COLLECTION_NAME"),
		"mongo_username_set":      os.Getenv("MONGO_ROOT_USERNAME") != "",
		"mongo_password_set":      os.Getenv("MONGO_ROOT_PASSWORD") != "",
	}).Debug("Loaded environment variables")

	// Initialize MongoDB connection
	mongodb, err := database.NewMongoDB(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	defer mongodb.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(mongodb.GetCollection(os.Getenv("MONGODB_COLLECTION_NAME")))

	// Initialize router
	router := gin.Default()

	// Configure CORS - Allow all origins
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "User-Agent", "Referer"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Must be false when AllowOrigins is ["*"]
		MaxAge:           12 * 60 * 60,
	}))

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.RateLimiter())
	router.Use(middleware.SecurityHeaders())

	// Initialize handlers
	registerHandler := handlers.NewRegisterHandler(logger, userRepo)
	getUsersHandler := handlers.NewGetUsersHandler(logger, userRepo)

	// Setup routes
	api := router.Group("/api")
	{
		api.POST("/register", registerHandler.Register)
		api.GET("/whosyourdaddy", getUsersHandler.GetAll)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := fmt.Sprintf(":%s", port)
	logger.WithField("port", port).Info("Starting server")
	if err := router.Run(serverAddr); err != nil {
		logger.Fatal("Failed to start server: ", err)
	}
}
