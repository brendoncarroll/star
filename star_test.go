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
		Flags []AnyParam

		Values map[Symbol][]any
		Extra  []string
	}
	tcs := []testCase{
		{
			Args: []string{"--set-int", "117", "extra", "stuff"},
			Flags: []AnyParam{
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
			Flags: []AnyParam{
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

func TestParsePos(t *testing.T) {
	type testCase struct {
		Args []string
		Pos  []AnyParam

		Values map[Symbol][]any
		Extra  []string
	}
	tcs := []testCase{
		{
			Args: []string{"1", "a", "b", "c"},
			Pos: []AnyParam{
				Param[string]{
					Name:     "must-have",
					Repeated: false,
					Parse:    ParseString,
				},
				Param[string]{
					Name:     "xs",
					Repeated: true,
					Parse:    ParseString,
				},
			},

			Extra: nil,
			Values: map[Symbol][]any{
				"must-have": {"1"},
				"xs":        {"a", "b", "c"},
			},
		},
		{
			Args: []string{"1"},
			Pos: []AnyParam{
				Param[string]{
					Name:     "must-have",
					Repeated: false,
					Parse:    ParseString,
				},
				Param[string]{
					Name:     "xs",
					Repeated: true,
					Parse:    ParseString,
				},
			},

			Extra: nil,
			Values: map[Symbol][]any{
				"must-have": {"1"},
			},
		},
		{
			Args: []string{},
			Pos: []AnyParam{
				Param[string]{
					Name:     "has-default",
					Default:  Ptr("default-value"),
					Repeated: false,
					Parse:    ParseString,
				},
			},
			Extra: []string{},
			Values: map[Symbol][]any{
				"has-default": {"default-value"},
			},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[Symbol][]any)
			extra, err := ParsePos(dst, tc.Pos, tc.Args)
			require.NoError(t, err)
			assert.Equal(t, tc.Values, dst)
			assert.Equal(t, tc.Extra, extra)
		})
	}
}
