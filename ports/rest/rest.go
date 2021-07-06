package rest

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go-workshop/config"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Application Rest func
func Application() error {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.App.Latency)*time.Millisecond,
	)
	defer cancel()

	r := chi.NewRouter()
	r.Group(func(c chi.Router) {
		hello := &Hello{}
		hello.Register(ctx, c)
	})

	log.Info().Msg(fmt.Sprintf("http server running:%d", cfg.Port.Http))

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port.Http), r)
}
