package teststar

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.brendoncarroll.net/star"
)

// OutIs runs the command with in as an input string, and checks that the output is out.
func OutIsString(t testing.TB, c *star.Command, args []string, expect string) {
	t.Helper()
	_, stdout := run(t, c, args)
	require.Equal(t, expect, string(stdout))
}

func OutContainsString(t testing.TB, c *star.Command, args []string, expect string) {
	t.Helper()
	stdout, _ := run(t, c, args)
	if !strings.Contains(string(stdout), expect) {
		t.Log(string(stdout))
		t.Fatalf("output: %v does not contain: %v", stdout, expect)
	}
}

// run runs c with cmdStr
func run(t testing.TB, c *star.Command, args []string) (stdout []byte, stderr []byte) {
	ctx := context.Background()
	var inbuf, outbuf, errbuf bytes.Buffer
	env := map[string]string{}
	err := star.Run(ctx, *c, env, t.Name(), args, &inbuf, &outbuf, &errbuf)
	require.NoError(t, err)
	return outbuf.Bytes(), errbuf.Bytes()
}
