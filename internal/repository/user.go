package repository

import (
	"context"
	"errors"
	"fmt"
	"kelasmirai_backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) WithTx(tx *gorm.DB) model.UserRepository {
	return &userRepository{db: tx}
}

func (r *userRepository) Store(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("userRepository.Store: %w", result.Error)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("userRepository.GetByID: user with id %s not found", id)
		}
		return nil, fmt.Errorf("userRepository.GetByID: %w", result.Error)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("userRepository.GetByEmail: user with email %s not found in tenant %s", email, tenantID)
		}
		return nil, fmt.Errorf("userRepository.GetByEmail: %w", result.Error)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("userRepository.Update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("userRepository.Update: user with id %s not found", user.ID)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.User{})
	if result.Error != nil {
		return fmt.Errorf("userRepository.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("userRepository.Delete: user with id %s not found", id)
	}
	return nil
}
