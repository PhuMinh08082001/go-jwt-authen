package routes

import (
	"github.com/PhuMinh08082001/go-jwt-authen/internal/controller"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouteParams struct {
	fx.In
	Route          *gin.Engine
	AuthController *controller.AuthController
}
type Hello struct {
	Message string
}

func InitAccountRoute(params RouteParams) {
	route := params.Route
	authController := params.AuthController
	indexRoute := route.Group("/")
	{
		indexRoute.POST("/login", authController.Login)
	}

	//accountRoute := route.Group("/hello")
	//accountRoute.Use(middlewares.JWT)
	//{
	//	accountRoute.POST("/", func(context *gin.Context) {
	//		context.JSON(http.StatusOK, &Hello{
	//			Message: "You are signed in",
	//		})
	//	})
	//}

}
