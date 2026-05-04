package star

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunPanicsWithoutPosName(t *testing.T) {
	cmd := Command{
		Pos: []Positional{&Required[string]{Parse: ParseString}},
		F: func(c Context) error {
			return nil
		},
	}

	require.PanicsWithValue(t, "positional parameter at index 0 must set non-empty PosName", func() {
		_ = Run(context.Background(), cmd, map[string]string{}, "test", nil, nil, nil, nil)
	})
}

func TestParseFlags(t *testing.T) {
	type testCase struct {
		Args  []string
		Flags map[string]Flag

		Values map[Parameter][]any
		Extra  []string
	}
	setInt := &Required[int]{Parse: strconv.Atoi}
	setInts := &Repeated[int]{Parse: strconv.Atoi}
	aa := &Required[int]{Parse: strconv.Atoi}
	bb := &Optional[int]{Parse: strconv.Atoi}
	cc := &Required[int]{Parse: strconv.Atoi}
	tcs := []testCase{
		{
			Args: []string{"--set-int", "117", "extra", "stuff"},
			Flags: map[string]Flag{
				"set-int": setInt,
			},
			Values: map[Parameter][]any{
				setInt: {117},
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
				"set-ints": setInts,
			},
			Values: map[Parameter][]any{
				setInts: {3, 6, 9},
			},
			Extra: []string{"damn", "she", "fine"},
		},
		{
			Args: []string{
				"--aa", "1",
				"--cc", "3",
			},
			Flags: map[string]Flag{
				"aa": aa,
				"bb": bb,
				"cc": cc,
			},
			Values: map[Parameter][]any{
				aa: []any{1},
				cc: []any{3},
			},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[Parameter][]any)
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

		Values map[Parameter][]any
		Extra  []string
	}
	mustHave := &Required[string]{PosName: "must-have", Parse: ParseString}
	xs := &Repeated[string]{PosName: "xs", Parse: ParseString}
	mustHave2 := &Required[string]{PosName: "must-have", Parse: ParseString}
	xs2 := &Repeated[string]{PosName: "xs", Parse: ParseString}
	optional := &Optional[string]{PosName: "optional", Parse: ParseString}
	tcs := []testCase{
		{
			Args: []string{"1", "a", "b", "c"},
			Pos: []Positional{
				mustHave,
				xs,
			},

			Extra: nil,
			Values: map[Parameter][]any{
				mustHave: {"1"},
				xs:       {"a", "b", "c"},
			},
		},
		{
			Args: []string{"1"},
			Pos: []Positional{
				mustHave2,
				xs2,
			},

			Extra: nil,
			Values: map[Parameter][]any{
				mustHave2: {"1"},
			},
		},
		{
			Args: []string{},
			Pos: []Positional{
				optional,
			},
			Extra:  []string{},
			Values: map[Parameter][]any{},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			dst := make(map[Parameter][]any)
			extra, err := ParsePos(dst, tc.Pos, tc.Args)
			require.NoError(t, err)
			assert.Equal(t, tc.Values, dst)
			assert.Equal(t, tc.Extra, extra)
		})
	}
}
