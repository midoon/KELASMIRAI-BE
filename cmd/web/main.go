package main

import (
	"fmt"
	"kelasmirai_backend/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)

	r := mux.NewRouter()

	httpClient := http.Client{}

	config.NewBootstrapConfig(&config.BootstrapConfig{
		Database:   db,
		Router:     r,
		HttpClient: &httpClient,
	})

	// handler := middleware.CorsMiddleware(r) // harus diassign di awal, biar kepanggil pertamakali

	// addr := fmt.Sprintf("%s:%s", cnf.Server.Host, cnf.Server.Port)
	// server := &http.Server{
	// 	Addr:    addr,
	// 	Handler: handler,
	// }

	addr := fmt.Sprintf("%s:%s", viperConfig.GetString("web.host"), viperConfig.GetString("web.port"))
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Println("Server is running on port", addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
