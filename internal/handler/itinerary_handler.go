package handler

import (
	"go.uber.org/zap"
	"net/http"

	"flight-itinerary-go/internal/model"
	"flight-itinerary-go/internal/service"
	"github.com/labstack/echo/v4"
)

// ItineraryHandler handles HTTP requests for itinerary operations
type ItineraryHandler struct {
	itineraryService service.ItineraryService
	logger           *zap.Logger
}

// NewItineraryHandler creates a new itinerary handler
func NewItineraryHandler(itineraryService service.ItineraryService, logger *zap.Logger) *ItineraryHandler {
	return &ItineraryHandler{
		itineraryService: itineraryService,
		logger:           logger,
	}
}

// @Summary Reconstruct Itinerary
// @Description Reconstructs the travel itinerary from a list of source-destination pairs
// @Tags Itinerary
// @Accept json
// @Produce json
// @Param input body []Ticket true "Array of ticket pairs"
// @Success 200 {object} []string
// @Router /api/v1/itinerary/reconstruct [post]
func (itineraryHandlerV1 *ItineraryHandler) ReconstructItinerary(ctx echo.Context) error {
	var request []model.Ticket

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	response, err := itineraryHandlerV1.itineraryService.ReconstructItinerary(request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}
