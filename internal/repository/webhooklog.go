package repository

import (
	"context"
	"kelasmirai_backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type webhookLogRepository struct {
	db *gorm.DB
}

func NewWebhookLogRepository(db *gorm.DB) model.WebhookLogRepository {
	return &webhookLogRepository{
		db: db,
	}
}

func (r *webhookLogRepository) Store(ctx context.Context, log *model.WebhookLog) error {
	return r.db.WithContext(ctx).
		Create(log).Error
}

func (r *webhookLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.WebhookLog, error) {
	var log model.WebhookLog

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&log).Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *webhookLogRepository) Update(ctx context.Context, log *model.WebhookLog) error {
	return r.db.WithContext(ctx).
		Model(&model.WebhookLog{}).
		Where("id = ?", log.ID).
		Updates(log).Error
}

func (r *webhookLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.WebhookLog{}, "id = ?", id).Error
}

func (r *webhookLogRepository) GetUnprocessedByProvider(ctx context.Context, provider string) ([]model.WebhookLog, error) {
	var logs []model.WebhookLog

	err := r.db.WithContext(ctx).
		Where("provider = ?", provider).
		Where("processed = ?", false).
		Order("created_at ASC").
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *webhookLogRepository) GetByProviderAndExternalID(
	ctx context.Context,
	provider, externalID string,
) (*model.WebhookLog, error) {

	var log model.WebhookLog

	err := r.db.WithContext(ctx).
		Where("provider = ?", provider).
		Where("external_id = ?", externalID).
		First(&log).Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *webhookLogRepository) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.WebhookLog{}).
		Where("id = ?", id).
		Update("processed", true).Error
}
