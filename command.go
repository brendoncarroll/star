package star

import (
	"fmt"
	"strings"
)

type Metadata struct {
	Short string
	// Tags for grouping by category
	Tags []string
}

type Command struct {
	Flags []Flag
	Pos   []Positional
	F     func(c Context) error
	Metadata
}

func (c Command) HasParam(x Symbol) bool {
	for i := range c.Pos {
		if c.Pos[i].name() == x {
			return true
		}
	}
	for i := range c.Flags {
		if c.Flags[i].name() == x {
			return true
		}
	}
	return false
}

func (c Command) Doc(calledAs string) string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "%s ", calledAs)
	for _, pos := range c.Pos {
		sb.WriteString(pos.usagePositional())
	}
	sb.WriteString("\nFLAGS:\n")
	for _, flag := range c.Flags {
		fmt.Fprintf(sb, "  --%-20s", flag.name())
		sb.WriteString("  ")
		sb.WriteString(flag.usageFlag())
		sb.WriteString("\n")
	}
	return sb.String()
}

func pickLast[E any, S ~[]E](x S) E {
	return x[len(x)-1]
}

func Ptr[T any](x T) *T {
	return &x
}
