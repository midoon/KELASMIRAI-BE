package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusInactive  TenantStatus = "cancelled"
)

type Tenant struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Name        string       `gorm:"type:text;not null"`
	Slug        string       `gorm:"type:text;not null;uniqueIndex"`
	Email       string       `gorm:"type:text;not null;uniqueIndex"`
	Phone       *string      `gorm:"type:text"`
	Address     *string      `gorm:"type:text"`
	LogoURL     *string      `gorm:"type:text"`
	Status      TenantStatus `gorm:"type:tenant_status;not null;default:'active'"`
	TrialEndsAt time.Time    `gorm:"type:timestamptz;not null"`
	CreatedAt   time.Time    `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time    `gorm:"type:timestamptz;not null;default:now()"`
}

func (Tenant) TableName() string {
	return "tenants"
}

func (m *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

type TenantRepository interface {
	WithTx(tx *gorm.DB) TenantRepository
	Store(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}
