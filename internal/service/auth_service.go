package service

import (
	"fmt"
	"github.com/PhuMinh08082001/go-jwt-authen/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	"net/http"
	"os"
	"time"
)

type AuthService struct {
	AccountRepository *repository.UserRepository
	RedisClient       *redis.Client
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func NewAuthService(accountRepository *repository.UserRepository, redisClient *redis.Client) *AuthService {
	return &AuthService{
		AccountRepository: accountRepository,
		RedisClient:       redisClient,
	}
}

func (service *AuthService) Login(ctx *gin.Context) {
	var cred Credentials
	if err := ctx.ShouldBindJSON(&cred); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user := service.AccountRepository.GetUser(cred.Username)

	fmt.Println(user)
	if user.UserName != cred.Username || user.Password != cred.Password {
		ctx.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	token, err := CreateToken(user.UserName)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := CreateAuth(cred.Username, token, service.RedisClient)
	if saveErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	ctx.JSON(http.StatusOK, &LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

func CreateToken(username string) (*TokenDetails, error) {

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	err = os.Setenv("ACCESS_SECRET", "ramping-around-jwt")
	if err != nil {
		return nil, err
	}
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_name"] = username
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_name"] = username
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(username string, td *TokenDetails, client *redis.Client) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, username, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, username, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
