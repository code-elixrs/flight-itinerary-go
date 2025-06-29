package model

import "flight-itinerary-go/pkg/errors"

// Ticket represents a flight ticket with source and destination
type Ticket [2]string

// ItineraryResponse represents the response containing the reconstructed itinerary
type ItineraryResponse []string

// Source returns the source airport code
func (t Ticket) Source() string {
	return t[0]
}

// Destination returns the destination airport code
func (t Ticket) Destination() string {
	return t[1]
}

// ItineraryRequest represents the request for itinerary reconstruction
type ItineraryRequest struct {
	Tickets []Ticket `json:"tickets" validate:"required,min=1"`
}

// ToTickets converts string pairs to Ticket domain objects
func (itineraryRequest *ItineraryRequest) ToTickets() ([]Ticket, error) {
	tickets := make([]Ticket, 0, len(itineraryRequest.Tickets))

	for i, ticketPair := range itineraryRequest.Tickets {
		ticket := Ticket{ticketPair[0], ticketPair[1]}
		if err := ticket.Validate(); err != nil {
			return nil, errors.NewValidationError("ticket at index %d is invalid: %v", i, err)
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

// Validate checks if the ticket is valid
func (t Ticket) Validate() error {
	if len(t[0]) == 0 || len(t[1]) == 0 {
		return errors.ErrInvalidTicket
	}
	return nil
}
