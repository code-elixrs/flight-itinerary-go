package main

import (
	"context"
	"flight-itinerary-go/internal/service"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flight-itinerary-go/internal/handler"
	"flight-itinerary-go/internal/logger"
	customMiddleware "flight-itinerary-go/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "flight-itinerary-go/docs"
)

// @Summary Get health status
// @Description Simple health status api
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health/status [get]
func GetHealthStatus(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "flight-itinerary-go",
	})
}

func main() {
	logger := logger.NewLogger()
	defer logger.Sync()
	logger.Info("Initializing...")

	// Initialize services
	itineraryService := service.NewItineraryService(logger)

	// Initialize handlers
	itineraryHandler := handler.NewItineraryHandler(itineraryService, logger)

	itineraryRequestValidator := customMiddleware.NewItineraryValidator(logger)
	echoServer := echo.New()

	//Global middleware
	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.CORS())
	echoServer.Use(customMiddleware.LoggingMiddleware(logger))

	// Routes
	v1 := echoServer.Group("/api/v1")
	{
		v1.GET("/health/status", GetHealthStatus)
		v1.POST("/itinerary/reconstruct", itineraryHandler.ReconstructItinerary,
			itineraryRequestValidator.Validate())
	}
	echoServer.GET("/swagger/*", echoSwagger.WrapHandler)

	// Graceful shutdown
	go func() {
		if err := echoServer.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := echoServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
