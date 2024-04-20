package ratings

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type mockService struct {
	getallrating func(*[]models.RatingResponse) error
	createrating func(*models.Reviews) error
}

func (h *mockService) GetAllRatings(rating *[]models.RatingResponse) error {
	return h.getallrating(rating)
}
func (h *mockService) CreateRating(rating *models.Reviews) error {
	return h.createrating(rating)
}

func TestCreateRating(t *testing.T) {
	app := fiber.New()

	testCase := []struct {
		description  string
		requestbody  string
		expectStatus int
		cookieValue  string
	}{
		{
			description: "success",
			requestbody: `{
                "property_id" : "0bd03187-91ac-457d-957c-3ba2f6c0d24b",
                "rating" : 0,
                "review" : "rating is 0 in range 0 to 5"
            }`,
			expectStatus: fiber.StatusOK,
		},
		{
			description: "rating is not integer",
			requestbody: `{
				"property_id" : "0bd03187-91ac-457d-957c-3ba2f6c0d24b",
				"rating" : 1.5,
				"review" : "rating is not integer but float"
			}`,
			expectStatus: fiber.StatusBadRequest,
		},
		{
			description: "rating is not in range 1",
			requestbody: `{
				"property_id" : "0bd03187-91ac-457d-957c-3ba2f6c0d24b",
				"rating" : 6,
				"review" : "rating is not in range 0 to 5 but got 6" 
			}`,
			expectStatus: fiber.StatusBadRequest,
		},
		{
			description: "request body is not valid",
			requestbody: `{
				"property_id" : "0bd03187-91ac-457d-957c-3ba2f6c0d24b",
				"rating" : 5,
				"err" : "my error"
			`,
			expectStatus: fiber.StatusBadRequest,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.description, func(t *testing.T) {
			reqBody := []byte(tc.requestbody)
			service := mockService{
				createrating: func(rating *models.Reviews) error {
					return nil
				},
			}
			handler := NewHandler(&service)
			app.Post("/api/v1/ratings", handler.CreateRating)
			req := httptest.NewRequest("POST", "/api/v1/ratings", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "session", Value: tc.cookieValue})

			resp, _ := app.Test(req)
			if resp.StatusCode != tc.expectStatus {
				t.Fatalf("expected %d but got %d", tc.expectStatus, resp.StatusCode)
			}
			// You can also check the response body, headers, etc. here if needed
		})
	}
}
