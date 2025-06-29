package handler

import (
	"net/http"

	"flight-itinerary-go/internal/model"
	"flight-itinerary-go/pkg/errors"
	"github.com/labstack/echo/v4"
)

func reconstructItinerary(tickets []model.Ticket) ([]string, error) {
	if len(tickets) == 0 {
		return nil, errors.NewValidationError("no tickets provided")
	}

	// Build adjacencyGraph and track destinations
	adjacencyGraph := make(map[string]string)

	for _, ticket := range tickets {
		src, dst := ticket.Source(), ticket.Destination()
		// Check for duplicate routes
		if _, exists := adjacencyGraph[src]; exists {
			return nil, errors.NewValidationError("duplicate route from %s", src)
		}
		adjacencyGraph[src] = dst
	}

	// Find starting point
	startingPoint, err := findStartingPoint(tickets)
	if err != nil {
		return nil, err
	}

	// Reconstruct itinerary
	itinerary, err := buildItinerary(adjacencyGraph, startingPoint, len(tickets))
	if err != nil {
		return nil, err
	}

	// Validate completeness
	if len(itinerary) != len(tickets)+1 {
		return nil, errors.ErrDisconnectedRoute
	}

	return itinerary, nil
}

func buildItinerary(graph map[string]string, startingPoint string, expectedHops int) ([]string, error) {
	itinerary := []string{startingPoint}
	current := startingPoint
	visited := make(map[string]bool)

	for i := 0; i < expectedHops; i++ {
		if visited[current] {
			return nil, errors.ErrCircularRoute
		}
		visited[current] = true
		next, exists := graph[current]
		if !exists {
			break
		}
		itinerary = append(itinerary, next)
		current = next
	}
	return itinerary, nil
}

func findStartingPoint(tickets []model.Ticket) (string, error) {
	destinationSet := make(map[string]bool)
	for _, ticket := range tickets {
		_, dst := ticket.Source(), ticket.Destination()
		destinationSet[dst] = true
	}

	for _, ticket := range tickets {
		if !destinationSet[ticket.Source()] {
			return ticket.Source(), nil
		}
	}
	return "", errors.ErrNoStartingPoint
}

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
	response, err := reconstructItinerary(request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}
