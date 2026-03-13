package controller

import "kelasmirai_backend/internal/usecase"

type AuthController struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthController(authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: &authUsecase,
	}
}
