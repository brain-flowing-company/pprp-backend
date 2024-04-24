package ratings

import (
	"errors"
	"testing"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type mockRepository struct {
	mockCreateRating func(*models.Reviews) error
}

func (r *mockRepository) CreateRating(reviews *models.Reviews) error {
	return r.mockCreateRating(reviews)
}

func TestCreateRatingService(t *testing.T) {
	johnDoeId, _ := uuid.Parse("f38f80b3-f326-4825-9afc-ebc331626555")
	johnDoePropertyId, _ := uuid.Parse("0bd03187-91ac-457d-957c-3ba2f6c0d24b")
	falseId, _ := uuid.Parse("f38f80b3-f326-4825-9afc-ebc331629999")

	internalServerError := errors.New("Internal Server Error")
	testCase := []struct {
		description  string
		review       *models.Reviews
		expectResult error
	}{
		{
			description: "success",
			review: &models.Reviews{
				ReviewId:      uuid.New(),
				DwellerUserId: johnDoeId,
				PropertyId:    johnDoePropertyId,
				Rating:        3,
				Review:        "review review review",
			},
			expectResult: nil,
		},
		{
			description: "Internal Server Error",
			review: &models.Reviews{
				ReviewId:      uuid.New(),
				DwellerUserId: johnDoeId,
				PropertyId:    falseId,
				Rating:        3,
				Review:        "Internal Server Error",
			},
			// expectResult: error(apperror.New(apperror.InternalServerError).Describe("Failed to create rating")),
			// expectResult: errors.New("Internal Server Error"),
			expectResult: internalServerError,
		},
	}

	mockCreateRatingReturn := func(reviews *models.Reviews) error {
		if reviews.Review == "Internal Server Error" {
			// return errors.New("Internal Server Error")
			return internalServerError
		}
		return nil
	}

	for _, tc := range testCase {
		t.Run(tc.description, func(t *testing.T) {
			mockLogger := zap.NewNop()

			mockRepo := &mockRepository{
				mockCreateRating: mockCreateRatingReturn,
			}

			service := NewService(mockRepo, mockLogger, nil)
			err := service.CreateRating(tc.review)
			if err != tc.expectResult {
				t.Errorf("expected %v but got %v", tc.expectResult, err)
			}
		})
	}
}
