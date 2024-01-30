package models

type Users struct {
	UserId                    string       `json:"user_id" example:"f38f80b3-f326-4825-9afc-ebc331626875"`
	FirstName                 string       `json:"first_name" example:"John"`
	LastName                  string       `json:"last_name" example:"Doe"`
	Email                     string       `json:"email" example:"email@email.com"`
	PhoneNumber               string       `json:"phone_number" example:"0812345678"`
	ProfileImage              ProfileImage `gorm:"references:UserId" json:"profile_image"`
	CreditCardholderName      string       `json:"credit_cardholder_name" example:"JOHN DOE"`
	CreditCardNumber          string       `json:"credit_card_number" example:"1234567890123456"`
	CreditCardExpirationMonth string       `json:"credit_card_expiration_month" example:"12"`
	CreditCardExpirationYear  string       `json:"credit_card_expiration_year" example:"2023"`
	CreditCardCVV             string       `json:"credit_card_cvv" example:"123"`
	BankName                  string       `json:"bank_name" example:"KBANK"`
	BankAccountNumber         string       `json:"bank_account_number" example:"1234567890"`
	IsVerified                bool         `json:"is_verified" example:"false"`
}

// type BankName string

// const (
// 	KBANK BankName = "KASIKORN BANK"
// 	BBL   BankName = "BANGKOK BANK"
// 	KTB   BankName = "KRUNG THAI BANK"
// 	BAY   BankName = "BANK OF AYUDHYA"
// 	CIMB  BankName = "CIMB THAI BANK"
// 	TTB   BankName = "TMBTHANACHART BANK"
// 	SCB   BankName = "SIAM COMMERCIAL BANK"
// 	GSB   BankName = "GOVERNMENT SAVINGS BANK"
// )

type ProfileImage struct {
	UserId   string `json:"-"`
	ImageUrl string `json:"url" example:"https://image_url.com/abcd"`
}
