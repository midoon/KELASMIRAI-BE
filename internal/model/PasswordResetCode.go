package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetCode struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	Code      string     `gorm:"type:text;not null"`
	ExpiresAt time.Time  `gorm:"type:timestamptz;not null;index:idx_password_reset_expires_at"`
	UsedAt    *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (PasswordResetCode) TableName() string {
	return "password_reset_codes"
}

func (m *PasswordResetCode) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

type PasswordResetCodeRepository interface {
	Store(ctx context.Context, code *PasswordResetCode) error
	GetByID(ctx context.Context, id uuid.UUID) (*PasswordResetCode, error)
	Update(ctx context.Context, code *PasswordResetCode) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetValidByUserID(ctx context.Context, userID uuid.UUID) (*PasswordResetCode, error)
	GetByCode(ctx context.Context, code string) (*PasswordResetCode, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
