package usecase

import (
	"context"
	"fmt"
	"kelasmirai_backend/internal/dto"
	"kelasmirai_backend/internal/helper"
	"kelasmirai_backend/internal/model"
	"kelasmirai_backend/internal/util"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(ctx context.Context, request dto.TenantRegisterRequest) error
	VerifyRegistration(ctx context.Context, token string) error
}

type authUsecase struct {
	db             *gorm.DB
	viperCnf       *viper.Viper
	validate       *validator.Validate
	userRepo       model.UserRepository
	subsPlanrepo   model.SubscriptionPlanRepository
	passResetRepo  model.PasswordResetCodeRepository
	emailVerifRepo model.EmailVerificationTokenRepository
	tenantSubsRepo model.TenantSubscriptionRepository
	tenantRepo     model.TenantRepository
}

func NewAuthUsecase(
	db *gorm.DB,
	viperCnf *viper.Viper,
	validate *validator.Validate,
	userRepo model.UserRepository,
	subsPlanrepo model.SubscriptionPlanRepository,
	passResetRepo model.PasswordResetCodeRepository,
	emailVerifRepo model.EmailVerificationTokenRepository,
	tenantSubsRepo model.TenantSubscriptionRepository,
	tenantRepo model.TenantRepository,
) AuthUsecase {
	return &authUsecase{
		db:             db,
		userRepo:       userRepo,
		viperCnf:       viperCnf,
		validate:       validate,
		subsPlanrepo:   subsPlanrepo,
		passResetRepo:  passResetRepo,
		emailVerifRepo: emailVerifRepo,
		tenantSubsRepo: tenantSubsRepo,
		tenantRepo:     tenantRepo,
	}
}

func (au *authUsecase) Register(ctx context.Context, request dto.TenantRegisterRequest) error {
	if err := au.validate.Struct(request); err != nil {
		return helper.NewCustomError(http.StatusBadRequest, "validation error", err)
	}

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
	}

	var isTenantExists *model.Tenant

	isTenantExists, err = au.tenantRepo.GetByEmail(ctx, request.SchoolEmail)

	if isTenantExists != nil {
		return helper.NewCustomError(http.StatusConflict, "tenant with this email already exists", err)
	}

	isTenantExists, err = au.tenantRepo.GetBySlug(ctx, request.SchoolSlug)

	if isTenantExists != nil {
		return helper.NewCustomError(http.StatusConflict, "tenant with this slug already exists", err)
	}

	tx := au.db.WithContext(ctx).Begin()
	userRepoTx := au.userRepo.WithTx(tx)
	tenantRepoTx := au.tenantRepo.WithTx(tx)
	tenantSubscriptionRepoTx := au.tenantSubsRepo.WithTx(tx)

	if err := tenantRepoTx.Store(ctx, &tenant); err != nil {
		tx.Rollback()
		return helper.NewCustomError(http.StatusInternalServerError, "failed to create tenant", err)
	}

	tenantSubscription.TenantID = tenant.ID

	if err := tenantSubscriptionRepoTx.Store(ctx, &tenantSubscription); err != nil {
		tx.Rollback()
		return helper.NewCustomError(http.StatusInternalServerError, "failed to create tenant subscription", err)
	}

	user.TenantID = tenant.ID
	if err := userRepoTx.Store(ctx, &user); err != nil {
		tx.Rollback()
		return helper.NewCustomError(http.StatusInternalServerError, "failed to create user", err)
	}

	// create and store email verification token
	token := helper.GenerateRandomString(9)
	emailVerificationToken := model.EmailVerificationToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := au.emailVerifRepo.WithTx(tx).Store(ctx, &emailVerificationToken); err != nil {
		tx.Rollback()
		return helper.NewCustomError(http.StatusInternalServerError, "failed to create email verification token", err)
	}

	to := request.AdminEmail
	subject := "KelasMirai - Email Verification"
	templatePath := filepath.Join("internal", "template", "email", "registration-confirmation.html")
	data := struct {
		Name             string
		VerificationLink string
		Year             int
	}{
		Name:             user.Name,
		VerificationLink: fmt.Sprintf("%s/auth/tenant/verify-registration?token=%s", au.viperCnf.GetString("services.frontend.url"), token),
		Year:             time.Now().Year(),
	}
	emailCnf := au.viperCnf

	if err := util.SendEmail(to, subject, templatePath, data, emailCnf); err != nil {
		tx.Rollback()
		return helper.NewCustomError(http.StatusInternalServerError, "failed to send verification email", err)
	}

	if err := tx.Commit().Error; err != nil {
		return helper.NewCustomError(http.StatusInternalServerError, "failed to commit transaction", err)
	}

	return nil
}

func (au *authUsecase) VerifyRegistration(ctx context.Context, token string) error {
	emailVerifToken, err := au.emailVerifRepo.GetByToken(ctx, token)
	if err != nil {
		return helper.NewCustomError(http.StatusBadRequest, "invalid or expired token", err)
	}

	if emailVerifToken.UsedAt != nil || emailVerifToken.ExpiresAt.Before(time.Now()) {
		return helper.NewCustomError(http.StatusBadRequest, "invalid or expired token", nil)
	}

	tx := au.db.WithContext(ctx).Begin()
	emailVerifRepoTx := au.emailVerifRepo.WithTx(tx)
	userRepoTx := au.userRepo.WithTx(tx)

	now := time.Now()
	emailVerifToken.UsedAt = &now

	if err := emailVerifRepoTx.Update(ctx, emailVerifToken); err != nil {
		return helper.NewCustomError(http.StatusInternalServerError, "failed to verify email", err)
	}

	user, err := userRepoTx.GetByID(ctx, emailVerifToken.UserID)
	if err != nil {
		tx.Rollback()
		return helper.NewCustomError(http.StatusInternalServerError, "failed to verify email", err)
	}

	user.ActivatedAt = &now

	if err := userRepoTx.Update(ctx, user); err != nil {
		tx.Rollback()
		return helper.NewCustomError(http.StatusInternalServerError, "failed to verify email", err)
	}

	if err := tx.Commit().Error; err != nil {
		return helper.NewCustomError(http.StatusInternalServerError, "failed to verify email", err)
	}
	return nil
}
