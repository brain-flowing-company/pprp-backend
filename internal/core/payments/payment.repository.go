package payments

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreatePayment(*models.Payments) error
	GetPaymentByUserId(*models.MyPaymentsResponse, uuid.UUID) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (r *repositoryImpl) CreatePayment(payment *models.Payments) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		// Check if user_id exists in users table
		var userCount int64
		if err := tx.Model(&models.Users{}).Where("user_id = ?", payment.UserId).Count(&userCount).Error; err != nil {
			return err
		}
		if userCount == 0 {
			return apperror.
				New(apperror.UserNotFound).
				Describe("FK constraint in user table")
		}

		// Check if agreement_id exists in agreements table
		var agreementCount int64
		if err := tx.Model(&models.Agreements{}).Where("agreement_id = ?", payment.AgreementId).Count(&agreementCount).Error; err != nil {
			return err
		}
		if agreementCount == 0 {
			return apperror.
				New(apperror.AgreementNotFound).
				Describe("FK constraint in agreement table")
		}

		paymentQuery := `INSERT INTO payments (payment_id , user_id , price ,IsSuccess ,name,agreement_id,payment_method) VALUES (?,?,?,?,?,?,?)`
		if err := tx.Exec(paymentQuery, payment.PaymentId, payment.UserId, payment.Price, payment.IsSuccess, payment.Name, payment.AgreementId, payment.PaymentMethod).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *repositoryImpl) GetPaymentByUserId(payments *models.MyPaymentsResponse, userId uuid.UUID) error {
	paymentQuery := `SELECT * FROM payments WHERE user_id = ?`
	if err := r.db.Raw(paymentQuery, userId).Scan(&payments.Payments).Error; err != nil {
		return err
	}
	return nil
}
