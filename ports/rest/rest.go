package rest

import (
	"net/http"

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
	mux.Handle(graphQuery, bifrost.Graphql(graphql.Graphql, "schema", &graphql.Resolver{}))
	mux.Handle(qlQuery, bifrost.Graph(graphQuery))

	return serve.Run(mux)
}
