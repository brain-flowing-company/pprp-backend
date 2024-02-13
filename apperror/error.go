package apperror

import "net/http"

type AppError struct {
	errorType *AppErrorType
	message   *string
}

func New(errorType *AppErrorType) *AppError {
	return &AppError{
		errorType: errorType,
		message:   nil,
	}
}

func (apperr *AppError) Describe(message string) *AppError {
	apperr.message = new(string)
	*apperr.message = message

	return apperr
}

func (apperr *AppError) Error() string {
	return *apperr.message
}

func (apperr *AppError) Code() int {
	return apperr.errorType.Code
}

func (apperr *AppError) Name() string {
	return apperr.errorType.Name
}

type AppErrorType struct {
	Code int    `json:"code" example:"500"`
	Name string `json:"name" example:"internal-server-error"`
}

var (
	// server errors
	InternalServerError = &AppErrorType{http.StatusInternalServerError, "internal-server-error"}
	InvalidBody         = &AppErrorType{http.StatusBadRequest, "invalid-body"}
	BadRequest          = &AppErrorType{http.StatusBadRequest, "bad-request"}

	//property errors
	InvalidPropertyId = &AppErrorType{http.StatusBadRequest, "invalid-property-id"}
	PropertyNotFound  = &AppErrorType{http.StatusNotFound, "property-not-found"}

	// user errors
	InvalidUserId            = &AppErrorType{http.StatusBadRequest, "invalid-user-id"}
	UserNotFound             = &AppErrorType{http.StatusNotFound, "user-not-found"}
	EmailAlreadyExists       = &AppErrorType{http.StatusBadRequest, "email-already-exists"}
	PhoneNumberAlreadyExists = &AppErrorType{http.StatusBadRequest, "phone-number-already-exists"}
	InvalidEmail             = &AppErrorType{http.StatusBadRequest, "invalid-email"}
	InvalidPassword          = &AppErrorType{http.StatusBadRequest, "invalid-password"}
	InvalidCredentials       = &AppErrorType{http.StatusUnauthorized, "invalid-credentials"}
	ServiceUnavailable       = &AppErrorType{http.StatusServiceUnavailable, "service-unavailable"}

	Unauthorized = &AppErrorType{http.StatusUnauthorized, "unauthorized"}
)
