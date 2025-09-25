package main

import (
	"testing"

	"go.brendoncarroll.net/star/teststar"
)

func TestCommand(t *testing.T) {
	teststar.OutContainsString(t, &rootCmd, []string{"create"}, "CREATE")
	teststar.OutContainsString(t, &rootCmd, []string{"delete", "abc123"}, "DELETE")

	teststar.OutContainsString(t, &rootCmd, []string{"sub-dir-command", "echo", "ECHO123"}, "ECHO123")

	teststar.OutContainsString(t, &rootCmd, []string{"sub-dir-command", "echo-pos", "foobar1"}, "foobar1")
}
