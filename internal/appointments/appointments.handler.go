package appointments

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllAppointments(c *fiber.Ctx) error
	GetAppointmentById(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

func (h *handlerImpl) GetAllAppointments(c *fiber.Ctx) error {
	var apps []models.Appointments
	err := h.service.GetAllAppointments(&apps)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(apps)
}

func (h *handlerImpl) GetAppointmentById(c *fiber.Ctx) error {
	appointmentId := c.Params("appointmentId")

	var apps models.Appointments
	err := h.service.GetAppointmentsById(&apps, appointmentId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(apps)
}
