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
	Metadata
	Flags map[string]Flag
	Pos   []Positional
	F     func(c Context) error
}

func (c Command) HasParam(x ParamID) bool {
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

	sb.WriteString("\n\nPOSITIONAL:\n")
	if len(c.Pos) == 0 {
		sb.WriteString("  (this command does not accept any positional parameters)\n")
	} else {
		for _, pos := range c.Pos {
			fmt.Fprintf(sb, "  %-10s\t%s\n", pos.name(), pos.getShortDoc())
		}
	}

	sb.WriteString("\nFLAGS:\n")
	if len(c.Flags) == 0 {
		sb.WriteString("  (this command does not accept any parameters as flags)\n")
	} else {
		for key, flag := range c.Flags {
			fmt.Fprintf(sb, "  --%-20s %s\n", key, flag.getShortDoc())
		}
	}
	sb.WriteString("\n")
	return sb.String()
}

func pickLast[E any, S ~[]E](x S) E {
	return x[len(x)-1]
}

func Ptr[T any](x T) *T {
	return &x
}
