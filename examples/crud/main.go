package main

import (
	"fmt"

	"go.brendoncarroll.net/star"
)

func main() {
	star.Main(rootCmd)
}

var idArg = star.Required[string]{
	ID:       "id",
	Parse:    star.ParseString,
	ShortDoc: "the id of the entity",
}

var rootCmd = star.NewDir(star.Metadata{Short: "an example CLI app"}, map[string]star.Command{
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
		map[string]star.Command{
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
	"grouped-subc": star.NewGroupedDir(
		star.Metadata{Short: "grouped directory command"},
		[]star.Group{
			{
				Title: "Even Commands",
				Commands: []string{
					"c0",
					"c2",
				},
			},
			{
				Title: "Odd Commands",
				Commands: []string{
					"c1",
					"c3",
				},
			},
		},
		map[string]star.Command{
			"c0": echoCmd,
			"c1": echoCmd,
			"c2": echoCmd,
			"c3": echoCmd,
		},
	),
})

var echoParam = star.Required[string]{
	ID:       "param",
	Parse:    star.ParseString,
	ShortDoc: "the value to echo back",
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
		c.Printf("EXTRA:\n")
		for i, arg := range c.Extra {
			c.Printf("%d %s\n", i, arg)
		}
		return nil
	},
}
