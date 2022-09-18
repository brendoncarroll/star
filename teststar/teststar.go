package teststar

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/brendoncarroll/star"

	"github.com/stretchr/testify/require"
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
		t.Fatalf("output: %v does not contain: %v", stdout, expect)
	}
}

// run runs c with cmdStr
func run(t testing.TB, c *star.Command, args []string) (stdout []byte, stderr []byte) {
	var inbuf, outbuf, errbuf bytes.Buffer
	err := star.Execute(context.Background(), c, star.Env{In: &inbuf, Out: &outbuf, Err: &errbuf}, "TEST", args)
	require.NoError(t, err)
	return outbuf.Bytes(), errbuf.Bytes()
}
