package star

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.brendoncarroll.net/stdctx/logctx"
	"go.uber.org/zap"
)

// Main is a default entrypoint for a Command.
// e.g.
//
//	func main() {
//	  star.Main(yourCommandHere)
//	}
func Main(c Command, opts ...MainOption) {
	calledAs := os.Args[0]
	args := os.Args[1:]
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	// setup the default config
	cfg := mainConfig{
		Background: func() context.Context {
			ctx := context.Background()
			l, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			ctx = logctx.NewContext(ctx, l)
			return ctx
		}(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if err := Run(cfg.Background, c, cfg.Env, calledAs, args, stdin, stdout, stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type mainConfig struct {
	Background context.Context
	Env        map[string]string
}

// MainOption configures the behavior off Main
type MainOption = func(*mainConfig)

// MainBackground returns a MainOption that sets the background context to the provided context.
// Providing MainBackground multiple times will overwrite the previous options.
func MainBackground(bgCtx context.Context) MainOption {
	return func(cfg *mainConfig) {
		cfg.Background = bgCtx
	}
}

// MainIncludeEnv selects environment variables by key-name
// and includes them in the Env passed to commands.
// MainIncludeEnv must be specified to pass through any environment variables.
// Providing MainIncludeEnv multiple times.
func MainIncludeEnv(filter func(string) bool) MainOption {
	return func(cfg *mainConfig) {
		if cfg.Env == nil {
			cfg.Env = make(map[string]string)
		}
		OSEnv(cfg.Env, filter)
	}
}

// OSEnv reads from os.Environ() and copies items to dst if they match filter.
func OSEnv(dst map[string]string, filter func(string) bool) {
	for _, pair := range os.Environ() {
		parts := strings.SplitN(pair, "=", 2)
		if filter(parts[0]) {
			dst[parts[0]] = parts[1]
		}
	}
}

// MatchAll is a predicate for type T that always returns true.
func MatchAll[T any](x T) bool {
	return true
}
