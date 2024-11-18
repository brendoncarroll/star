package main

import (
	"fmt"

	"go.brendoncarroll.net/star"
)

func main() {
	star.Main(rootCmd)
}

var idArg = star.Param[string]{
	Name:  "id",
	Parse: star.ParseString,
}

var rootCmd = star.NewDir(star.Metadata{Short: "an example CLI app"}, map[star.Symbol]star.Command{
	"create": {
		Metadata: star.Metadata{Short: "creates a new entity"},
		F: func(ctx star.Context) error {
			_, err := fmt.Fprintln(ctx.StdOut, "CREATE")
			return err
		},
	},
	"read": {
		Metadata: star.Metadata{Short: "reads the value of an entity"},
		Pos:      []star.IParam{idArg},
		F: func(ctx star.Context) error {
			_, err := fmt.Fprintln(ctx.StdOut, "READ "+idArg.Load(ctx))
			return err
		},
	},
	"update": {
		Metadata: star.Metadata{Short: "update the value of an entity"},
		Pos:      []star.IParam{idArg},
		F: func(ctx star.Context) error {
			_, err := fmt.Fprintln(ctx.StdOut, "UPDATE "+idArg.Load(ctx))
			return err
		},
	},
	"delete": {
		Metadata: star.Metadata{Short: "delete an entity by id"},
		Pos:      []star.IParam{idArg},
		F: func(ctx star.Context) error {
			_, err := fmt.Fprintln(ctx.StdOut, "DELETE "+idArg.Load(ctx))
			return err
		},
	},
})
