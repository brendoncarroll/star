package star

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	type testCase struct {
		Args  []string
		Flags []IParam

		Values map[Symbol][]any
		Extra  []string
	}
	tcs := []testCase{
		{
			Args: []string{"--set-int", "117", "extra", "stuff"},
			Flags: []IParam{
				Param[int]{
					Name:  "set-int",
					Parse: strconv.Atoi,
				},
			},
			Values: map[Symbol][]any{
				"set-int": {117},
			},
			Extra: []string{"extra", "stuff"},
		},
		{
			Args: []string{
				"--set-ints", "3",
				"--set-ints", "6",
				"--set-ints", "9",
				"damn", "she", "fine",
			},
			Flags: []IParam{
				Param[int]{
					Name:     "set-ints",
					Repeated: true,
					Parse:    strconv.Atoi,
				},
			},
			Values: map[Symbol][]any{
				"set-ints": {3, 6, 9},
			},
			Extra: []string{"damn", "she", "fine"},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[Symbol][]any)
			extra, err := ParseFlags(dst, tc.Flags, tc.Args)
			require.NoError(t, err)
			assert.Equal(t, tc.Values, dst)
			assert.Equal(t, tc.Extra, extra)
		})
	}
}
