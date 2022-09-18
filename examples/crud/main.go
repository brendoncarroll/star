package main

import (
	"fmt"

	"github.com/brendoncarroll/star"
)

func main() {
	star.Main("crud", rootCmd)
}

var idArg = star.NewPos[string]("id", "", true)

var rootCmd = star.NewParent("an example CLI app", nil, map[string]*star.Command{
	"create": &star.Command{
		Short: "creates a new entity",
		F: func(ctx *star.Context) error {
			_, err := fmt.Fprintln(ctx.Out, "CREATE")
			return err
		},
	},
	"read": &star.Command{
		Short: "reads the value of an entity",
		Pos:   []star.Pos{idArg},
		F: func(ctx *star.Context) error {
			_, err := fmt.Fprintln(ctx.Out, "READ "+ctx.String("id"))
			return err
		},
	},
	"update": &star.Command{
		Short: "update the value of an entity",
		Pos:   []star.Pos{idArg},
		F: func(ctx *star.Context) error {
			_, err := fmt.Fprintln(ctx.Out, "UPDATE "+ctx.String("id"))
			return err
		},
	},
	"delete": &star.Command{
		Short: "delete an entity by id",
		Pos:   []star.Pos{idArg},
		F: func(ctx *star.Context) error {
			_, err := fmt.Fprintln(ctx.Out, "DELETE "+ctx.String("id"))
			return err
		},
	},
})
