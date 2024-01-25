package greeting

import "github.com/brain-flowing-company/pprp-backend/internal/dto"

type Service interface {
	Greeting() dto.GreetingResponse
}

type serviceImpl struct{}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) Greeting() dto.GreetingResponse {
	return dto.GreetingResponse{Message: "Hello, World"}
}
