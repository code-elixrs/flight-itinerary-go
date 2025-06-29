package handler

import (
	"flight-itinerary-go/pkg/errors"
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
	requestID := ctx.Response().Header().Get(echo.HeaderXRequestID)
	logger := itineraryHandlerV1.logger.With(zap.String("request_id", requestID))

	validatedRequest := ctx.Get("validated_request")
	if validatedRequest == nil {
		logger.Error("Validated request not found in context")
		return itineraryHandlerV1.handleError(ctx, errors.NewInternalError("request validation failed"))
	}

	tickets := validatedRequest.([]model.Ticket)
	request := model.ItineraryRequest{
		Tickets: tickets,
	}

	logger.Info("Processing itinerary reconstruction request",
		zap.Int("ticket_count", len(request.Tickets)))
	tickets, err := request.ToTickets()
	if err != nil {
		logger.Warn("Failed to convert request to tickets", zap.Error(err))
		return itineraryHandlerV1.handleError(ctx, err)
	}

	response, err := itineraryHandlerV1.itineraryService.ReconstructItinerary(tickets)

	if err != nil {
		logger.Error("Failed to reconstruct itinerary", zap.Error(err))
		return itineraryHandlerV1.handleError(ctx, err)
	}
	logger.Info("Successfully reconstructed itinerary",
		zap.Strings("result", response))
	return ctx.JSON(http.StatusOK, response)
}

func (itineraryHandlerV1 *ItineraryHandler) handleError(ctx echo.Context, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return ctx.JSON(appErr.Code, appErr)
	}

	// Handle unexpected errors
	itineraryHandlerV1.logger.Error("Unexpected error", zap.Error(err))
	internalErr := errors.NewInternalError("internal server error")
	return ctx.JSON(internalErr.Code, internalErr)
}
