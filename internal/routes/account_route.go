package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouteParams struct {
	fx.In
	Route *gin.Engine
}

func InitAccountRoute(params RouteParams) {
	route := params.Route
	accountRoute := route.Group("/users")
	{
		accountRoute.GET("/login")
	}
}
