package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubuskotak/boilerplate-go-project/ports/grpc"
	"github.com/kubuskotak/boilerplate-go-project/ports/rest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	flags = flag.NewFlagSet("baldr", flag.ExitOnError)
	help  = flags.Bool("h", false, "print help")
	//version = flags.Bool("version", false, "print version")
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Error().Err(err)
		return
	}

	args := flags.Args()
	if len(args) == 0 || *help {
		flags.Usage()
		return
	}

	switch args[0] {
	case "rest":
		if err := rest.Application(); err != nil {
			log.Error().Err(err)
		}
		return
	case "grpc":
		if err := grpc.Application(); err != nil {
			log.Error().Err(err)
		}
		return
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
}

var (
	usagePrefix = `Usage: baldr [OPTIONS] COMMAND

Examples:
	baldr rest
	baldr event-store
	baldr dispatcher

Options:`
)
