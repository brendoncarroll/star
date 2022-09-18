package star

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// Func is the type of functions which can be run as commands
type Func = func(ctx *Context) error

// Command is a Func and associated metadata.
type Command struct {
	// Short is a short doc string
	Short string

	// Pos is a slice of specifications for positional arguments
	Pos []Pos
	// Flags are the flags the command is expecting
	Flags []Flag

	// F is the function to run
	F Func
}

func (c *Command) Clone() *Command {
	c2 := *c
	c2.Flags = slices.Clone(c2.Flags)
	c2.Pos = slices.Clone(c2.Pos)
	return &c2
}

func (c *Command) Usage() string {
	return usage(c.Pos)
}

// usage
func usage(poss []Pos) string {
	sb := strings.Builder{}
	for i, a := range poss {
		if i > 0 {
			sb.WriteString(" ")
		}
		if !a.IsRequired {
			sb.WriteString("[")
		}
		sb.WriteString(a.Name)
		if a.IsRepeated {
			sb.WriteString(" ...")
		}
		if !a.IsRequired {
			sb.WriteString("]")
		}
	}
	return sb.String()
}

func describeUsage(calledAs string, poss []Pos) string {
	sb := strings.Builder{}
	sb.WriteString("Usage:\n  ")
	sb.WriteString(calledAs)
	sb.WriteString("  ")
	sb.WriteString(usage(poss))
	sb.WriteString("\n\n")
	return sb.String()
}

func describeFlags(flags []Flag) string {
	sb := strings.Builder{}
	sb.WriteString("Flags:\n")
	for _, f := range flags {
		sb.WriteString("\t")
		sb.WriteString(f.Name)
		sb.WriteString(" :")
		sb.WriteString(fmt.Sprint(f.Type))
	}
	return sb.String()
}

func describeChildren(children map[string]*Command) string {
	sb := &strings.Builder{}
	sb.WriteString("Available Commands:\n")
	for name, child := range children {
		fmt.Fprintf(sb, "  %-10s %s\n", name, child.Short)
	}
	return sb.String()
}

func parentLong(calledAs string, short string, pos []Pos, children map[string]*Command) string {
	sb := &strings.Builder{}
	sb.WriteString(short)
	sb.WriteString("\n\n")

	sb.WriteString(calledAs)
	sb.WriteString(" ")
	sb.WriteString(usage(pos))
	sb.WriteString("\n\n")

	sb.WriteString(describeChildren(children))

	return sb.String()
}

// NewParent constructs a parent command grouping several child commans
func NewParent(short string, flags []Flag, children map[string]*Command) *Command {
	pos := []Pos{NewPos[string]("command", "", false)}
	return &Command{
		Short: short,
		Flags: flags,
		Pos:   pos,
		F: func(ctx *Context) error {
			name, ok := ctx.MaybeString("command")
			if !ok {
				_, err := fmt.Fprintln(ctx.Out, parentLong(ctx.CalledAs, short, pos, children))
				return err
			}
			childCmd, exists := children[name]
			if !exists {
				return fmt.Errorf("%s has no sub-command %s", ctx.CalledAs, name)
			}
			return Execute(ctx.Context, childCmd, ctx.env(), name, ctx.Args)
		},
	}
}

// func WrapHelp(cmd *Command) *Command {
// 	return cmd
// }

// PreExec returns a function x which calls preF before x.
func PreExec(x Func, preF func(ctx *Context) (*Context, error)) Func {
	return func(ctx *Context) error {
		var err error
		ctx, err = preF(ctx)
		if err != nil {
			return err
		}
		return x(ctx)
	}
}

// PostExec returns a function which calls postF after x.
// PostExec always runs after x, even if x returns an error.
// x's error is prioritized in the return value
func PostExec(x Func, postF func(ctx *Context) error) Func {
	return func(ctx *Context) (retErr error) {
		defer func() {
			err := postF(ctx)
			if retErr == nil {
				retErr = err
			}
		}()
		return x(ctx)
	}
}

func PrePostExec(x Func, preF func(ctx *Context) (*Context, error), postF func(ctx *Context) error) Func {
	return PostExec(PreExec(x, preF), postF)
}
