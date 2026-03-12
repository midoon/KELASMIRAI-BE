package repository

import (
	"context"
	"kelasmirai_backend/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tenantSubscriptionRepository struct {
	db *gorm.DB
}

func NewTenantSubscriptionRepository(db *gorm.DB) model.TenantSubscriptionRepository {
	return &tenantSubscriptionRepository{
		db: db,
	}
}

func (r *tenantSubscriptionRepository) Store(ctx context.Context, subscription *model.TenantSubscription) error {
	return r.db.WithContext(ctx).
		Create(subscription).Error
}

func (r *tenantSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TenantSubscription, error) {
	var subscription model.TenantSubscription

	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Plan").
		Where("id = ?", id).
		First(&subscription).Error

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *tenantSubscriptionRepository) Update(ctx context.Context, subscription *model.TenantSubscription) error {
	return r.db.WithContext(ctx).
		Model(&model.TenantSubscription{}).
		Where("id = ?", subscription.ID).
		Updates(subscription).Error
}

func (r *tenantSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.TenantSubscription{}, "id = ?", id).Error
}

func (r *tenantSubscriptionRepository) GetActiveByTenantID(ctx context.Context, tenantID uuid.UUID) (*model.TenantSubscription, error) {
	var subscription model.TenantSubscription

	err := r.db.WithContext(ctx).
		Preload("Plan").
		Where("tenant_id = ?", tenantID).
		Where("status IN ?", []model.SubscriptionStatus{
			model.SubscriptionStatusTrial,
			model.SubscriptionStatusActive,
		}).
		Order("created_at DESC").
		First(&subscription).Error

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *tenantSubscriptionRepository) GetByMidtransSubscriptionID(
	ctx context.Context,
	id string,
) (*model.TenantSubscription, error) {

	var sub model.TenantSubscription

	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Plan").
		Where("midtrans_subscription_id = ?", id).
		First(&sub).Error

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *tenantSubscriptionRepository) Activate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.TenantSubscription{}).
		Where("id = ?", id).
		Update("status", model.SubscriptionStatusActive).Error
}

func (r *tenantSubscriptionRepository) MarkPastDue(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.TenantSubscription{}).
		Where("id = ?", id).
		Update("status", model.SubscriptionStatusPastDue).Error
}

func (r *tenantSubscriptionRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.TenantSubscription{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":   model.SubscriptionStatusCanceled,
			"ended_at": now,
		}).Error
}
