package repository

import (
	"context"
	"kelasmirai_backend/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

func NewEmailVerificationTokenRepository(db *gorm.DB) model.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{
		db: db,
	}
}

func (r *emailVerificationTokenRepository) WithTx(tx *gorm.DB) model.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: tx}
}

func (r *emailVerificationTokenRepository) Store(ctx context.Context, token *model.EmailVerificationToken) error {
	return r.db.WithContext(ctx).
		Create(token).Error
}

func (r *emailVerificationTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.EmailVerificationToken, error) {
	var token model.EmailVerificationToken

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&token).Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *emailVerificationTokenRepository) Update(ctx context.Context, token *model.EmailVerificationToken) error {
	return r.db.WithContext(ctx).
		Model(&model.EmailVerificationToken{}).
		Where("id = ?", token.ID).
		Updates(token).Error
}

func (r *emailVerificationTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.EmailVerificationToken{}, "id = ?", id).Error
}

func (r *emailVerificationTokenRepository) GetValidByUserID(ctx context.Context, userID uuid.UUID) (*model.EmailVerificationToken, error) {
	var token model.EmailVerificationToken

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("used_at IS NULL").
		Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		First(&token).Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *emailVerificationTokenRepository) GetByToken(ctx context.Context, tokenStr string) (*model.EmailVerificationToken, error) {
	var token model.EmailVerificationToken

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ?", tokenStr).
		First(&token).Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}
