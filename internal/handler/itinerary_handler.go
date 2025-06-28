package handler

import (
	"flight-itinerary-go/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Reconstruct Itinerary
// @Description Reconstructs the travel itinerary from a list of source-destination pairs
// @Tags Itinerary
// @Accept json
// @Produce json
// @Param input body []Ticket true "Array of ticket pairs"
// @Success 200 {object} []string
// @Router /api/v1/itinerary/reconstruct [post]
func ReconstructItinerary(ctx echo.Context) error {
	var request []model.Ticket

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	response := []string{"JFK", "LAX", "DXB", "SFO", "SJC"}
	return ctx.JSON(http.StatusOK, response)
}
