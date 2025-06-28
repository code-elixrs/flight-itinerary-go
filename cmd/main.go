package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetHealthStatus(context echo.Context) error {
	return context.JSON(http.StatusOK, "Healthy")
}

func main() {
	fmt.Println("Itinerary service!!")
	echoServer := echo.New()
	echoServer.GET("/api/v1/health/status", GetHealthStatus)
	echoServer.Logger.Fatal(echoServer.Start(":8080"))
}
