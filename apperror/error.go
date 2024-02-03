package apperror

import "net/http"

type AppError struct {
	Code int    `json:"code" example:"500"`
	Name string `json:"name" example:"internal-server-error"`
}

var (
	InternalServerError = &AppError{http.StatusInternalServerError, "internal-server-error"}
	InvalidPropertyId   = &AppError{http.StatusBadRequest, "invalid-property-id"}
	PropertyNotFound    = &AppError{http.StatusNotFound, "property-not-found"}
	InvalidUserId       = &AppError{http.StatusBadRequest, "invalid-user-id"}
	UserNotFound        = &AppError{http.StatusNotFound, "user-not-found"}
)
