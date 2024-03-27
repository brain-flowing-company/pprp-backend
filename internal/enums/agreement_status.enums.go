package enums

type AgreementStatus string

const (
	AwaitingDepositAgreement AgreementStatus = "AWAITING_DEPOSIT"
	AwaitingPaymentAgreement AgreementStatus = "AWAITING_PAYMENT"
	RentingAgreement         AgreementStatus = "RENTING"
	CancelledAgreement       AgreementStatus = "CANCELLED"
	OverdueAgreement         AgreementStatus = "OVERDUE"
	ArchivedAgreement        AgreementStatus = "ARCHIVED"
)

var AgreementStatusMap = map[string]AgreementStatus{
	"AWAITING_DEPOSIT": AwaitingDepositAgreement,
	"AWAITING_PAYMENT": AwaitingPaymentAgreement,
	"RENTING":          RentingAgreement,
	"CANCELLED":        CancelledAgreement,
	"OVERDUE":          OverdueAgreement,
	"ARCHIVED":         ArchivedAgreement,
}
