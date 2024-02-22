package greetings

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
)

type Service interface {
	Greeting(*models.Greetings)
	UserGreeting(*models.Greetings, string)
}

type serviceImpl struct{}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) Greeting(msg *models.Greetings) {
	msg.Message = "Hello, World!"
}

func (s *serviceImpl) UserGreeting(msg *models.Greetings, email string) {
	msg.Message = fmt.Sprintf("Hello, %v!", email)
}
