package models

type EmailType interface {
	Path() string
}

type VerificationEmails struct {
	VerificationLink string
}

func (v VerificationEmails) Path() string {
	return "templates/VerificationEmail.html"
}
