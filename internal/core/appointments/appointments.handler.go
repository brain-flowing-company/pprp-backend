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
	GetMyAppointments(c *fiber.Ctx) error
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
// @summary     Get all appointments *use cookies*
// @description Get all appointments
// @tags        appointments
// @produce     json
// @success     200	{object} []models.AppointmentLists
// @failure     500 {object} models.ErrorResponses "Could not get all appointments"
func (h *handlerImpl) GetAllAppointments(c *fiber.Ctx) error {
	var appointments []models.AppointmentLists
	err := h.service.GetAllAppointments(&appointments)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(appointments)
}

// @router      /api/v1/appointments/:appointmentId [get]
// @summary     Get an appointment by id *use cookies*
// @description Get the appointment and other related information by id
// @tags        appointments
// @produce     json
// @param       appointmentId path string true "Appointment ID"
// @success     200	{object} []models.AppointmentDetails
// @failure     400 {object} models.ErrorResponses "Invalid appointment id"
// @failure     404 {object} models.ErrorResponses "Could not find the specified appointment"
// @failure     500 {object} models.ErrorResponses "Could not get appointment by id"
func (h *handlerImpl) GetAppointmentById(c *fiber.Ctx) error {
	appointmentId := c.Params("appointmentId")

	var appointment models.AppointmentDetails
	err := h.service.GetAppointmentById(&appointment, appointmentId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(appointment)
}

// @router      /api/v1/user/me/appointments [get]
// @summary     Get my appointments *use cookies*
// @description Get all appointments related to the user
// @tags        appointments
// @produce     json
// @success     200	{object} models.MyAppointmentResponses
// @failure     500 {object} models.ErrorResponses "Could not get my appointments"
func (h *handlerImpl) GetMyAppointments(c *fiber.Ctx) error {
	userId := c.Locals("session").(models.Sessions).UserId.String()

	var appointments models.MyAppointmentResponses
	err := h.service.GetMyAppointments(&appointments, userId)
	if err != nil {
		return utils.ResponseError(c, err)
	}

	return c.JSON(appointments)
}

// @router      /api/v1/appointments [post]
// @summary     Create an appointment *use cookies*
// @description Create an appointment by parsing the body (note is optional)
// @tags        appointments
// @produce     json
// @param       body body models.CreatingAppointments true "Appointment details"
// @success     201	{object} models.MessageResponses "Appointments created"
// @failure     400 {object} models.ErrorResponses "Empty dates or some of appointments duplicate with existing one"
// @failure     500 {object} models.ErrorResponses "Could not create appointments"
func (h *handlerImpl) CreateAppointment(c *fiber.Ctx) error {
	appointment := &models.CreatingAppointments{
		DwellerUserId: c.Locals("session").(models.Sessions).UserId,
	}

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

// @router      /api/v1/appointments/:appointmentId [delete]
// @summary     Delete an appointment by id *use cookies*
// @description Delete an appointment by id
// @tags        appointments
// @produce     json
// @param       appointmentId path string true "Appointment ID"
// @success     200	{object} models.MessageResponses "Appointments deleted"
// @failure     500 {object} models.ErrorResponses "Could not delete appointments"
func (h *handlerImpl) DeleteAppointment(c *fiber.Ctx) error {
	appointmentId := c.Params("appointmentId")

	apperr := h.service.DeleteAppointment(appointmentId)
	if apperr != nil {
		return utils.ResponseError(c, apperr)
	}

	return utils.ResponseMessage(c, http.StatusCreated, "Appointments deleted")
}

// @router      /api/v1/appointments/:appointmentId [patch]
// @summary     Update an appointment status by id *use cookies*
// @description Update an appointment status by id with **status** and **cancelled_message**(optional)
// @tags        appointments
// @produce     json
// @param       appointmentId path string true "Appointment ID"
// @param       body body models.UpdatingAppointmentStatus true "Appointment status and cancelled message(optional)"
// @success     200	{object} models.MessageResponses "Appointment state updated"
// @failure     400 {object} models.ErrorResponses "Invalid appointment id"
// @failure     500 {object} models.ErrorResponses "Could not update appointment status"
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
