package payments

import (
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
		paymentQuery := `INSERT INTO payments (payment_id , user_id , price ,IsSuccess ,name) VALUES (?,?,?,?,?)`
		if err := tx.Exec(paymentQuery, payment.PaymentId, payment.UserId, payment.Price, payment.IsSuccess, payment.Name).Error; err != nil {
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
