package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	_ "github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"strings"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AccessDetails struct {
	AccessUuid string
	UserName   string
}

type ErrorResponse struct {
	ErrorCode string
	Code      int
}

func (middleware *Middleware) JWT(ctx *gin.Context) {
	tokenAuth, err := ExtractTokenMetadata(ctx.Request)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
			ErrorCode: http.StatusText(http.StatusUnauthorized),
			Code:      http.StatusUnauthorized,
		})
		return
	}
	_, err = FetchAuth(tokenAuth, middleware.Client)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
			ErrorCode: http.StatusText(http.StatusUnauthorized),
			Code:      http.StatusUnauthorized,
		})
		return
	}
	log.Println(tokenAuth.UserName)
	ctx.Next()
}

func FetchAuth(authD *AccessDetails, client *redis.Client) (string, error) {
	userName, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}

	return userName, nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userName, ok := claims["user_name"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserName:   userName,
		}, nil
	}
	return nil, err
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
