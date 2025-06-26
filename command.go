package star

import (
	"fmt"
	"strings"

	"go.brendoncarroll.net/exp/slices2"
)

// Symbol is the name of a parameter
type Symbol string

type Param[T any] struct {
	// Name identifies the parameter
	Name Symbol
	// Default is used as input to Parse if the parameter was not specified
	Default *string
	// Repeated means the parameter can be specified multiple times
	Repeated bool
	// Parse is called to convert strings
	Parse func(string) (T, error)
}

func (Param[T]) isParam() {}

// Load acceses a parameter from the context
func (p Param[T]) Load(c Context) T {
	if !c.self.HasParam(p.Name) {
		panic(fmt.Sprintf("Command does not take param %q", p.Name))
	}
	return pickLast(c.Params[p.Name]).(T)
}

// LoadOpt acceses an optional parameter from the context
func (p Param[T]) LoadOpt(c Context) (T, bool) {
	if !c.self.HasParam(p.Name) {
		panic(fmt.Sprintf("Command does not take param %q", p.Name))
	}
	if y, exists := c.Params[p.Name]; exists && len(y) > 0 {
		return pickLast(y).(T), exists
	} else {
		var zero T
		return zero, false
	}
}

// LoadAll returns a slice containing every specified value for the flag
func (p Param[T]) LoadAll(c Context) []T {
	if !c.self.HasParam(p.Name) {
		panic(fmt.Sprintf("Command does not take param %q", p.Name))
	}
	if !p.Repeated {
		panic(fmt.Sprintf("LoadAll on non-repated Param %q", p.Name))
	}
	return slices2.Map(c.Params[p.Name], func(x any) T {
		return x.(T)
	})
}

func (p Param[T]) name() Symbol {
	return p.Name
}

func (p Param[T]) parse(x string) (any, error) {
	return p.Parse(x)
}

func (p Param[T]) hasDefault() bool {
	return p.Default != nil
}

func (p Param[T]) defaultString() string {
	return *p.Default
}

func (p Param[T]) makeDefault() (any, error) {
	return p.parse(*p.Default)
}

func (p Param[T]) parseVar(xs []string) (any, error) {
	var ret []T
	for _, x := range xs {
		y, err := p.Parse(x)
		if err != nil {
			return nil, fmt.Errorf("parsing vararg %v: %w", x, err)
		}
		ret = append(ret, y)
	}
	return ret, nil
}

func (p Param[T]) isRepeated() bool {
	return p.Repeated
}

type AnyParam interface {
	isParam()
	name() Symbol
	parse(string) (any, error)
	parseVar([]string) (any, error)
	hasDefault() bool
	defaultString() string
	makeDefault() (any, error)
	isRepeated() bool
}

func NewString(name Symbol) Param[string] {
	return Param[string]{Name: name, Parse: ParseString}
}

type Metadata struct {
	Short string
	// Tags for grouping by category
	Tags []string
}

type Command struct {
	Flags []AnyParam
	Pos   []AnyParam
	F     func(c Context) error
	Metadata
}

func (c Command) HasParam(x Symbol) bool {
	for i := range c.Pos {
		if c.Pos[i].name() == x {
			return true
		}
	}
	for i := range c.Flags {
		if c.Flags[i].name() == x {
			return true
		}
	}
	return false
}

func (c Command) Doc(calledAs string) string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "%s ", calledAs)
	for _, pos := range c.Pos {
		if pos.hasDefault() {
			fmt.Fprintf(sb, "[%s]", pos.name())
		} else {
			fmt.Fprintf(sb, "<%s>", pos.name())
		}
	}
	sb.WriteString("\nFLAGS:\n")
	for _, flag := range c.Flags {
		fmt.Fprintf(sb, "  --%-20s", flag.name())
		if flag.isRepeated() {
			sb.WriteString("  (repeated)")
		}
		if !flag.isRepeated() && !flag.hasDefault() {
			sb.WriteString("  (required)")
		}
		if flag.hasDefault() {
			fmt.Fprintf(sb, "  (default=%q)", flag.defaultString())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func pickLast[E any, S ~[]E](x S) E {
	return x[len(x)-1]
}

func Ptr[T any](x T) *T {
	return &x
}
