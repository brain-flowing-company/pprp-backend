package greeting

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
)

type Service interface {
	Greeting(*models.Greeting)
	UserGreeting(*models.Greeting, string)
}

type serviceImpl struct{}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) Greeting(msg *models.Greeting) {
	msg.Message = "Hello, World!"
}

func (s *serviceImpl) UserGreeting(msg *models.Greeting, email string) {
	msg.Message = fmt.Sprintf("Hello, %v!", email)
}
