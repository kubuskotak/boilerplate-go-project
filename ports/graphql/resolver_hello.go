package graphql

import (
	"context"
	"fmt"
)

func (r *Resolver) Hello(ctx context.Context) string {
	return "Hello Graphql"
}

func (r *Resolver) SetHello(ctx context.Context, args struct{ Name string }) string {
	return fmt.Sprintf("Hello meet again, %s", args.Name)
}
