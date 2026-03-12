package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey"`
	Name         string          `gorm:"type:text;not null;unique"`
	PriceMonthly decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0"`
	PriceYearly  decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0"`
	DurationDays int             `gorm:"type:integer;not null;default:30"`
	MaxStudents  *int            `gorm:"type:integer"`
	MaxTeachers  *int            `gorm:"type:integer"`
	FeaturesJSON datatypes.JSON  `gorm:"type:jsonb"`
	IsActive     bool            `gorm:"not null;default:true"`
	CreatedAt    time.Time       `gorm:"type:timestamptz;not null;default:now()"`
}

func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

func (m *SubscriptionPlan) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

type SubscriptionPlanRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*SubscriptionPlan, error)
	GetAll(ctx context.Context) ([]SubscriptionPlan, error)
}
