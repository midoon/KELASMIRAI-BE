package repository

import (
	"context"
	"kelasmirai_backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) model.InvoiceRepository {
	return &invoiceRepository{
		db: db,
	}
}

func (r *invoiceRepository) Store(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).
		Create(invoice).Error
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	var invoice model.Invoice

	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Subscription").
		Where("id = ?", id).
		First(&invoice).Error

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepository) Update(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).
		Model(&model.Invoice{}).
		Where("id = ?", invoice.ID).
		Updates(invoice).Error
}

func (r *invoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.Invoice{}, "id = ?", id).Error
}

func (r *invoiceRepository) GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]model.Invoice, error) {
	var invoices []model.Invoice

	err := r.db.WithContext(ctx).
		Where("subscription_id = ?", subscriptionID).
		Order("created_at DESC").
		Find(&invoices).Error

	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *invoiceRepository) GetPendingByTenantID(ctx context.Context, tenantID uuid.UUID) ([]model.Invoice, error) {
	var invoices []model.Invoice

	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("status = ?", model.InvoiceStatusPending).
		Order("due_date ASC").
		Find(&invoices).Error

	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *invoiceRepository) GetByCode(ctx context.Context, code string) (*model.Invoice, error) {
	var invoice model.Invoice

	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Subscription").
		Where("code = ?", code).
		First(&invoice).Error

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepository) MarkPaid(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Invoice{}).
		Where("id = ?", id).
		Update("status", model.InvoiceStatusPaid).Error
}

func (r *invoiceRepository) MarkExpired(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Invoice{}).
		Where("id = ?", id).
		Update("status", model.InvoiceStatusExpired).Error
}
