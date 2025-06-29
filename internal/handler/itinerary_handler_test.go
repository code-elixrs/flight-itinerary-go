package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"flight-itinerary-go/internal/handler"
	"flight-itinerary-go/internal/model"
	"flight-itinerary-go/pkg/errors"
)

func TestItineraryHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ItineraryHandler Suite")
}

// Mock service for testing
type mockItineraryService struct {
	reconstructFunc func([]model.Ticket) ([]string, error)
}

func (m *mockItineraryService) ReconstructItinerary(tickets []model.Ticket) ([]string, error) {
	if m.reconstructFunc != nil {
		return m.reconstructFunc(tickets)
	}
	return []string{"JFK", "LAX"}, nil
}

var _ = Describe("ItineraryHandler", func() {
	var (
		handler1    *handler.ItineraryHandler
		mockService *mockItineraryService
		logger      *zap.Logger
		e           *echo.Echo
	)

	BeforeEach(func() {
		logger = zap.NewExample()
		mockService = &mockItineraryService{}
		handler1 = handler.NewItineraryHandler(mockService, logger)
		e = echo.New()
	})

	Describe("ReconstructItinerary", func() {
		Context("when given valid request", func() {
			It("should return reconstructed itinerary", func() {
				mockService.reconstructFunc = func(tickets []model.Ticket) ([]string, error) {
					return []string{"JFK", "LAX", "DXB"}, nil
				}

				tickets := []model.Ticket{
					{"JFK", "LAX"},
					{"LAX", "DXB"},
				}

				reqBody, _ := json.Marshal(tickets)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				// Set validated_request as []model.Ticket, not *model.ItineraryRequest
				c.Set("validated_request", tickets)

				err := handler1.ReconstructItinerary(c)

				Expect(err).Should(BeNil())
				Expect(rec.Code).To(Equal(http.StatusOK))

				var response []string
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				Expect(err).Should(BeNil())
				Expect(response).To(Equal([]string{"JFK", "LAX", "DXB"}))
			})
		})

		Context("when service returns error", func() {
			It("should return error response", func() {
				mockService.reconstructFunc = func(tickets []model.Ticket) ([]string, error) {
					return nil, errors.ErrNoStartingPoint
				}

				tickets := []model.Ticket{
					{"JFK", "LAX"},
					{"LAX", "DXB"},
				}

				reqBody, _ := json.Marshal(tickets)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				c.Set("validated_request", tickets)

				err := handler1.ReconstructItinerary(c)

				Expect(err).Should(BeNil())
				Expect(rec.Code).To(Equal(http.StatusBadRequest))

				var response errors.AppError
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				Expect(err).Should(BeNil())
				Expect(response.Message).To(Equal("no valid starting point found"))
			})
		})

		Context("when validated request is missing from context", func() {
			It("should return internal server error", func() {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				// Not setting validated_request in context

				err := handler1.ReconstructItinerary(c)

				Expect(err).Should(BeNil())
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("when request contains invalid ticket format", func() {
			It("should handle conversion error", func() {
				// This would typically be caught by middleware, but testing the handler's robustness
				tickets := []model.Ticket{
					{"JFK"}, // Invalid: only one element
				}

				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				c.Set("validated_request", tickets)

				err := handler1.ReconstructItinerary(c)

				Expect(err).Should(BeNil())
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
