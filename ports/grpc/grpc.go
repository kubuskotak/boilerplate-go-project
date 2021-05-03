package grpc

import (
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/boilerplate-go-project/config"
	person "github.com/kubuskotak/boilerplate-go-project/ports/grpc/proto"
	"google.golang.org/grpc"
)

type PersonServer struct {
	person.PersonServiceServer
}

// Application Rest func
func Application() error {
	cfg := config.GetConfig()
	ps := PersonServer{}

	serve := bifrost.NewServerGRPC(bifrost.GRPCOpts{
		Port: bifrost.GRPCPort(cfg.Port.Grpc),
	})

	return serve.Run(func(s *grpc.Server) {
		person.RegisterPersonServiceServer(s, ps)
	})
}
