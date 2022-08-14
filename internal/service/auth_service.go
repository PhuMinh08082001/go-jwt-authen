package service

import (
	"fmt"
	"github.com/PhuMinh08082001/go-jwt-authen/internal/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AuthService struct {
	AccountRepository *repository.UserRepository
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserName string `json:"user_name"`
}

func NewAuthService(accountRepository *repository.UserRepository) *AuthService {
	return &AuthService{
		AccountRepository: accountRepository,
	}
}

func (service *AuthService) Login(ctx *gin.Context) {
	var cred Credentials
	if err := ctx.ShouldBindJSON(&cred); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cred)

	ctx.JSON(http.StatusOK, &LoginResponse{
		Token: "Success",
	})
}
