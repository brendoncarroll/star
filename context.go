package star

import (
	"bufio"
	"context"
	"fmt"
	"strings"
)

// Context is the context in which a command is run.
// It provides parsed parameters, input and output streams, and a go context.Context
type Context struct {
	context.Context
	Params   map[Symbol][]any
	Env      map[string]string
	StdIn    *bufio.Reader
	StdOut   *bufio.Writer
	StdErr   *bufio.Writer
	CalledAs string
	Extra    []string

	self *Command
}

func (c Context) Printf(format string, a ...any) {
	_, err := fmt.Fprintf(c.StdOut, format, a...)
	if err != nil {
		panic(err)
	}
}

func (c Context) Param(x Symbol) (any, bool) {
	if !c.self.HasParam(x) {
		panic(fmt.Sprintf("Command does not take param %q", x))
	}
	v, exists := c.Params[x]
	return v, exists
}

func Run(ctx context.Context, cmd Command, env map[string]string, calledAs string, args []string, stdin *bufio.Reader, stdout *bufio.Writer, stderr *bufio.Writer) error {
	params := make(map[Symbol][]any)
	args, err := ParseFlags(params, cmd.Flags, args)
	if err != nil {
		stderr.WriteString(cmd.Doc(calledAs))
		return err
	}
	args, err = ParsePos(params, cmd.Pos, args)
	if err != nil {
		stderr.WriteString(cmd.Doc(calledAs))
		return err
	}
	if err := checkParams(params, cmd.Flags, cmd.Pos); err != nil {
		stderr.WriteString(cmd.Doc(calledAs))
		return err
	}
	defer stdout.Flush()
	defer stderr.Flush()
	return cmd.F(Context{
		Context:  ctx,
		StdOut:   stdout,
		StdIn:    stdin,
		StdErr:   stderr,
		Params:   params,
		Extra:    args,
		CalledAs: calledAs,

		self: &cmd,
	})
}

func checkParams(valueMap map[Symbol][]any, flags, pos []AnyParam) error {
	for _, params := range [][]AnyParam{flags, pos} {
		for _, param := range params {
			vals := valueMap[param.name()]
			if len(vals) < 1 && !param.isRepeated() {
				return fmt.Errorf("missing value for parameter %q", param.name())
			}
			if len(vals) > 1 && !param.isRepeated() {
				return fmt.Errorf("multiple values provided for parameter %q", param.name())
			}
		}
	}
	return nil
}

// ParsePos parses positional arguments
func ParsePos(dst map[Symbol][]any, params []AnyParam, args []string) (rest []string, err error) {
	for _, param := range params {
		switch {
		case param.isRepeated():
			// if the param is repeated, consume continuously, only if args are available
			for len(args) > 0 {
				val, rest, err := parseOnePos(param, args)
				if err != nil {
					return nil, err
				}
				dst[param.name()] = append(dst[param.name()], val)
				args = rest
			}
		case param.hasDefault():
			if len(args) > 0 {
				val, rest, err := parseOnePos(param, args)
				if err != nil {
					return nil, err
				}
				dst[param.name()] = append(dst[param.name()], val)
				args = rest
			} else {
				val, err := param.makeDefault()
				if err != nil {
					return nil, err
				}
				dst[param.name()] = append(dst[param.name()], val)
			}
		default:
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

func parseOnePos(p AnyParam, args []string) (vals any, rest []string, err error) {
	for i := 0; i < len(args); i++ {
		if isFlag(args[i]) {
			// ignore flags
			rest = append(rest, args[i])
			i += 2
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

const flagPrefix = "--"

func isFlag(x string) bool {
	return strings.HasPrefix(x, flagPrefix)
}

// ParseFlags takes a slice of args, and parses paramaeters in the list of flags.
// ParseFlags writes values to dst.
func ParseFlags(dst map[Symbol][]any, flags []AnyParam, args []string) (rest []string, err error) {
	flagIndex := make(map[Symbol]AnyParam)
	for _, flag := range flags {
		flagIndex[flag.name()] = flag
	}

	for len(args) > 0 {
		arg := args[0]
		if k, yes := strings.CutPrefix(arg, flagPrefix); yes {
			sym := Symbol(k)
			param := flagIndex[sym]
			if _, exists := flagIndex[sym]; exists {
				// TODO handle equals
				if len(args) < 2 {
					return nil, fmt.Errorf("arg named but not provided for %q", k)
				}
				v, err := param.parse(args[1])
				if err != nil {
					return nil, err
				}
				dst[sym] = append(dst[sym], v)
				args = args[2:]
				continue
			}
		}
		args = args[1:]
		rest = append(rest, arg)
	}

	for name, param := range flagIndex {
		if len(dst[name]) == 0 && !param.isRepeated() {
			if param.hasDefault() {
				val, err := param.makeDefault()
				if err != nil {
					return rest, err
				}
				dst[name] = append(dst[name], val)
				return rest, nil
			}
			return nil, fmt.Errorf("missing flag %q", param.name())
		}
	}
	return rest, nil
}

// ParseString is a parser for strings, it is the identity function on strings, and never errors.
func ParseString(x string) (string, error) {
	return x, nil
}
