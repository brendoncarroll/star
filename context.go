package star

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// Context is the context in which a command is run.
// It provides parsed parameters, input and output streams, and a go context.Context
type Context struct {
	context.Context
	// Values are parsed Values filling a Parameter by Name
	Values   map[ParamID][]any
	Env      map[string]string
	StdIn    io.Reader
	StdOut   io.Writer
	StdErr   io.Writer
	CalledAs string
	Extra    []string

	self *Command
}

// Printf is a convenience function for writing to stdout.
func (c Context) Printf(format string, a ...any) {
	_, err := fmt.Fprintf(c.StdOut, format, a...)
	if err != nil {
		panic(err)
	}
}

func Run(ctx context.Context, cmd Command, env map[string]string, calledAs string, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	params := make(map[ParamID][]any)
	args, err := ParseFlags(params, cmd.Flags, args)
	if err != nil {
		fmt.Fprint(stderr, cmd.Doc(calledAs))
		return err
	}
	args, err = ParsePos(params, cmd.Pos, args)
	if err != nil {
		fmt.Fprint(stderr, cmd.Doc(calledAs))
		return err
	}
	if err := checkParams(params, cmd.Flags, cmd.Pos); err != nil {
		fmt.Fprint(stderr, cmd.Doc(calledAs))
		return err
	}
	return cmd.F(Context{
		Context:  ctx,
		Env:      env,
		StdOut:   stdout,
		StdIn:    stdin,
		StdErr:   stderr,
		Values:   params,
		Extra:    args,
		CalledAs: calledAs,

		self: &cmd,
	})
}

func checkParams(valueMap map[ParamID][]any, flags map[string]Flag, pos []Positional) error {
	var allParams []Parameter
	for _, x := range flags {
		allParams = append(allParams, x)
	}
	for _, x := range pos {
		allParams = append(allParams, x)
	}
	for _, param := range allParams {
		vals := valueMap[param.name()]
		if len(vals) < param.minCount() {
			return fmt.Errorf("missing value for parameter %q", param.name())
		}
		if len(vals) > param.maxCount() {
			return fmt.Errorf("multiple values provided for parameter %q", param.name())
		}
	}
	return nil
}

// ParsePos parses positional arguments
func ParsePos(dst map[ParamID][]any, params []Positional, args []string) (rest []string, err error) {
	for _, param := range params {
		for i := 0; i < param.maxCount() && len(args) > 0; i++ {
			val, rest, err := parseOnePos(param, args)
			if err != nil {
				return nil, err
			}
			dst[param.name()] = append(dst[param.name()], val)
			args = rest
		}
	}
	return args, nil
}

func parseOnePos(p Parameter, args []string) (vals any, rest []string, err error) {
	for i := 0; i < len(args); i++ {
		if isFlag(args[i]) {
			// ignore flags
			rest = append(rest, args[i])
			i += 1
			continue
		}
		val, err := p.parse(args[i])
		if err != nil {
			return nil, nil, err
		}
		return val, append(rest, args[i+1:]...), nil
	}
	return nil, args, fmt.Errorf("no args left to parse for positional argument %q", p.name())
}

const (
	flagPrefix      = "--"
	shortFlagPrefix = "-"
)

func isFlag(x string) bool {
	return strings.HasPrefix(x, flagPrefix)
}

func isShortFlag(x string) bool {
	return strings.HasPrefix(x, "-") && len(x) == 2
}

// ParseFlags takes a slice of args, and parses paramaeters in the list of flags.
// ParseFlags writes values to dst.
func ParseFlags(dst map[ParamID][]any, flags map[string]Flag, args []string) (rest []string, err error) {
	flagIndex := make(map[string]Flag)
	for k, flag := range flags {
		flagIndex[k] = flag
	}

	for len(args) > 0 {
		arg := args[0]
		if k, yes := strings.CutPrefix(arg, flagPrefix); yes {
			if param, exists := flagIndex[k]; exists {
				// TODO handle equals
				if len(args) < 2 {
					return nil, fmt.Errorf("arg named but not provided for %q", k)
				}
				v, err := param.parse(args[1])
				if err != nil {
					return nil, err
				}
				name := param.name()
				dst[name] = append(dst[name], v)
				args = args[2:]
				continue
			}
		}
		args = args[1:]
		rest = append(rest, arg)
	}

	for _, param := range flagIndex {
		name := param.name()
		if len(dst[name]) < param.minCount() {
			return nil, fmt.Errorf("missing flag %q", param.name())
		}
	}
	return rest, nil
}

// ParseString is a parser for strings, it is the identity function on strings, and never errors.
func ParseString(x string) (string, error) {
	return x, nil
}
