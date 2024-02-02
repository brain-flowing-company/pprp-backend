package apperror

import "net/http"

type AppError struct {
	Code int    `json:"code" example:"500"`
	Name string `json:"name" example:"internal-server-error"`
}

// Error implements error.
func (*AppError) Error() string {
	panic("unimplemented")
}

var (
	InternalServerError = &AppError{http.StatusInternalServerError, "internal-server-error"}
	InvalidPropertyId   = &AppError{http.StatusBadRequest, "invalid-property-id"}
	PropertyNotFound    = &AppError{http.StatusNotFound, "property-not-found"}
	UserNotFound        = &AppError{http.StatusNotFound, "user-not-found"}
	InvalidCredentials  = &AppError{http.StatusUnauthorized, "invalid-credentials"}
	BadRequest          = &AppError{http.StatusBadRequest, "bad-request"}
	InvalidEmail        = &AppError{http.StatusBadRequest, "invalid-email"}
	EmailAlreadyExists  = &AppError{http.StatusBadRequest, "email-already-exists"}
)
