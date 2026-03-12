package repository

import (
	"context"
	"kelasmirai_backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type subscriptionPlanRepository struct {
	db *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) model.SubscriptionPlanRepository {
	return &subscriptionPlanRepository{
		db: db,
	}
}

func (r *subscriptionPlanRepository) Get(ctx context.Context, id uuid.UUID) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan

	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", id, true).
		First(&plan).Error

	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (r *subscriptionPlanRepository) GetAll(ctx context.Context) ([]model.SubscriptionPlan, error) {
	var plans []model.SubscriptionPlan

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("price_monthly ASC").
		Find(&plans).Error

	if err != nil {
		return nil, err
	}

	return plans, nil
}
