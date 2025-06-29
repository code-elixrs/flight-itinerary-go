package main_test

import (
	"bytes"
	"encoding/json"
	"flight-itinerary-go/internal/service"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"

	"flight-itinerary-go/internal/handler"
	customMiddleware "flight-itinerary-go/internal/middleware"
	"flight-itinerary-go/internal/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = Describe("Integration Tests", func() {
	var (
		echoServer       *echo.Echo
		itineraryService service.ItineraryService
		itineraryHandler *handler.ItineraryHandler
		logger           *zap.Logger
	)

	BeforeEach(func() {
		logger = zap.NewExample()
		echoServer = echo.New()
		itineraryService = service.NewItineraryService(logger)
		itineraryHandler = handler.NewItineraryHandler(itineraryService, logger)
		itineraryRequestValidator := customMiddleware.NewItineraryValidator(logger)
		echoServer.POST("/api/v1/itinerary/reconstruct",
			itineraryHandler.ReconstructItinerary,
			itineraryRequestValidator.Validate(),
		)
	})

	Describe("End-to-End API Tests", func() {
		Context("Itinerary Reconstruction Endpoint", func() {
			It("should successfully reconstruct a simple itinerary", func() {
				request := []model.Ticket{
					{"JFK", "LAX"},
					{"LAX", "DXB"},
				}

				reqBody, _ := json.Marshal(request)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				echoServer.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusOK))

				var response []string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				Expect(err).Should(BeNil())
				Expect(response).To(Equal([]string{"JFK", "LAX", "DXB"}))
			})

			It("should handle complex multi-stop itinerary", func() {
				request := []model.Ticket{
					{"BOM", "DEL"},
					{"JFK", "BOM"},
					{"DEL", "BKK"},
					{"BKK", "SIN"},
					{"SIN", "SYD"},
				}
				reqBody, _ := json.Marshal(request)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				echoServer.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusOK))

				var response model.ItineraryResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				Expect(err).Should(BeNil())
				Expect(response).To(Equal(model.ItineraryResponse{"JFK", "BOM", "DEL", "BKK", "SIN", "SYD"}))
			})

			It("should return error for circular routes", func() {
				request := model.ItineraryRequest{
					Tickets: []model.Ticket{
						{"A", "B"},
						{"B", "C"},
						{"C", "A"},
					},
				}

				reqBody, _ := json.Marshal(request)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				echoServer.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return error for invalid JSON", func() {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader([]byte("invalid json")))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				echoServer.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return error for empty ticket array", func() {
				request := model.ItineraryRequest{
					Tickets: []model.Ticket{},
				}

				reqBody, _ := json.Marshal(request)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				echoServer.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return error for disconnected routes", func() {
				request := model.ItineraryRequest{
					Tickets: []model.Ticket{
						{"JFK", "LAX"},
						{"DXB", "SFO"},
					},
				}

				reqBody, _ := json.Marshal(request)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/itinerary/reconstruct", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				echoServer.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
