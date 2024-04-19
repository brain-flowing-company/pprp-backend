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

// func TestGetAllRatings(t *testing.T) {
// 	testCase := []struct {
// 		name     string
// 		route    string
// 		err      *apperror.AppError
// 		expected int
// 	}{
// 		{
// 			name:     "success",
// 			err:      nil,
// 			expected: 200,
// 			route:    "/api/v1/ratings",
// 		},
// 	}
// 	app := fiber.New()
// 	// Create route with GET method for test
// 	app.Get("/hello", func(c *fiber.Ctx) error {
// 		// Return simple string as response
// 		return c.SendString("Hello, World!")
// 	})
// 	for _, tc := range testCase {
// 		t.Run(tc.name, func(t *testing.T) {
// 			service := mockService{
// 				getallrating: func(rating *[]models.RatingResponse) error {
// 					return tc.err
// 				},
// 			}
// 			handler := NewHandler(&service)
// 			fmt.Println("handler = ", handler)

// 			// Create a new HTTP request
// 			app.Get(tc.route, handler.GetAllRatings)
// 			req := httptest.NewRequest(http.MethodGet, tc.route, nil)
// 			req.Header.Set("Content-Type", "application/json")

// 			resp, err := app.Test(req)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			assert.Equal(t, resp.StatusCode, tc.expected)

// 		})
// 	}

// }
//
//	{
//		"property_id" : "0bd03187-91ac-457d-957c-3ba2f6c0d24b",
//		"rating" : 1,
//		"review" : "kuayfkdfjl"
//	 }

// func TestCreateRating(t *testing.T) {
// 	app := fiber.New()
// 	testCase := []struct {
// 		description  string
// 		requestbody  string
// 		expectStatus int
// 		cookieValue  string
// 	}{
// 		{
// 			description: "success",
// 			requestbody: `{
// 				"property_id" : "0bd03187-91ac-457d-957c-3ba2f6c0d24b",
// 				"rating" : 1,
// 				"review" : "kuayfkdfjl"
// 			}`,
// 			expectStatus: fiber.StatusOK,
// 			cookieValue:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2MjM1MjIsImlhdCI6MTcxMzUzNzEyMiwic2Vzc2lvbiI6eyJ1c2VyX2lkIjoiZGM0YmQwOTAtMWM3OS00N2M1LTlmOWQtNzUyOGZlMGMyOWQxIiwiZW1haWwiOiJuaW5ldGVuNjA5QGdtYWlsLmNvbSJ9fQ.mr4gvZFPICkrnv0r3Tter1r43Uv4S27hue2WZYm8PgA",
// 		},
// 	}
// 	for _, tc := range testCase {
// 		t.Run(tc.description, func(t *testing.T) {
// 			reqBody, _ := json.Marshal(tc.requestbody)
// 			service := mockService{
// 				createrating: func(rating *models.Reviews) error {
// 					return nil
// 				},
// 			}
// 			handler := NewHandler(&service)
// 			app.Post("/api/v1/ratings", handler.CreateRating)
// 			req := httptest.NewRequest("POST", "/api/v1/ratings", bytes.NewReader(reqBody))
// 			req.Header.Set("Content-Type", "application/json")
// 			// Set session cookie
// 			req.AddCookie(&http.Cookie{Name: "session", Value: tc.cookieValue})
// 			resp, _ := app.Test(req)
// 			if resp.StatusCode != tc.expectStatus {
// 				t.Fatalf("expected %d but got %d", tc.expectStatus, resp.StatusCode)
// 			}

// 		})
// 	}

// }

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
                "rating" : 1,
                "review" : "kuayfkdfjlklfjakl"
            }`,
			expectStatus: fiber.StatusOK,
			cookieValue:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2MjUxODAsImlhdCI6MTcxMzUzODc4MCwic2Vzc2lvbiI6eyJ1c2VyX2lkIjoiZmZhNWZmNjItNGVmMi00ODAzLWEwOTktNGYyMzliNTE5MmRkIiwiZW1haWwiOiJuaW5ldGVuNjA5QGdtYWlsLmNvbSJ9fQ.pXZ0wOneV98F1i69VzFu2Ukfe7xfzlsx3g_nCpyoqgw",
		},
		{
			description: "failed",
			requestbody: `{
				"property_id" : "0bd03187-91ac-457d-957c-3ba2f6c0d24b",
				"rating" : 1,
				"review" : "kuayfkdfjl"
			}`,
			expectStatus: fiber.StatusUnauthorized,
			cookieValue:  "",
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

			// Set session cookie
			req.AddCookie(&http.Cookie{Name: "session", Value: tc.cookieValue})
			resp, _ := app.Test(req)
			if resp.StatusCode != tc.expectStatus {
				t.Fatalf("expected %d but got %d", tc.expectStatus, resp.StatusCode)
			}
			// You can also check the response body, headers, etc. here if needed
		})
	}
}
