package route

import (
	"kelasmirai_backend/internal/controller"
	"net/http"

	"github.com/gorilla/mux"
)

type RouteConfig struct {
	Router *mux.Router

	AuthController *controller.AuthController
}

func (rc *RouteConfig) Setup() {

	rc.setupPublicRoute()
	rc.setupPrivateRoute()
}

func (rc *RouteConfig) setupPublicRoute() {

	rc.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Health route"))
	}).Methods("GET")
}

func (rc *RouteConfig) setupPrivateRoute() {

	// rc.Router.HandleFunc("/private", privateHandler).Methods("GET")
}
