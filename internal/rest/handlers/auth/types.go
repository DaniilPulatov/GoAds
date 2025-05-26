package auth

import "ads-service/internal/usecase/auth"

type AuthHandler struct {
	userAuthService auth.AuthService
}

func NewAuthHandler(userAuthService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		userAuthService: userAuthService,
	}
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
