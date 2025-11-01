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
		Flags map[string]Flag

		Values map[ParamID][]any
		Extra  []string
	}
	tcs := []testCase{
		{
			Args: []string{"--set-int", "117", "extra", "stuff"},
			Flags: map[string]Flag{
				"set-int": Required[int]{
					ID:    "set-int",
					Parse: strconv.Atoi,
				},
			},
			Values: map[ParamID][]any{
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
			Flags: map[string]Flag{
				"set-ints": Repeated[int]{
					ID:    "set-ints",
					Parse: strconv.Atoi,
				},
			},
			Values: map[ParamID][]any{
				"set-ints": {3, 6, 9},
			},
			Extra: []string{"damn", "she", "fine"},
		},
		{
			Args: []string{
				"--aa", "1",
				"--cc", "3",
			},
			Flags: map[string]Flag{
				"aa": Required[int]{ID: "aa", Parse: strconv.Atoi},
				"bb": Optional[int]{
					ID:    "bb",
					Parse: strconv.Atoi,
				},
				"cc": Required[int]{ID: "cc", Parse: strconv.Atoi},
			},
			Values: map[ParamID][]any{
				"aa": []any{1},
				"cc": []any{3},
			},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[ParamID][]any)
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
		Pos  []Positional

		Values map[ParamID][]any
		Extra  []string
	}
	tcs := []testCase{
		{
			Args: []string{"1", "a", "b", "c"},
			Pos: []Positional{
				Required[string]{
					ID:    "must-have",
					Parse: ParseString,
				},
				Repeated[string]{
					ID:    "xs",
					Parse: ParseString,
				},
			},

			Extra: nil,
			Values: map[ParamID][]any{
				"must-have": {"1"},
				"xs":        {"a", "b", "c"},
			},
		},
		{
			Args: []string{"1"},
			Pos: []Positional{
				Required[string]{
					ID:    "must-have",
					Parse: ParseString,
				},
				Repeated[string]{
					ID:    "xs",
					Parse: ParseString,
				},
			},

			Extra: nil,
			Values: map[ParamID][]any{
				"must-have": {"1"},
			},
		},
		{
			Args: []string{},
			Pos: []Positional{
				Optional[string]{
					ID:    "optional",
					Parse: ParseString,
				},
			},
			Extra:  []string{},
			Values: map[ParamID][]any{},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[ParamID][]any)
			extra, err := ParsePos(dst, tc.Pos, tc.Args)
			require.NoError(t, err)
			assert.Equal(t, tc.Values, dst)
			assert.Equal(t, tc.Extra, extra)
		})
	}
}
