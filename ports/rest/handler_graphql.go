package rest

import (
	"context"

	"github.com/go-chi/chi/v5"
	graph "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/trace"
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/boilerplate-go-project/ports/graphql"
)

type Graphql struct {
}

func (g *Graphql) Register(ctx context.Context, router chi.Router) {
	opts := []graph.SchemaOpt{graph.Tracer(&trace.OpenTracingTracer{}), graph.UseFieldResolvers()}
	router.Handle("/query", bifrost.Graphql(
		graphql.Graphql,
		"schema",
		&graphql.Resolver{}, opts...))
	router.Handle("/graphql", bifrost.Graph("/query"))
}
