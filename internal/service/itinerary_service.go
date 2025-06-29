package service

import (
	"flight-itinerary-go/internal/model"
	"flight-itinerary-go/pkg/errors"
	"go.uber.org/zap"
)

// ItineraryService defines the interface for itinerary operations
type ItineraryService interface {
	ReconstructItinerary(tickets []model.Ticket) ([]string, error)
}

// ItineraryServiceV1 implements the ItineraryService interface
type ItineraryServiceV1 struct {
	logger *zap.Logger
}

// NewItineraryService creates a new instance of ItineraryService
func NewItineraryService(logger *zap.Logger) ItineraryService {
	return &ItineraryServiceV1{
		logger: logger,
	}
}

func (itineraryService *ItineraryServiceV1) ReconstructItinerary(tickets []model.Ticket) ([]string, error) {
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
	startingPoint, err := itineraryService.findStartingPoint(tickets)
	if err != nil {
		return nil, err
	}

	// Reconstruct itinerary
	itinerary, err := itineraryService.buildItinerary(adjacencyGraph, startingPoint, len(tickets))
	if err != nil {
		return nil, err
	}

	// Validate completeness
	if len(itinerary) != len(tickets)+1 {
		return nil, errors.ErrDisconnectedRoute
	}

	return itinerary, nil
}

func (itineraryService *ItineraryServiceV1) buildItinerary(graph map[string]string, startingPoint string,
	expectedHops int) ([]string, error) {
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

func (itineraryService *ItineraryServiceV1) findStartingPoint(tickets []model.Ticket) (string, error) {
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
