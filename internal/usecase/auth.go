package usecase

import "kelasmirai_backend/internal/model"

type AuthUsecase interface{}

type authUsecase struct {
	userRepo       model.UserRepository
	subsPlanrepo   model.SubscriptionPlanRepository
	passResetRepo  model.PasswordResetCodeRepository
	emailVerifRepo model.EmailVerificationTokenRepository
	tenantSubsRepo model.TenantSubscriptionRepository
}

func NewAuthUsecase(
	userRepo model.UserRepository,
	subsPlanrepo model.SubscriptionPlanRepository,
	passResetRepo model.PasswordResetCodeRepository,
	emailVerifRepo model.EmailVerificationTokenRepository,
	tenantSubsRepo model.TenantSubscriptionRepository,
) AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		subsPlanrepo:   subsPlanrepo,
		passResetRepo:  passResetRepo,
		emailVerifRepo: emailVerifRepo,
		tenantSubsRepo: tenantSubsRepo,
	}
}
