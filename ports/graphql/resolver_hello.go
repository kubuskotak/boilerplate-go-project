package graphql

import (
	"context"
	"time"

	person "github.com/kubuskotak/boilerplate-go-project/ports/grpc/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func (r *Resolver) Hello(ctx context.Context) string {
	return "Hello Graphql"
}

func (r *Resolver) SetHello(ctx context.Context, args struct{ Name string }) []*person.Person {
	conn, err := grpc.Dial(
		"localhost:5077",
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second),
	)

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
