package rest

import (
	"net/http"

	graph "github.com/graph-gophers/graphql-go"
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/boilerplate-go-project/config"
	"github.com/kubuskotak/boilerplate-go-project/ports/graphql"
	"github.com/rs/zerolog/log"
)

// Application Rest func
func Application() error {
	cfg := config.GetConfig()
	log.Info().Interface("Config", &cfg).Msg("Application rest")

	serve := bifrost.NewServerMux(bifrost.ServeOpts{
		Port: bifrost.WebPort(8077),
	})
	mux := http.NewServeMux()
	graphQuery := "/query"
	qlQuery := "/qraphql"

	opts := []graph.SchemaOpt{graph.UseFieldResolvers()}

	mux.Handle(graphQuery, bifrost.Graphql(
		graphql.Graphql,
		"schema",
		&graphql.Resolver{}, opts...))
	mux.Handle(qlQuery, bifrost.Graph(graphQuery))

	return serve.Run(mux)
}
