package repository

import (
	"context"
	"kelasmirai_backend/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) model.PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Store(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).
		Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	var payment model.Payment

	err := r.db.WithContext(ctx).
		Preload("Invoice").
		Where("id = ?", id).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("id = ?", payment.ID).
		Updates(payment).Error
}

func (r *paymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.Payment{}, "id = ?", id).Error
}

func (r *paymentRepository) GetByMidtransOrderID(ctx context.Context, midtransOrderID string) (*model.Payment, error) {
	var payment model.Payment

	err := r.db.WithContext(ctx).
		Preload("Invoice").
		Where("midtrans_order_id = ?", midtransOrderID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) GetByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment

	err := r.db.WithContext(ctx).
		Where("invoice_id = ?", invoiceID).
		Order("created_at DESC").
		Find(&payments).Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *paymentRepository) MarkSettlement(ctx context.Context, id uuid.UUID, transactionID string) error {
	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":                  model.PaymentStatusSettlement,
			"midtrans_transaction_id": transactionID,
			"paid_at":                 now,
		}).Error
}

func (r *paymentRepository) MarkExpired(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("id = ?", id).
		Update("status", model.PaymentStatusExpired).Error
}

func (r *paymentRepository) MarkCancel(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("id = ?", id).
		Update("status", model.PaymentStatusCancel).Error
}

func (r *paymentRepository) MarkDeny(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("id = ?", id).
		Update("status", model.PaymentStatusDeny).Error
}
