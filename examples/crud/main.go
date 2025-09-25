package main

import (
	"fmt"

	"go.brendoncarroll.net/star"
)

func main() {
	star.Main(rootCmd)
}

var idArg = star.Required[string]{
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
		Pos:      []star.Positional{idArg},
		F: func(ctx star.Context) error {
			_, err := fmt.Fprintln(ctx.StdOut, "READ "+idArg.Load(ctx))
			return err
		},
	},
	"update": {
		Metadata: star.Metadata{Short: "update the value of an entity"},
		Pos:      []star.Positional{idArg},
		F: func(ctx star.Context) error {
			_, err := fmt.Fprintln(ctx.StdOut, "UPDATE "+idArg.Load(ctx))
			return err
		},
	},
	"delete": {
		Metadata: star.Metadata{Short: "delete an entity by id"},
		Pos:      []star.Positional{idArg},
		F: func(ctx star.Context) error {
			_, err := fmt.Fprintln(ctx.StdOut, "DELETE "+idArg.Load(ctx))
			return err
		},
	},
	"sub-dir-command": star.NewDir(
		star.Metadata{Short: "directory command below the root"},
		map[star.Symbol]star.Command{
			"c1": {
				Metadata: star.Metadata{Short: "command1"},
				F: func(c star.Context) error {
					fmt.Println("C1")
					return nil
				},
			},
			"c2": {
				Metadata: star.Metadata{Short: "command1"},
				F: func(c star.Context) error {
					fmt.Println("C1")
					return nil
				},
			},
			"echo":     echoCmd,
			"echo-pos": echoPosCmd,
		},
	),
})

var echoParam = star.Required[string]{
	Name:  "param",
	Parse: star.ParseString,
}

var echoPosCmd = star.Command{
	Metadata: star.Metadata{Short: "echos back the args"},
	Pos:      []star.Positional{echoParam},
	F: func(c star.Context) error {
		c.Printf("%s\n", echoParam.Load(c))
		return nil
	},
}

var echoCmd = star.Command{
	Metadata: star.Metadata{Short: "echos back the args"},
	F: func(c star.Context) error {
		c.Printf("EXTRA")
		for i, arg := range c.Extra {
			c.Printf("%d %s\n", i, arg)
		}
		return nil
	},
}
