package appointments

import (
	"fmt"
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllAppointments(c *fiber.Ctx) error
	GetAppointmentById(c *fiber.Ctx) error
	CreateAppointments(c *fiber.Ctx) error
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

func (h *handlerImpl) CreateAppointments(c *fiber.Ctx) error {
	apps := &models.CreatingAppointments{}
	err := c.BodyParser(apps)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse form data: %v", err.Error())))
	}

	apperr := h.service.CreateAppointments(apps)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "Appointments created")
}
