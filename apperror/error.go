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
	// server errors
	InternalServerError = &AppError{http.StatusInternalServerError, "internal-server-error"}
	BadRequest          = &AppError{http.StatusBadRequest, "bad-request"}

	//property errors
	InvalidPropertyId = &AppError{http.StatusBadRequest, "invalid-property-id"}
	PropertyNotFound  = &AppError{http.StatusNotFound, "property-not-found"}

	// user errors
	InvalidUserId          = &AppError{http.StatusBadRequest, "invalid-user-id"}
	UserNotFound           = &AppError{http.StatusNotFound, "user-not-found"}
	EmailAlreadyExists     = &AppError{http.StatusBadRequest, "email-already-exists"}
	PasswordCannotBeHashed = &AppError{http.StatusInternalServerError, "password-cannot-be-hashed"}
	InvalidEmail           = &AppError{http.StatusBadRequest, "invalid-email"}
	InvalidPassword        = &AppError{http.StatusUnauthorized, "invalid-password"}
	InvalidCredentials     = &AppError{http.StatusUnauthorized, "invalid-credentials"}
)
