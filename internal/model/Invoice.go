package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusPending  InvoiceStatus = "pending"
	InvoiceStatusPaid     InvoiceStatus = "paid"
	InvoiceStatusCanceled InvoiceStatus = "canceled"
	InvoiceStatusExpired  InvoiceStatus = "expired"
)

type Invoice struct {
	ID             uuid.UUID          `gorm:"type:uuid;primaryKey"`
	TenantID       uuid.UUID          `gorm:"type:uuid;not null;index:idx_invoices_tenant_id"`
	SubscriptionID uuid.UUID          `gorm:"type:uuid;not null"`
	Code           string             `gorm:"type:text;not null;uniqueIndex:unique_invoice_code"`
	Amount         decimal.Decimal    `gorm:"type:numeric(15,2);not null;check:amount >= 0"`
	Status         InvoiceStatus      `gorm:"type:invoice_status;not null;default:'pending';index:idx_invoices_status"`
	DueDate        time.Time          `gorm:"type:timestamptz;not null;index:idx_invoices_due_date"`
	CreatedAt      time.Time          `gorm:"type:timestamptz;not null;default:now()"`
	Tenant         Tenant             `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Subscription   TenantSubscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE"`
}

func (Invoice) TableName() string {
	return "invoices"
}

func (m *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

type InvoiceRepository interface {
	Store(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	Update(ctx context.Context, invoice *Invoice) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]Invoice, error)
	GetPendingByTenantID(ctx context.Context, tenantID uuid.UUID) ([]Invoice, error)
	GetByCode(ctx context.Context, code string) (*Invoice, error)
	MarkPaid(ctx context.Context, id uuid.UUID) error
	MarkExpired(ctx context.Context, id uuid.UUID) error
}
