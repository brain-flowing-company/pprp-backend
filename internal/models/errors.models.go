package models

type ErrorResponses struct {
	Code    int    `json:"code"    example:"500"`
	Name    string `json:"name"    example:"internal-server-error"`
	Message string `json:"message,omitempty" example:"internal server error"`
}

type MessageResponses struct {
	Message string `json:"message" example:"Message sent successfully"`
}
