package service_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"flight-itinerary-go/internal/model"
	"flight-itinerary-go/internal/service"
	"flight-itinerary-go/pkg/errors"
)

func TestItineraryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ItineraryService Suite")
}

var _ = Describe("ItineraryService", func() {
	var (
		itineraryService service.ItineraryService
		logger           *zap.Logger
	)

	BeforeEach(func() {
		logger = zap.NewExample()
		// zapLogger := zap.NewExample()
		itineraryService = service.NewItineraryService(logger)
	})

	Describe("ReconstructItinerary", func() {
		Context("when given valid linear tickets", func() {
			It("should correctly reconstruct the itinerary", func() {
				tickets := []model.Ticket{
					{"JFK", "LAX"},
					{"LAX", "DXB"},
					{"DXB", "SFO"},
					{"SFO", "SJC"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(BeNil())
				Expect(itinerary).To(Equal([]string{"JFK", "LAX", "DXB", "SFO", "SJC"}))
			})
		})

		Context("when given tickets in random order", func() {
			It("should correctly reconstruct the itinerary", func() {
				tickets := []model.Ticket{
					{"DXB", "SFO"},
					{"JFK", "LAX"},
					{"SFO", "SJC"},
					{"LAX", "DXB"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(BeNil())
				Expect(itinerary).To(Equal([]string{"JFK", "LAX", "DXB", "SFO", "SJC"}))
			})
		})

		Context("when given a single ticket", func() {
			It("should return a two-city itinerary", func() {
				tickets := []model.Ticket{
					{"NYC", "LAX"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(BeNil())
				Expect(itinerary).To(Equal([]string{"NYC", "LAX"}))
			})
		})

		Context("when given an empty ticket list", func() {
			It("should return a validation error", func() {
				tickets := []model.Ticket{}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(HaveOccurred())
				Expect(itinerary).Should(BeNil())
				Expect(err.Error()).To(ContainSubstring("no tickets provided"))
			})
		})

		Context("when no valid starting point is found", func() {
			It("should return an error for circular routes", func() {
				tickets := []model.Ticket{
					{"LAX", "DXB"},
					{"DXB", "SFO"},
					{"SFO", "LAX"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(HaveOccurred())
				Expect(itinerary).Should(BeNil())
				Expect(err).To(Equal(errors.ErrNoStartingPoint))
			})
		})

		Context("when given disconnected routes", func() {
			It("should return an error", func() {
				tickets := []model.Ticket{
					{"JFK", "LAX"},
					{"DXB", "SFO"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(HaveOccurred())
				Expect(itinerary).Should(BeNil())
				Expect(err).To(Equal(errors.ErrDisconnectedRoute))
			})
		})

		Context("when given duplicate routes from same source", func() {
			It("should return a validation error", func() {
				tickets := []model.Ticket{
					{"JFK", "LAX"},
					{"JFK", "DXB"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(HaveOccurred())
				Expect(itinerary).Should(BeNil())
				Expect(err.Error()).To(ContainSubstring("duplicate route"))
			})
		})

		Context("when given tickets with circular path detection needed", func() {
			It("should detect and return circular route error", func() {
				tickets := []model.Ticket{
					{"A", "B"},
					{"B", "C"},
					{"C", "A"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(HaveOccurred())
				Expect(itinerary).Should(BeNil())
				Expect(err).To(Equal(errors.ErrNoStartingPoint))
			})
		})

		Context("when given a complex valid route", func() {
			It("should handle multiple stops correctly", func() {
				tickets := []model.Ticket{
					{"BOM", "DEL"},
					{"JFK", "BOM"},
					{"DEL", "BKK"},
					{"BKK", "SIN"},
					{"SIN", "SYD"},
				}

				itinerary, err := itineraryService.ReconstructItinerary(tickets)

				Expect(err).Should(BeNil())
				Expect(itinerary).To(Equal([]string{"JFK", "BOM", "DEL", "BKK", "SIN", "SYD"}))
			})
		})
	})
})
