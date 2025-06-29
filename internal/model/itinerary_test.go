package model

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestItineraryModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ItineraryModel Suite")
}

var _ = Describe("Ticket", func() {
	Describe("Validate", func() {
		Context("when ticket has valid source and destination", func() {
			It("should return no error", func() {
				ticket := Ticket{"JFK", "LAX"}
				err := ticket.Validate()
				Expect(err).Should(BeNil())
			})
		})

		Context("when ticket has empty source", func() {
			It("should return validation error", func() {
				ticket := Ticket{"", "LAX"}
				err := ticket.Validate()
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("when ticket has empty destination", func() {
			It("should return validation error", func() {
				ticket := Ticket{"JFK", ""}
				err := ticket.Validate()
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("when both source and destination are empty", func() {
			It("should return validation error", func() {
				ticket := Ticket{"", ""}
				err := ticket.Validate()
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Source and Destination methods", func() {
		It("should return correct source and destination", func() {
			ticket := Ticket{"JFK", "LAX"}
			Expect(ticket.Source()).To(Equal("JFK"))
			Expect(ticket.Destination()).To(Equal("LAX"))
		})
	})
})

var _ = Describe("ItineraryRequest", func() {
	Describe("ToTickets", func() {
		Context("when given valid ticket pairs", func() {
			It("should convert to Ticket objects successfully", func() {
				request := &ItineraryRequest{
					Tickets: []Ticket{
						{"JFK", "LAX"},
						{"LAX", "DXB"},
					},
				}

				tickets, err := request.ToTickets()

				Expect(err).Should(BeNil())
				Expect(tickets).To(HaveLen(2))
				Expect(tickets[0]).To(Equal(Ticket{"JFK", "LAX"}))
				Expect(tickets[1]).To(Equal(Ticket{"LAX", "DXB"}))
			})
		})
		Context("when given ticket pairs with empty values", func() {
			It("should return validation error", func() {
				request := &ItineraryRequest{
					Tickets: []Ticket{
						{"JFK", ""},
						{"LAX", "DXB"},
					},
				}

				tickets, err := request.ToTickets()

				Expect(err).Should(HaveOccurred())
				Expect(tickets).Should(BeNil())
				Expect(err.Error()).To(ContainSubstring("is invalid"))
			})
		})

		Context("when given single element ticket pairs", func() {
			It("should return validation error", func() {
				request := &ItineraryRequest{
					Tickets: []Ticket{
						{"JFK"},
					},
				}

				tickets, err := request.ToTickets()

				Expect(err).Should(HaveOccurred())
				Expect(tickets).Should(BeNil())
				Expect(err.Error()).To(ContainSubstring("ticket at index 0 is invalid: invalid ticket: source and destination cannot be empty"))
			})
		})
	})
})
