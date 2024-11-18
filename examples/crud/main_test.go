package main

import (
	"testing"

	"go.brendoncarroll.net/star/teststar"
)

func TestCommand(t *testing.T) {
	teststar.OutContainsString(t, &rootCmd, []string{"create"}, "CREATE")
	teststar.OutContainsString(t, &rootCmd, []string{"delete", "abc123"}, "DELETE")
}
