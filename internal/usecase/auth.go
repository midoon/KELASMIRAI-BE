package usecase

import (
	"context"
	"kelasmirai_backend/internal/dto"
	"kelasmirai_backend/internal/model"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(ctx context.Context, request dto.TenantRegisterRequest) error
}

type authUsecase struct {
	db             *gorm.DB
	userRepo       model.UserRepository
	subsPlanrepo   model.SubscriptionPlanRepository
	passResetRepo  model.PasswordResetCodeRepository
	emailVerifRepo model.EmailVerificationTokenRepository
	tenantSubsRepo model.TenantSubscriptionRepository
}

func NewAuthUsecase(
	db *gorm.DB,
	userRepo model.UserRepository,
	subsPlanrepo model.SubscriptionPlanRepository,
	passResetRepo model.PasswordResetCodeRepository,
	emailVerifRepo model.EmailVerificationTokenRepository,
	tenantSubsRepo model.TenantSubscriptionRepository,
) AuthUsecase {
	return &authUsecase{
		db:             db,
		userRepo:       userRepo,
		subsPlanrepo:   subsPlanrepo,
		passResetRepo:  passResetRepo,
		emailVerifRepo: emailVerifRepo,
		tenantSubsRepo: tenantSubsRepo,
	}
}

func (au *authUsecase) Register(ctx context.Context, request dto.TenantRegisterRequest) error {

	plan, err := au.subsPlanrepo.GetAll(ctx)
	if err != nil {
		return err
	}

	freePlanId := plan[0].ID

	for _, p := range plan {
		if p.PriceYearly.IsZero() && p.PriceMonthly.IsZero() {
			freePlanId = p.ID
			break
		}
	}

	// Create tenant
	tenant := model.Tenant{
		Name:        request.SchoolName,
		Slug:        request.SchoolSlug,
		Email:       request.SchoolEmail,
		Phone:       &request.SchoolPhone,
		Address:     &request.SchoolAddress,
		Status:      model.TenantStatusActive,
		TrialEndsAt: time.Now().AddDate(0, 0, 30),
	}

	// create tenant subscription with default plan
	tenantSubscription := model.TenantSubscription{
		PlanID:                 freePlanId,
		BillingCycle:           "monthly",
		Price:                  decimal.NewFromInt(0),
		Status:                 model.SubscriptionStatusTrial,
		StartedAt:              time.Now(),
		NextBillingAt:          time.Now().AddDate(0, 0, 30),
		EndedAt:                nil,
		MidtransSubscriptionID: nil,
	}

	// create user with role tenant admin
	user := model.User{
		Name:         request.AdminName,
		Email:        request.AdminEmail,
		PasswordHash: request.AdminPassword,
		Role:         model.UserRoleAdmin,
		IsActive:     true,
	}

	tx := au.db.WithContext(ctx).Begin()
	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		return err
	}

	tenantSubscription.TenantID = tenant.ID
	if err := tx.Create(&tenantSubscription).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.TenantID = tenant.ID
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
