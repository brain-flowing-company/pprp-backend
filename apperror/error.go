package apperror

import "net/http"

type AppError struct {
	Code int
	Name string
}

var (
	InternalServerError = &AppError{http.StatusInternalServerError, "internal-server-error"}
	InvalidPropertyId   = &AppError{http.StatusBadRequest, "invalid-property-id"}
)
