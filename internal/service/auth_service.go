package service

import (
	"fmt"
	"github.com/PhuMinh08082001/go-jwt-authen/common"
	"github.com/PhuMinh08082001/go-jwt-authen/common/constants"
	"github.com/PhuMinh08082001/go-jwt-authen/config"
	"github.com/PhuMinh08082001/go-jwt-authen/internal/middleware"
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
	Config            config.Configuration
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

func NewAuthService(accountRepository *repository.UserRepository, redisClient *redis.Client, config config.Configuration) *AuthService {
	return &AuthService{
		AccountRepository: accountRepository,
		RedisClient:       redisClient,
		Config:            config,
	}
}

func (service *AuthService) Login(ctx *gin.Context) {
	var cred Credentials
	if err := ctx.ShouldBindJSON(&cred); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user := service.AccountRepository.GetUser(cred.Username)

	if user.UserName != cred.Username || user.Password != cred.Password {
		ctx.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	token, err := CreateToken(user.UserName)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := service.CreateAuth(cred.Username, token)
	if saveErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	ctx.JSON(http.StatusOK, &LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

func (service *AuthService) Logout(ctx *gin.Context) {
	req := ctx.Request
	au, err := middleware.ExtractTokenMetadata(req)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorResponse{
			ErrorCode: http.StatusText(http.StatusUnauthorized),
			Code:      http.StatusUnauthorized,
		})
		return
	}

	deleted, delErr := service.DeleteAuth(au.AccessUuid)

	if delErr != nil || deleted == 0 { //if any goes wrong
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorResponse{
			ErrorCode: http.StatusText(http.StatusUnauthorized),
			Code:      http.StatusUnauthorized,
		})
		return
	}

	ctx.JSON(http.StatusOK, common.SuccessResponse{
		SuccessCode: "Logout successfully",
		Code:        http.StatusOK,
	})

}

func (service *AuthService) RefreshToken(ctx *gin.Context) {
	mapToken := map[string]string{}

	if err := ctx.ShouldBindJSON(&mapToken); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	refreshToken := mapToken["refresh_token"]
	refreshSecret := service.Config.Server.RefreshSecret

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse{
			ErrorCode: "Refresh token expired",
			Code:      http.StatusUnauthorized,
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			ctx.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userName, err := claims["user_name"].(string)
		if !ok {
			ctx.JSON(http.StatusUnprocessableEntity, err)
			return
		}

		//Delete the previous Refresh Token
		deleted, delErr := service.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			ctx.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(userName)
		if createErr != nil {
			ctx.JSON(http.StatusForbidden, createErr.Error())
			return
		}

		//save the tokens metadata to redis
		saveErr := service.CreateAuth(userName, ts)
		if saveErr != nil {
			ctx.JSON(http.StatusForbidden, saveErr.Error())
			return
		}

		tokens := map[string]string{
			constants.ACCESS_TOKEN:  ts.AccessToken,
			constants.REFRESH_TOKEN: ts.RefreshToken,
		}
		ctx.JSON(http.StatusCreated, tokens)
	} else {
		ctx.JSON(http.StatusUnauthorized, "refresh expired")
	}

}

func (service *AuthService) DeleteAuth(accessUuid string) (int64, error) {
	deleted, err := service.RedisClient.Del(accessUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func CreateToken(username string) (*TokenDetails, error) {

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	err = os.Setenv("ACCESS_SECRET", "ACCESS_SECRET")
	if err != nil {
		return nil, err
	}
	atClaims := jwt.MapClaims{}
	atClaims[constants.AUTHORIZED] = true
	atClaims[constants.ACCESS_UUID] = td.AccessUuid
	atClaims[constants.USER_NAME] = username
	atClaims[constants.EXPIRED] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "REFRESH_SECRET") //this should be in an env file
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

func (service *AuthService) CreateAuth(username string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := service.RedisClient.Set(td.AccessUuid, username, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := service.RedisClient.Set(td.RefreshUuid, username, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
