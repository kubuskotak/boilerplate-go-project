package rest

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/boilerplate-go-project/config"
)

// Application Rest func
func Application() error {
	cfg := config.GetConfig()
	serve := bifrost.NewServerMux(bifrost.ServeOpts{
		Port: bifrost.WebPort(cfg.Port.Http),
	})
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.App.Latency)*time.Millisecond,
	)
	defer cancel()

	r := chi.NewRouter()
	hello := &Hello{}
	hello.Register(ctx, r)

	graphql := Graphql{}
	graphql.Register(ctx, r)

	return serve.Run(r)
}
