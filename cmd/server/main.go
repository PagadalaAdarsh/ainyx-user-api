package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	db "github.com/adarsh/ainyx-task/db/sqlc"
	"github.com/adarsh/ainyx-task/internal/handler"
	"github.com/adarsh/ainyx-task/internal/logger"
	"github.com/adarsh/ainyx-task/internal/service"
	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	defer logger.Log.Sync() // flush logs

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn(".env file not found, using system environment variables")
	}

	// Read DB environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect using database/sql with pgx driver
	dbConn, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Fatal("DB connection error", zap.Error(err))
	}
	defer dbConn.Close()

	// Ping DB
	if err := dbConn.Ping(); err != nil {
		logger.Log.Fatal("DB ping error", zap.Error(err))
	}
	logger.Log.Info("Connected to PostgreSQL")

	// Initialize SQLC queries
	store := db.New(dbConn)

	// Initialize service and handler
	userService := service.NewUserService(store)
	userHandler := handler.NewUserHandler(userService)

	// Initialize Fiber app
	app := fiber.New()

	// Middleware to log request info and duration
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		logger.Log.Info("Request completed",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
		)
		return err
	})

	// Setup routes
	api := app.Group("/users")
	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUser)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Log.Info("Server running", zap.String("port", port))
	log.Fatal(app.Listen(":" + port))
}
