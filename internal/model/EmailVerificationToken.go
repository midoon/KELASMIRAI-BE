package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailVerificationToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_email_verification_user_id"`
	Token     string     `gorm:"type:text;not null;uniqueIndex:unique_email_verification_token"`
	ExpiresAt time.Time  `gorm:"type:timestamptz;not null;index:idx_email_verification_expires_at"`
	UsedAt    *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

func (m *EmailVerificationToken) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

type EmailVerificationTokenRepository interface {
	Store(ctx context.Context, token *EmailVerificationToken) error
	GetByID(ctx context.Context, id uuid.UUID) (*EmailVerificationToken, error)
	Update(ctx context.Context, token *EmailVerificationToken) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetValidByUserID(ctx context.Context, userID uuid.UUID) (*EmailVerificationToken, error)
	GetByToken(ctx context.Context, token string) (*EmailVerificationToken, error)
}
