package main

import (
	"flag"
	"fmt"
	"go-workshop/ports/rest"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	flags = flag.NewFlagSet("baldr", flag.ExitOnError)
	help  = flags.Bool("h", false, "print help")
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
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
}

var (
	usagePrefix = `Usage: goWork [OPTIONS] COMMAND

Examples:
	goWork rest

Options:`
)
