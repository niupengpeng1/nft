package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dapp/internal/config"
	"dapp/internal/handler"
	"dapp/internal/middleware"
	"dapp/internal/repository"
	"dapp/internal/service"
	"dapp/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.Init()
	logger.Info("Application starting...")

	// Initialize database
	if err := repository.InitDB(&cfg.Database); err != nil {
		logger.Error("Database initialization failed: %v", err)
		os.Exit(1)
	}
	defer repository.CloseDB()

	// Auto migrate database tables
	if err := repository.AutoMigrate(); err != nil {
		logger.Error("Database migration failed: %v", err)
		os.Exit(1)
	}
	logger.Info("Database migration completed")

	// Initialize event listener service
	eventService := service.NewEventListenerService(&cfg.Web3)
	if err := eventService.Init(); err != nil {
		logger.Error("Web3 service initialization failed: %v", err)
		os.Exit(1)
	}
	defer eventService.Close()

	// Initialize HTTP handler
	eventHandler := handler.NewEventHandler(eventService)

	// Create Gin router with default middleware (logger + recovery)
	r := gin.Default()

	// Apply custom middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// Setup API routes
	api := r.Group("/api")
	{
		api.GET("/status", eventHandler.Status)
		api.GET("/events", eventHandler.GetEvents)
		api.POST("/start", eventHandler.StartListening)
		api.POST("/stop", eventHandler.StopListening)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "nft-event-listener",
		})
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "NFT Contract Event Listener",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Start HTTP server in a goroutine
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting HTTP server on %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	eventService.Stop()

	if err := server.Close(); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
	log.Println("Server exiting")
}
