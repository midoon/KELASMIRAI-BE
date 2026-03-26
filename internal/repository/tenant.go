package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kelasmirai_backend/internal/model"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) model.TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

func (r *tenantRepository) WithTx(tx *gorm.DB) model.TenantRepository {
	return &tenantRepository{db: tx}
}

func (r *tenantRepository) Store(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&tenant).Error

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *tenantRepository) GetBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&tenant).Error

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).
		Model(&model.Tenant{}).
		Where("id = ?", tenant.ID).
		Updates(map[string]interface{}{
			"name":          tenant.Name,
			"slug":          tenant.Slug,
			"email":         tenant.Email,
			"phone":         tenant.Phone,
			"address":       tenant.Address,
			"logo_url":      tenant.LogoURL,
			"status":        tenant.Status,
			"trial_ends_at": tenant.TrialEndsAt,
			"updated_at":    gorm.Expr("NOW()"),
		}).Error
}

func (r *tenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Tenant{}).Error
}

func (r *tenantRepository) GetActiveByID(ctx context.Context, id uuid.UUID) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.WithContext(ctx).
		Where("id = ? AND status = ?", id, model.TenantStatusActive).
		First(&tenant).Error

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}
