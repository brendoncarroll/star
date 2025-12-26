package star

import (
	"fmt"
	"path/filepath"
	"slices"
	"strconv"

	"golang.org/x/exp/maps"
)

func NewDir(md Metadata, children map[string]Command) Command {
	return Command{
		Pos:   []Positional{},
		Flags: map[string]Flag{},
		F: func(ctx Context) error {
			var childName string
			var rest []string
			for i := 0; i < len(ctx.Extra); i++ {
				arg := ctx.Extra[i]
				if isFlag(arg) {
					i++
					continue
				}
				childName, rest = arg, slices.Delete(ctx.Extra, i, i+1)
				break
			}

			if childName == "" {
				keys := maps.Keys(children)
				slices.Sort(keys)
				ctx.Printf("%s\n\n", filepath.Base(ctx.CalledAs))
				ctx.Printf("%s\n\n", md.Short)
				ctx.Printf("COMMANDS:\n")
				fmtStr := "  %-" + strconv.Itoa(maxLen(keys)) + "s  %s\n"
				for _, k := range keys {
					child := children[k]
					ctx.Printf(fmtStr, k, child.Metadata.Short)
				}
				ctx.Printf("\n")
				return nil
			}
			child, ok := children[childName]
			if !ok {
				return fmt.Errorf("no command found for %q", childName)
			}
			return Run(ctx.Context, child, ctx.Env, childName, rest, ctx.StdIn, ctx.StdOut, ctx.StdErr)
		},
		Metadata: md,
	}
}

// Group is a named set of Commands presented together.
// Each command name in the list is assumed to be unique.
// The list of commands will be sorted, so the order is not meaningful.
type Group struct {
	Title    string
	Commands []string
}

func NewGroupedDir(md Metadata, groups []Group, children map[string]Command) Command {
	return Command{
		F: func(ctx Context) error {
			var childName string
			var rest []string
			for i := 0; i < len(ctx.Extra); i++ {
				arg := ctx.Extra[i]
				if isFlag(arg) {
					i++
					continue
				}
				childName, rest = arg, slices.Delete(ctx.Extra, i, i+1)
				break
			}

			if childName == "" {
				ctx.Printf("%s\n\n", filepath.Base(ctx.CalledAs))
				ctx.Printf("%s\n\n", md.Short)
				for _, g := range groups {
					ctx.Printf("%s:\n", g.Title)
					slices.Sort(g.Commands)
					fmtStr := "  %-" + strconv.Itoa(maxLen(g.Commands)) + "s  %s\n"
					for _, cmdName := range g.Commands {
						child, ok := children[cmdName]
						if !ok {
							panic(fmt.Sprintf("No child command %q exists.  This is a bug.", cmdName))
						}
						ctx.Printf(fmtStr, cmdName, child.Metadata.Short)
					}
					ctx.Printf("\n")
				}
				return nil
			} else {
				child, ok := children[childName]
				if !ok {
					return fmt.Errorf("no command found for %q", childName)
				}
				return Run(ctx.Context, child, ctx.Env, childName, rest, ctx.StdIn, ctx.StdOut, ctx.StdErr)
			}
		}}
}

func maxLen[T ~string](xs []T) (ret int) {
	for _, x := range xs {
		ret = max(ret, len(x))
	}
	return ret
}
