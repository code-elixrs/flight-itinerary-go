package main_test

import (
	"bytes"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"

	"flight-itinerary-go/internal/handler"
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
		echoServer *echo.Echo
	)

	BeforeEach(func() {
		echoServer = echo.New()
		echoServer.POST("/api/v1/itinerary/reconstruct",
			handler.ReconstructItinerary,
		)
	})

	Describe("End-to-End API Tests", func() {
		Context("Itinerary Reconstruction Endpoint", func() {
			It("should successfully reconstruct a simple itinerary", func() {
				request := []model.Ticket{
					{"LAX", "DXB"},
					{"JFK", "LAX"},
					{"SFO", "SJC"},
					{"DXB", "SFO"},
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
				Expect(response).To(Equal([]string{"JFK", "LAX", "DXB", "SFO", "SJC"}))
			})
		})
	})
})
