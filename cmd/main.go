package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	_ "flight-itinerary-go/docs"
)

// @Summary Get health status
// @Description Simple health status api
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health/status [get]
func GetHealthStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "flight-itinerary-go",
	})
}

func main() {
	fmt.Println("Itinerary service!!")
	echoServer := echo.New()
	echoServer.GET("/api/v1/health/status", GetHealthStatus)
	echoServer.GET("/swagger/*", echoSwagger.WrapHandler)
	echoServer.Logger.Fatal(echoServer.Start(":8080"))
}
