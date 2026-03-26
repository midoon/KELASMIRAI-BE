package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleTeacher UserRole = "teacher"
	UserRoleStaff   UserRole = "staff"
	UserRoleStudent UserRole = "student"
	UserRoleParent  UserRole = "parent"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:unique_user_email_per_tenant"`
	Name         string    `gorm:"type:text;not null"`
	Email        string    `gorm:"type:text;not null;uniqueIndex:unique_user_email_per_tenant"`
	PasswordHash string    `gorm:"type:text;not null"`
	Role         UserRole  `gorm:"type:user_role;not null"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()"`
	Tenant       Tenant    `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "users"
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	Store(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
