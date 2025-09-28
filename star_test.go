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
		Flags []Flag

		Values map[Name][]any
		Extra  []string
	}
	tcs := []testCase{
		{
			Args: []string{"--set-int", "117", "extra", "stuff"},
			Flags: []Flag{
				Required[int]{
					Name:  "set-int",
					Parse: strconv.Atoi,
				},
			},
			Values: map[Name][]any{
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
			Flags: []Flag{
				Repeated[int]{
					Name:  "set-ints",
					Parse: strconv.Atoi,
				},
			},
			Values: map[Name][]any{
				"set-ints": {3, 6, 9},
			},
			Extra: []string{"damn", "she", "fine"},
		},
		{
			Args: []string{
				"--aa", "1",
				"--cc", "3",
			},
			Flags: []Flag{
				Required[int]{Name: "aa", Parse: strconv.Atoi},
				Optional[int]{
					Name:  "bb",
					Parse: strconv.Atoi,
				},
				Required[int]{Name: "cc", Parse: strconv.Atoi},
			},
			Values: map[Name][]any{
				"aa": []any{1},
				"cc": []any{3},
			},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[Name][]any)
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

		Values map[Name][]any
		Extra  []string
	}
	tcs := []testCase{
		{
			Args: []string{"1", "a", "b", "c"},
			Pos: []Positional{
				Required[string]{
					Name:  "must-have",
					Parse: ParseString,
				},
				Repeated[string]{
					Name:  "xs",
					Parse: ParseString,
				},
			},

			Extra: nil,
			Values: map[Name][]any{
				"must-have": {"1"},
				"xs":        {"a", "b", "c"},
			},
		},
		{
			Args: []string{"1"},
			Pos: []Positional{
				Required[string]{
					Name:  "must-have",
					Parse: ParseString,
				},
				Repeated[string]{
					Name:  "xs",
					Parse: ParseString,
				},
			},

			Extra: nil,
			Values: map[Name][]any{
				"must-have": {"1"},
			},
		},
		{
			Args: []string{},
			Pos: []Positional{
				Optional[string]{
					Name:  "optional",
					Parse: ParseString,
				},
			},
			Extra:  []string{},
			Values: map[Name][]any{},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[Name][]any)
			extra, err := ParsePos(dst, tc.Pos, tc.Args)
			require.NoError(t, err)
			assert.Equal(t, tc.Values, dst)
			assert.Equal(t, tc.Extra, extra)
		})
	}
}
