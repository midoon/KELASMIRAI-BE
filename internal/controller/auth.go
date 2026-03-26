package controller

import (
	"encoding/json"
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

	request := dto.TenantRegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, dto.MessageResponse{
			Status:  false,
			Message: "invalid request body",
		})
		return
	}

	if err := ac.authUsecase.Register(ctx, request); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, dto.MessageResponse{
			Status:  false,
			Message: "failed to register tenant",
		})
		return
	}

	res := dto.MessageResponse{
		Status:  true,
		Message: "oke",
	}
	helper.WriteJSON(w, http.StatusOK, res)
}
