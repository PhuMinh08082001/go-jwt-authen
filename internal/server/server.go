package server

import (
	"context"
	"github.com/PhuMinh08082001/go-jwt-authen/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
)

var Module = fx.Invoke(registerServer)

func registerServer(lifecycle fx.Lifecycle, config config.Configuration, route *gin.Engine) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Println("Starting Http Server")
				go func() {
					err := route.Run(config.Server.Address)
					if err != nil {
						log.Fatal(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Http Server Stop ...")
				return nil
			},
		})
}
