package model

import (
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
