package config

import (
	"kelasmirai_backend/internal/controller"
	"kelasmirai_backend/internal/delivery/http/route"
	"kelasmirai_backend/internal/repository"
	"kelasmirai_backend/internal/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Database    *gorm.DB
	Router      *mux.Router
	HttpClient  *http.Client
	Validate    *validator.Validate
	RedisClient *redis.Client
}

func NewBootstrapConfig(bsCnf *BootstrapConfig) {

	// ==================== dependency injection
	emailVerifyRepo := repository.NewEmailVerificationTokenRepository(bsCnf.Database)
	// invoiceRepo := repository.NewInvoiceRepository(bsCnf.Database)
	passResetCodeRepo := repository.NewPasswordResetCodeRepository(bsCnf.Database)
	// paymentRepo := repository.NewPaymentRepository(bsCnf.Database)
	subsPlanRepo := repository.NewSubscriptionPlanRepository(bsCnf.Database)
	tenantSubsRepo := repository.NewTenantSubscriptionRepository(bsCnf.Database)
	userRepo := repository.NewUserRepository(bsCnf.Database)
	// webhookLogRepo := repository.NewWebhookLogRepository(bsCnf.Database)

	authUsecase := usecase.NewAuthUsecase(bsCnf.Database, userRepo, subsPlanRepo, passResetCodeRepo, emailVerifyRepo, tenantSubsRepo)

	authController := controller.NewAuthController(authUsecase)

	routeConfig := route.RouteConfig{
		Router:         bsCnf.Router,
		AuthController: authController,
	}
	routeConfig.Setup()
}
