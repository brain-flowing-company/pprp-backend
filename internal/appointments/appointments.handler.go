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
	DeleteAppointments(c *fiber.Ctx) error
	UpdateAppointmentStatus(c *fiber.Ctx) error
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{
		service,
	}
}

// @router      /api/v1/appointments [get]
// @summary     Get all appointments
// @description Get all appointments
// @tags        appointments
// @produce     json
// @success     200	{object} []models.Appointments
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) GetAllAppointments(c *fiber.Ctx) error {
	var apps []models.Appointments
	err := h.service.GetAllAppointments(&apps)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(apps)
}

// @router      /api/v1/appointments/:appointmentId [get]
// @summary     Get appointments by id
// @description Get appointments by id
// @tags        appointments
// @produce     json
// @success     200	{object} []models.Appointments
// @failure     400 {object} models.ErrorResponse "Invalid appointment id"
// @failure     404 {object} models.ErrorResponse "Appointment id not found"
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) GetAppointmentById(c *fiber.Ctx) error {
	appointmentId := c.Params("appointmentId")

	var apps models.Appointments
	err := h.service.GetAppointmentsById(&apps, appointmentId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(apps)
}

// @router      /api/v1/appointments [post]
// @summary     Create appointments
// @description Create appointments
// @tags        appointments
// @produce     json
// @param       body body models.CreatingAppointments true "Appointment details"
// @success     201	{object} models.Appointments
// @failure     400 {object} models.ErrorResponse "Empty dates or some of appointments duplicate with existing one"
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) CreateAppointments(c *fiber.Ctx) error {
	apps := &models.CreatingAppointments{}
	err := c.BodyParser(apps)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse body: %v", err.Error())))
	}

	apperr := h.service.CreateAppointments(apps)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "Appointments created")
}

// @router      /api/v1/appointments/ [delete]
// @summary     Delete appointments
// @description Delete **all appointments** in body.
// @tags        appointments
// @produce     json
// @param       body body models.DeletingAppointments true "Appointment id deleting lists"
// @success     200	{object} []models.Appointments
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) DeleteAppointments(c *fiber.Ctx) error {
	appIds := &[]string{}
	err := c.BodyParser(appIds)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse body: %v", err.Error())))
	}

	apperr := h.service.DeleteAppointments(appIds)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "Appointments deleted")
}

// @router      /api/v1/appointments/:appointmentId [patch]
// @summary     Update appointment status
// @description Update appointment status
// @tags        appointments
// @produce     json
// @param       body body models.DeletingAppointments true "Appointment id deleting lists"
// @success     200	{object} []models.Appointments
// @failure     400 {object} models.ErrorResponse
// @failure     500 {object} models.ErrorResponse
func (h *handlerImpl) UpdateAppointmentStatus(c *fiber.Ctx) error {
	status := models.UpdatingAppointmentStatus{}
	err := c.BodyParser(&status)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse body: %v", err.Error())))
	}

	appId := c.Params("appointmentId")

	apperr := h.service.UpdateAppointmentStatus(appId, status.Status)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Appointment state updated")
}
