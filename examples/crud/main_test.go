package main

import "testing"
import "github.com/brendoncarroll/star/teststar"

func TestCommand(t *testing.T) {
	teststar.OutContainsString(t, rootCmd, []string{"create"}, "CREATE")
	teststar.OutContainsString(t, rootCmd, []string{"delete", "abc123"}, "DELETE")
}
