package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSettlement PaymentStatus = "settlement"
	PaymentStatusExpired    PaymentStatus = "expired"
	PaymentStatusCancel     PaymentStatus = "cancel"
	PaymentStatusDeny       PaymentStatus = "deny"
)

type Payment struct {
	ID                    uuid.UUID       `gorm:"type:uuid;primaryKey"`
	InvoiceID             uuid.UUID       `gorm:"type:uuid;not null;index:idx_payments_invoice_id"`
	MidtransOrderID       string          `gorm:"type:text;not null;uniqueIndex:unique_midtrans_order_id"`
	MidtransTransactionID *string         `gorm:"type:text;index:idx_payments_midtrans_transaction_id"`
	Amount                decimal.Decimal `gorm:"type:numeric(15,2);not null;check:amount >= 0"`
	Status                PaymentStatus   `gorm:"type:payment_status;not null;default:'pending';index:idx_payments_status"`
	PaidAt                *time.Time      `gorm:"type:timestamptz"`
	RawResponseJSON       datatypes.JSON  `gorm:"type:jsonb"`
	CreatedAt             time.Time       `gorm:"type:timestamptz;not null;default:now()"`
	Invoice               Invoice         `gorm:"foreignKey:InvoiceID;constraint:OnDelete:RESTRICT"`
}

func (Payment) TableName() string {
	return "payments"
}

func (m *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
