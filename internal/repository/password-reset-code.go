package repository

import (
	"context"
	"kelasmirai_backend/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type passwordResetCodeRepository struct {
	db *gorm.DB
}

func NewPasswordResetCodeRepository(db *gorm.DB) model.PasswordResetCodeRepository {
	return &passwordResetCodeRepository{
		db: db,
	}
}

func (r *passwordResetCodeRepository) Store(ctx context.Context, code *model.PasswordResetCode) error {
	return r.db.WithContext(ctx).
		Create(code).Error
}

func (r *passwordResetCodeRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PasswordResetCode, error) {
	var resetCode model.PasswordResetCode

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&resetCode).Error

	if err != nil {
		return nil, err
	}

	return &resetCode, nil
}

func (r *passwordResetCodeRepository) Update(ctx context.Context, code *model.PasswordResetCode) error {
	return r.db.WithContext(ctx).
		Model(&model.PasswordResetCode{}).
		Where("id = ?", code.ID).
		Updates(code).Error
}

func (r *passwordResetCodeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.PasswordResetCode{}, "id = ?", id).Error
}

func (r *passwordResetCodeRepository) GetValidByUserID(ctx context.Context, userID uuid.UUID) (*model.PasswordResetCode, error) {
	var resetCode model.PasswordResetCode

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("used_at IS NULL").
		Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		First(&resetCode).Error

	if err != nil {
		return nil, err
	}

	return &resetCode, nil
}

func (r *passwordResetCodeRepository) GetByCode(ctx context.Context, code string) (*model.PasswordResetCode, error) {
	var resetCode model.PasswordResetCode

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("code = ?", code).
		First(&resetCode).Error

	if err != nil {
		return nil, err
	}

	return &resetCode, nil
}

func (r *passwordResetCodeRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.PasswordResetCode{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}

func (r *passwordResetCodeRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.PasswordResetCode{}).Error
}
