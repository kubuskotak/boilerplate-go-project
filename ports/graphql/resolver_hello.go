package graphql

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/kubuskotak/boilerplate-go-project/config"
	person "github.com/kubuskotak/boilerplate-go-project/ports/grpc/proto"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func (r *Resolver) Hello(ctx context.Context) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "resolver.hello")
	defer span.Finish()

	val := "Hello Graphql"
	span.SetTag("resolver.result", val)

	return val
}

func (r *Resolver) SetHello(ctx context.Context, args struct{ Name string }) []*person.Person {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "resolver.setHello")
	defer span.Finish()
	cfg := config.GetConfig()
	conn, err := grpc.DialContext(
		spanContext,
		fmt.Sprintf("localhost:%d", cfg.Port.Grpc),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithInsecure(),
	)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	if err != nil {
		log.Error().Err(err)
	}
	persons, err := person.NewPersonServiceClient(conn).List(ctx, &person.ListPersonsRequest{PageSize: 1})
	if err != nil {
		log.Error().Err(err)
	}
	log.Info().Msg("aneeeeehhhh")
	log.Info().Interface("Persons", persons.Persons).Msg("kudu ada msg")
	return persons.Persons
}
