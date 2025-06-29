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
	itineraryService.logger.Info("Starting itinerary reconstruction", zap.Int("ticket_count", len(tickets)))
	if len(tickets) == 0 {
		itineraryService.logger.Warn("Empty ticket list provided")
		return nil, errors.NewValidationError("no tickets provided")
	}

	// Build adjacencyGraph and track destinations
	adjacencyGraph := make(map[string]string)

	for _, ticket := range tickets {
		src, dst := ticket.Source(), ticket.Destination()
		// Check for duplicate routes
		if existing, exists := adjacencyGraph[src]; exists {
			itineraryService.logger.Warn("Duplicate route found", zap.String("source", src),
				zap.String("existing_dest", existing), zap.String("new_dest", dst))
			return nil, errors.NewValidationError("duplicate route from %s", src)
		}
		adjacencyGraph[src] = dst
	}
	itineraryService.logger.Debug("Graph built", zap.Int("nodes", len(adjacencyGraph)))
	// Find starting point
	startingPoint, err := itineraryService.findStartingPoint(tickets)
	if err != nil {
		itineraryService.logger.Error("Failed to find starting point", zap.Error(err))
		return nil, err
	}
	itineraryService.logger.Info("Starting point found", zap.String("start", startingPoint))

	// Reconstruct itinerary
	itinerary, err := itineraryService.buildItinerary(adjacencyGraph, startingPoint, len(tickets))
	if err != nil {
		return nil, err
	}

	// Validate completeness
	if len(itinerary) != len(tickets)+1 {
		itineraryService.logger.Error("Failed to build itinerary as itinerary route disconnected!!")
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
			itineraryService.logger.Warn("Circular route detected", zap.String("city", current),
				zap.Int("step", i))
			return nil, errors.ErrCircularRoute
		}
		visited[current] = true
		next, exists := graph[current]
		if !exists {
			itineraryService.logger.Warn("Route ends prematurely", zap.String("city", current),
				zap.Int("step", i))
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
