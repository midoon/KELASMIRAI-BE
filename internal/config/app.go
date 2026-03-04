package config

import (
	"kelasmirai_backend/internal/delivery/http/route"
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
	routeConfig := route.RouteConfig{
		Router: bsCnf.Router,
	}
	routeConfig.Setup()
}
