package star

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUsage(t *testing.T) {
	c := Command{
		Pos: []Pos{
			NewPos[string]("arg1", "", true),
			NewPos[string]("arg2", "", false),
		},
	}
	require.Equal(t, c.Usage(), "arg1 [arg2]")
}
