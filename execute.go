package star

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

type mainConfig struct{}

// MainOpt configures Main
type MainOpt = func(*mainConfig)

// Main calls Execute with context.Background and IO streams and args from the OS.
//
// You can call Main from main. e.g.
// func main() {
//   star.Main(cmd)
// }
func Main(name string, c *Command, opts ...MainOpt) {
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	env := Env{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stdout,
	}
	if err := Execute(context.Background(), c, env, name, args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err.Error())
		if ecErr, ok := err.(ExitCodeError); ok {
			os.Exit(ecErr.ExitCode())
		} else {
			os.Exit(1)
		}
	} else {
		os.Exit(0)
	}
}

// Env is the total environment in which the command runs.
// TODO: add FS interface, and environment variables.
type Env struct {
	// In is stdin
	In io.Reader
	// Out is stdout
	Out io.Writer
	Err io.Writer
}

// Execute executes the command c.
func Execute(ctx context.Context, c *Command, env Env, calledAs string, args []string) error {
	vs := map[string]value{}
	args2, err := parseArgs(vs, c.Pos, c.Flags, args)
	if err != nil {
		return err
	}
	if c.F == nil {
		return nil
	}
	return c.F(&Context{
		CalledAs: calledAs,
		Args:     args2,

		In:  env.In,
		Out: env.Out,
		Err: env.Err,

		values: vs,
	})
}

// parseArgs populates m with values parsed from input
func parseArgs(m map[string]value, poss []Pos, flags []Flag, input []string) ([]string, error) {
	params := map[string]Param{}
	for _, flag := range flags {
		// TODO: merge params
		params[flag.Name] = flag.Param
	}
	for _, pos := range poss {
		// TODO: merge params
		params[pos.Name] = pos.Param
	}

	fset := pflag.FlagSet{}
	for _, flag := range flags {
		switch flag.Type {
		case TypeOf[string]():
			fset.String(flag.Name, "", "")
		default:
			return nil, fmt.Errorf("unrecognized type: %v", flag.Type)
		}
	}
	if err := fset.ParseAll(input, func(flag *pflag.Flag, in string) error {
		param, exists := params[flag.Name]
		if !exists {
			return fmt.Errorf("unexpected flag %q", flag.Name)
		}
		v, err := newValue(param.Type, in, param.IsRequired)
		if err != nil {
			return err
		}
		m[flag.Name] = v
		return nil
	}); err != nil {
		return nil, err
	}

	out, err := parsePositional(m, poss, input)
	if err != nil {
		return nil, err
	}

	for name, param := range params {
		if _, exists := m[name]; param.IsRequired && !exists {
			return nil, fmt.Errorf("missing required parameter %s: %v", name, param)
		}
	}
	return out, nil
}

// parsePositional parses positional arguments
func parsePositional(m map[string]value, poss []Pos, input []string) ([]string, error) {
	var optionalCount, minLength, maxLength int
	for i, pos := range poss {
		if i < len(poss)-1 && pos.IsRepeated {
			return nil, errors.New("repeatable positional arguments must be last")
		}
		if pos.IsRepeated {
			maxLength = math.MaxInt
		}
		if pos.IsRequired {
			minLength++
		} else {
			optionalCount++
		}
		if maxLength != math.MaxInt {
			maxLength++
		}
	}
	for _, pos := range poss {
		if len(input) == 0 && pos.IsRequired {
			return nil, fmt.Errorf("missing positional %v", pos.Name)
		} else if len(input) == 0 {
			break
		}
		v, err := newValue(pos.Type, input[0], pos.IsRequired)
		if err != nil {
			return nil, err
		}
		m[pos.Name] = v
		input = input[1:]
	}
	return input, nil
}

// ExitCodeError is an error which can set an exit code.
type ExitCodeError interface {
	ExitCode() int
}

func newValue(ty Type, in string, isReq bool) (value, error) {
	var v any
	switch ty {
	case TypeOf[string]():
		v = in
	case TypeOf[int64]():
		n, err := strconv.ParseInt(in, 10, 64)
		if err != nil {
			return value{}, err
		}
		v = n
	case TypeOf[uint64]():
		n, err := strconv.ParseUint(in, 10, 64)
		if err != nil {
			return value{}, err
		}
		v = n
	case TypeOf[time.Duration]():
		d, err := time.ParseDuration(in)
		if err != nil {
			return value{}, err
		}
		v = d
	default:
		panic(ty)
	}
	return value{
		Type:        ty,
		IsRequired:  isReq,
		WasProvided: true,
		X:           v,
	}, nil
}

func newDefault(ty Type, x interface{}, isReq bool) value {
	return value{
		Type:        ty,
		IsRequired:  isReq,
		WasProvided: false,
		X:           x,
	}
}
