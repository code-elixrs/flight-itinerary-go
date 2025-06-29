package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"flight-itinerary-go/internal/model"
	"flight-itinerary-go/pkg/errors"
)

type ItineraryRequestValidator interface {
	Validate() echo.MiddlewareFunc
}

type ItineraryRequestValidatorV1 struct {
	logger *zap.Logger
}

// NewItineraryValidator creates a new itinerary handler
func NewItineraryValidator(logger *zap.Logger) ItineraryRequestValidator {
	return &ItineraryRequestValidatorV1{
		logger: logger,
	}
}

// Validate method validates the itinerary reconstruction request
func (itineraryRequestValidatorV1 *ItineraryRequestValidatorV1) Validate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var tickets []model.Ticket

			if err := ctx.Bind(&tickets); err != nil {
				appErr := errors.NewValidationError("invalid JSON format: %v", err)
				return ctx.JSON(appErr.Code, appErr)
			}
			//itineraryRequestValidatorV1.logger.Debug("Received request: ", zap.Any("tickets", tickets))

			if len(tickets) == 0 {
				appErr := errors.NewValidationError("at least one ticket is required")
				return ctx.JSON(appErr.Code, appErr)
			}

			for i, ticket := range tickets {
				if ticket[0] == "" || ticket[1] == "" {
					appErr := errors.NewValidationError("ticket at index %d has empty source or destination", i)
					return ctx.JSON(appErr.Code, appErr)
				}
			}

			// Store validated request in context
			ctx.Set("validated_request", tickets)
			return next(ctx)
		}
	}
}
