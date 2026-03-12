package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	SubscriptionStatusTrial    SubscriptionStatus = "trial"
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
)

type TenantSubscription struct {
	ID                     uuid.UUID          `gorm:"type:uuid;primaryKey"`
	TenantID               uuid.UUID          `gorm:"type:uuid;not null;index:idx_tenant_subscriptions_tenant_id"`
	PlanID                 uuid.UUID          `gorm:"type:uuid;not null"`
	BillingCycle           string             `gorm:"type:text;not null;check:billing_cycle IN ('monthly','yearly')"`
	Price                  decimal.Decimal    `gorm:"type:numeric(15,2);not null;check:price >= 0"`
	Status                 SubscriptionStatus `gorm:"type:subscription_status;not null;default:'trial';index:idx_tenant_subscriptions_status"`
	StartedAt              time.Time          `gorm:"type:timestamptz;not null"`
	EndedAt                *time.Time         `gorm:"type:timestamptz"`
	NextBillingAt          time.Time          `gorm:"type:timestamptz;not null;index:idx_tenant_subscriptions_next_billing"`
	MidtransSubscriptionID *string            `gorm:"type:text"`
	CreatedAt              time.Time          `gorm:"type:timestamptz;not null;default:now()"`
	Tenant                 Tenant             `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Plan                   SubscriptionPlan   `gorm:"foreignKey:PlanID"`
}

func (TenantSubscription) TableName() string {
	return "tenant_subscriptions"
}

func (m *TenantSubscription) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

type TenantSubscriptionRepository interface {
	Store(ctx context.Context, subscription *TenantSubscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*TenantSubscription, error)
	Update(ctx context.Context, subscription *TenantSubscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveByTenantID(ctx context.Context, tenantID uuid.UUID) (*TenantSubscription, error)
	GetByMidtransSubscriptionID(ctx context.Context, id string) (*TenantSubscription, error)
	Activate(ctx context.Context, id uuid.UUID) error
	MarkPastDue(ctx context.Context, id uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID) error
}
