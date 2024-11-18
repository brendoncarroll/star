package star

import (
	"fmt"
	"path/filepath"
	"slices"
	"strconv"

	"golang.org/x/exp/maps"
)

func NewDir(md Metadata, children map[Symbol]Command) Command {
	return Command{
		Pos: []IParam{},
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
				ctx.Printf("%s\t%s\n\n", filepath.Base(ctx.CalledAs), md.Short)
				ctx.Printf("COMMANDS:\n")
				fmtStr := "  %-" + strconv.Itoa(maxLen(keys)) + "s  %s\n"
				for _, k := range keys {
					child := children[k]
					ctx.Printf(fmtStr, k, child.Metadata.Short)
				}
				ctx.Printf("\n")
				return nil
			}
			child, ok := children[Symbol(childName)]
			if !ok {
				return fmt.Errorf("no command found for %q", childName)
			}
			return Run(ctx.Context, child, ctx.Env, childName, rest, ctx.StdIn, ctx.StdOut, ctx.StdErr)
		},
		Metadata: md,
	}
}

func maxLen[T ~string](xs []T) (ret int) {
	for _, x := range xs {
		ret = max(ret, len(x))
	}
	return ret
}
