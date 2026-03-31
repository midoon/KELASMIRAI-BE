package controller

import (
	"kelasmirai_backend/internal/dto"
	"kelasmirai_backend/internal/helper"
	"kelasmirai_backend/internal/usecase"
	"net/http"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthController(authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, dto.MessageResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	request := dto.TenantRegisterRequest{
		SchoolName:    r.FormValue("school_name"),
		SchoolSlug:    r.FormValue("school_slug"),
		SchoolEmail:   r.FormValue("school_email"),
		SchoolPhone:   r.FormValue("school_phone"),
		SchoolAddress: r.FormValue("school_address"),
		AdminName:     r.FormValue("admin_name"),
		AdminEmail:    r.FormValue("admin_email"),
		AdminPassword: r.FormValue("admin_password"),
	}

	if err := ac.authUsecase.Register(ctx, request); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, dto.MessageResponse{
			Status:  false,
			Message: "failed to register tenant: " + err.Error(),
		})
		return
	}

	res := dto.MessageResponse{
		Status:  true,
		Message: "oke",
	}
	helper.WriteJSON(w, http.StatusOK, res)
}

func (ac *AuthController) VerifyRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.URL.Query().Get("token")

	if err := ac.authUsecase.VerifyRegistration(ctx, token); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, dto.MessageResponse{
			Status:  false,
			Message: "invalid or expired token: " + err.Error(),
		})
		return
	}

	helper.WriteJSON(w, http.StatusOK, dto.MessageResponse{
		Status:  true,
		Message: "Email verified. You can now log in.",
	})

}
