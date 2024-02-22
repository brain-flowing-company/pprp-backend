package chats

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllChats(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router      /api/v1/chats [get]
// @summary     Get current users chat *use cookies*
// @description Get current users chat
// @tags        chats
// @produce     json
// @success     200	{object} []models.ChatsResponses
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) GetAllChats(c *fiber.Ctx) error {
	session, ok := c.Locals("session").(models.Sessions)
	if !ok {
		session = models.Sessions{}
	}

	var chats []models.ChatsResponses
	err := h.service.GetAllChats(&chats, session.UserId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(chats)
}
