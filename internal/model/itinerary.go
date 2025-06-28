package model

// Ticket represents a flight ticket with source and destination
type Ticket [2]string

// Source returns the source airport code
func (t Ticket) Source() string {
	return t[0]
}

// Destination returns the destination airport code
func (t Ticket) Destination() string {
	return t[1]
}
