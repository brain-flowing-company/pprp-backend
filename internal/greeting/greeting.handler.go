package greeting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Greeting(c *gin.Context)
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router      /greeting [get]
// @summary     Greeting
// @description hello, world endpoint
// @produce     json
// @success     200	{object}	dto.GreetingResponse
func (h *handlerImpl) Greeting(c *gin.Context) {
	res := h.service.Greeting()

	c.JSON(http.StatusOK, res)
}
