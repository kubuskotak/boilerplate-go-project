package grpc

import (
	"context"
	"fmt"

	person "github.com/kubuskotak/boilerplate-go-project/ports/grpc/proto"
	"github.com/rs/zerolog/log"
)

func (ps PersonServer) List(ctx context.Context, r *person.ListPersonsRequest) (*person.ListPersonsResponse, error) {
	persons := make([]*person.Person, 0)
	for i := 1; i < 10; i++ {
		p := &person.Person{}
		p.Id = fmt.Sprintf("id#%d", i)
		p.Name = fmt.Sprintf("name#%d", i)
		p.Email = fmt.Sprintf("email#%d", i)
		persons = append(persons, p)
	}

	log.Info().Msg("Person List")

	return &person.ListPersonsResponse{Persons: persons}, nil
}
