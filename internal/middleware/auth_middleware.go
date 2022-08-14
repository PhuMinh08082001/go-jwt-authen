package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/labstack/echo/v4/middleware"
	"log"
)

var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (middleware *Middleware) JWT(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	var cred Credentials
	log.Printf("Hello %s", authorization)

	if err := ctx.ShouldBindJSON(&cred); err != nil {

		return
	}
	fmt.Println(cred)

}
