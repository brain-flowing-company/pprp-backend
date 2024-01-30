package greeting

import "github.com/brain-flowing-company/pprp-backend/internal/models"

type Service interface {
	Greeting(msg *models.Greeting)
}

type serviceImpl struct{}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) Greeting(msg *models.Greeting) {
	msg.Message = "Hello, World"
}
