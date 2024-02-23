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
	if apperr.message == nil {
		return ""
	} else {
		return *apperr.message
	}
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

	// property errors
	InvalidPropertyId = &AppErrorType{http.StatusBadRequest, "invalid-property-id"}
	PropertyNotFound  = &AppErrorType{http.StatusNotFound, "property-not-found"}

	// appointment errors
	InvalidAppointmentId     = &AppErrorType{http.StatusBadRequest, "invalid-appointment-id"}
	AppointmentNotFound      = &AppErrorType{http.StatusNotFound, "appointment-not-found"}
	DuplicateAppointment     = &AppErrorType{http.StatusBadRequest, "duplicate-appointment"}
	InvalidAppointmentStatus = &AppErrorType{http.StatusBadRequest, "invalid-appointment-status"}

	// user errors
	InvalidUserId                 = &AppErrorType{http.StatusBadRequest, "invalid-user-id"}
	UserNotFound                  = &AppErrorType{http.StatusNotFound, "user-not-found"}
	EmailAlreadyExists            = &AppErrorType{http.StatusBadRequest, "email-already-exists"}
	PhoneNumberAlreadyExists      = &AppErrorType{http.StatusBadRequest, "phone-number-already-exists"}
	InvalidEmail                  = &AppErrorType{http.StatusBadRequest, "invalid-email"}
	InvalidEmailVerificationCode  = &AppErrorType{http.StatusBadRequest, "invalid-email-verification-code"}
	EmailVerificationCodeExpired  = &AppErrorType{http.StatusBadRequest, "email-verification-code-expired"}
	EmailVerificationDataNotFound = &AppErrorType{http.StatusNotFound, "email-verification-data-not-found"}
	InvalidPassword               = &AppErrorType{http.StatusBadRequest, "invalid-password"}
	InvalidCredentials            = &AppErrorType{http.StatusUnauthorized, "invalid-credentials"}
	ServiceUnavailable            = &AppErrorType{http.StatusServiceUnavailable, "service-unavailable"}
	InvalidProfileImageExtension  = &AppErrorType{http.StatusBadRequest, "invalid-profile-image-extensions"}

	InvalidAgreementId = &AppErrorType{http.StatusBadRequest, "invalid-agreement-id"}
	AgreementNotFound  = &AppErrorType{http.StatusNotFound, "agreement-not-found"}
	DuplicateAgreement = &AppErrorType{http.StatusBadRequest, "duplicate-agreement"}

	InvalidCallbackRequest = &AppErrorType{http.StatusBadRequest, "invalid-callback-request"}

	Unauthorized = &AppErrorType{http.StatusUnauthorized, "unauthorized"}
)
