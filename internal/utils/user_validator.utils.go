package utils

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/go-playground/validator/v10"
)

func RegisterTypeValidator(registerType validator.FieldLevel) bool {
	switch registerType.Field().String() {
	case "EMAIL", "GOOGLE":
		return true
	}
	return false
}

func PhoneValidator(phone validator.FieldLevel) bool {
	return len(phone.Field().String()) == 10 && phone.Field().String()[0] == '0'
}

func CardNumberValidator(cardNumber validator.FieldLevel) bool {
	if len(cardNumber.Field().String()) != 16 {
		return false
	} else {
		for _, c := range cardNumber.Field().String() {
			if c < '0' || c > '9' {
				return false
			}
		}
	}

	return true
}

func ExpireMonthValidator(expireMonth validator.FieldLevel) bool {
	if len(expireMonth.Field().String()) != 2 {
		return false
	} else if expireMonth.Field().String() < "01" || expireMonth.Field().String() > "12" {
		return false
	} else {
		c0 := expireMonth.Field().String()[0]
		c1 := expireMonth.Field().String()[1]
		if c0 < '0' || c0 > '9' || c1 < '0' || c1 > '9' {
			return false
		}
	}

	return true
}

func ExpireYearValidator(expireYear validator.FieldLevel) bool {
	if len(expireYear.Field().String()) != 4 {
		return false
	} else if expireYear.Field().String() < "2024" {
		return false
	} else {
		for _, c := range expireYear.Field().String() {
			if c < '0' || c > '9' {
				return false
			}
		}
	}

	return true
}

func CVVValidator(cvv validator.FieldLevel) bool {
	if len(cvv.Field().String()) != 3 {
		return false
	} else {
		for _, c := range cvv.Field().String() {
			if c < '0' || c > '9' {
				return false
			}
		}
	}

	return true
}

func CardColorValidator(cardColor validator.FieldLevel) bool {
	switch cardColor.Field().String() {
	case "LIGHT_BLUE", "BLUE", "DARK_BLUE", "VERY_DARK_BLUE":
		return true
	}

	return false
}

func CreditCardsValidator(creditCards validator.FieldLevel) bool {
	v, validatorErr := NewCreditCardValidator()
	if validatorErr != nil {
		return false
	}

	ccs := creditCards.Field().Interface().([]models.CreditCards)
	if len(ccs) > 4 {
		fmt.Println("CreditCards length is not in range")
		return false
	} else if len(ccs) > 0 {
		var tagNumbers = make(map[int]bool)

		for i := 0; i < len(ccs); i++ {
			tagNumbers[i+1] = false
		}

		for _, cc := range ccs {
			if err := v.Struct(cc); err != nil {
				fmt.Println("Credit Card Error:", err)
				return false
			}

			if check := tagNumbers[cc.TagNumber]; !check {
				tagNumbers[cc.TagNumber] = true
			} else {
				fmt.Println("Tag Number is duplicated")
				return false
			}
		}

		for _, check := range tagNumbers {
			if !check {
				fmt.Println("Tag Number is missing")
				return false
			}
		}
	}

	return true
}

func BankNameValidator(bankName validator.FieldLevel) bool {
	switch bankName.Field().String() {
	case "KBANK", "BBL", "KTB", "BAY", "CIMB", "TTB", "SCB", "GSB", "":
		return true
	}
	return false
}

func BankAccountNumberValidator(bankAccountNumber validator.FieldLevel) bool {
	return len(bankAccountNumber.Field().String()) == 10
}

func NewUserValidator() (*validator.Validate, error) {
	v := validator.New()
	if err := v.RegisterValidation("register_type", RegisterTypeValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("phone", PhoneValidator); err != nil {
		return nil, err
	}

	return v, nil
}

func NewCreditCardValidator() (*validator.Validate, error) {
	v := validator.New()
	if err := v.RegisterValidation("card_number", CardNumberValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("expire_month", ExpireMonthValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("expire_year", ExpireYearValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("cvv", CVVValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("card_color", CardColorValidator); err != nil {
		return nil, err
	}

	return v, nil
}

func NewUserFinancialInformationValidator() (*validator.Validate, error) {
	v := validator.New()
	if err := v.RegisterValidation("credit_cards", CreditCardsValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("bank_name", BankNameValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("bank_account_number", BankAccountNumberValidator); err != nil {
		return nil, err
	}

	return v, nil
}
