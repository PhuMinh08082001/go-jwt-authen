package routes

import (
	"github.com/PhuMinh08082001/go-jwt-authen/internal/controller"
	"github.com/PhuMinh08082001/go-jwt-authen/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

type RouteParams struct {
	fx.In
	Route          *gin.Engine
	Middleware     *middleware.Middleware
	AuthController *controller.AuthController
}
type Hello struct {
	Message string
}

func InitAccountRoute(params RouteParams) {

	route := params.Route
	authController := params.AuthController
	middlewares := params.Middleware

	indexRoute := route.Group("/")
	{
		indexRoute.POST("/login", authController.Login)
	}

	accountRoute := route.Group("/hello")
	accountRoute.Use(middlewares.JWT)
	{
		accountRoute.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, &Hello{
				Message: "Hello World",
			})
		})
	}

}
