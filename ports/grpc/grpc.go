package grpc

import (
	"github.com/kubuskotak/bifrost"
	person "github.com/kubuskotak/boilerplate-go-project/ports/grpc/proto"
	"google.golang.org/grpc"
)

type PersonServer struct {
	person.PersonServiceServer
}

// Application Rest func
func Application() error {
	ps := PersonServer{}

	serve := bifrost.NewServerGRPC(bifrost.GRPCOpts{
		Port: bifrost.GRPCPort(5077),
	})

	return serve.Run(func(s *grpc.Server) {
		person.RegisterPersonServiceServer(s, ps)
	})
}
