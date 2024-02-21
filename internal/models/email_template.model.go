package models

type EmailType interface {
	Path() string
}

type VerificationEmail struct {
	VerificationLink string
}

func (v VerificationEmail) Path() string {
	return "templates/VerificationEmail.html"
}
