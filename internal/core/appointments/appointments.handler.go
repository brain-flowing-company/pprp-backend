package appointments

import (
	"fmt"
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	GetAllAppointments(c *fiber.Ctx) error
	GetAppointmentById(c *fiber.Ctx) error
	CreateAppointment(c *fiber.Ctx) error
	DeleteAppointment(c *fiber.Ctx) error
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
// @success     200	{object} []models.AppointmentLists
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) GetAllAppointments(c *fiber.Ctx) error {
	var appointments []models.AppointmentLists
	err := h.service.GetAllAppointments(&appointments)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(appointments)
}

// @router      /api/v1/appointments/:appointmentId [get]
// @summary     Get appointments by id
// @description Get appointments and other related information by id
// @tags        appointments
// @produce     json
// @success     200	{object} []models.AppointmentDetails
// @failure     400 {object} models.ErrorResponses "Invalid appointment id"
// @failure     404 {object} models.ErrorResponses "Appointment id not found"
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) GetAppointmentById(c *fiber.Ctx) error {
	appointmentId := c.Params("appointmentId")

	var appointment models.AppointmentDetails
	err := h.service.GetAppointmentById(&appointment, appointmentId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(appointment)
}

// @router      /api/v1/appointments [post]
// @summary     Create appointments
// @description Create appointments with **property_id**, **owner_user_id**, **dweller_user_id**, **appointment_date** **note**(optional)
// @tags        appointments
// @produce     json
// @param       body body models.Appointments true "Appointment details"
// @success     201	{object} models.Appointments
// @failure     400 {object} models.ErrorResponses "Empty dates or some of appointments duplicate with existing one"
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) CreateAppointment(c *fiber.Ctx) error {
	appointment := &models.Appointments{}
	err := c.BodyParser(appointment)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse body: %v", err.Error())))
	}

	apperr := h.service.CreateAppointment(appointment)
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
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) DeleteAppointment(c *fiber.Ctx) error {
	appointmentId := c.Params("appointmentId")

	apperr := h.service.DeleteAppointment(appointmentId)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "Appointments deleted")
}

// @router      /api/v1/appointments/:appointmentId
// @summary     Update appointment status
// @description Update appointment status
// @tags        appointments
// @produce     json
// @param       body body models.UpdatingAppointmentStatus true
// @success     200	{object} []models.Appointments
// @failure     400 {object} models.ErrorResponses
// @failure     500 {object} models.ErrorResponses
func (h *handlerImpl) UpdateAppointmentStatus(c *fiber.Ctx) error {
	updatingAppointment := models.UpdatingAppointmentStatus{}
	err := c.BodyParser(&updatingAppointment)
	if err != nil {
		return utils.ResponseError(c, apperror.
			New(apperror.BadRequest).
			Describe(fmt.Sprintf("Could not parse body: %v", err.Error())))
	}

	appointmentId := c.Params("appointmentId")

	apperr := h.service.UpdateAppointmentStatus(&updatingAppointment, appointmentId)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusOK, "Appointment state updated")
}
