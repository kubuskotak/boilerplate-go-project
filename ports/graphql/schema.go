package graphql

import "embed"

//go:embed schema/*
var Graphql embed.FS //nolint
