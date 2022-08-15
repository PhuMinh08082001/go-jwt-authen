package middleware

import (
	"github.com/go-redis/redis/v7"
	"go.uber.org/fx"
)

type Middleware struct {
	Client *redis.Client
}

var Module = fx.Provide(NewMiddleware)

func NewMiddleware(client *redis.Client) *Middleware {
	return &Middleware{
		Client: client,
	}
}
