package controller

import (
	"github.com/PhuMinh08082001/go-jwt-authen/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

type TokenResponse struct {
	Message  string        `json:"message"`
	Response *TokenWrapper `json:"response"`
}

type TokenWrapper struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	AccessToken   string `json:"access_token"`
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		authService: service,
	}
}

func (controller *AuthController) Login(ctx *gin.Context) {
	controller.authService.Login(ctx)
}
