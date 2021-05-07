package rest

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/boilerplate-go-project/config"
	"github.com/kubuskotak/valkyrie"
	"github.com/opentracing/opentracing-go"
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

	tracer, cleanup, err := valkyrie.Tracer("hello", "0.2.0")
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	r := chi.NewRouter()
	r.Group(func(c chi.Router) {
		graphql := Graphql{}
		graphql.Register(ctx, c)
	})

	r.Group(func(c chi.Router) {
		hello := &Hello{Tracer: tracer}
		hello.Register(ctx, c)
	})

	errServer := serve.Run(r)
	cleanup()
	return errServer
}
