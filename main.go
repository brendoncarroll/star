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
func Main(c Command) {
	calledAs := os.Args[0]
	args := os.Args[1:]
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	env := OSEnv(strings.ToUpper(calledAs) + "_")

	ctx := context.Background()
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	ctx = logctx.NewContext(ctx, l)
	if err := Run(ctx, c, env, calledAs, args, stdin, stdout, stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func OSEnv(prefix string) map[string]string {
	ret := make(map[string]string)
	for _, pair := range os.Environ() {
		parts := strings.SplitN(pair, "=", 2)
		if strings.HasPrefix(parts[0], prefix) {
			ret[parts[0]] = parts[1]
		}
	}
	return ret
}
