package rest

import (
	"net/http"

	graph "github.com/graph-gophers/graphql-go"
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/boilerplate-go-project/ports/graphql"
)

// Application Rest func
func Application() error {
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
