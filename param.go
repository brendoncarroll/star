package star

import (
	"fmt"
	"math"

	"go.brendoncarroll.net/exp/slices2"
)

// ParamID uniquely identifies a parameter.
type ParamID string

// Parser converts strings to a value of type T.
type Parser[T any] = func(string) (T, error)

// Parameter is the common interface implemented by all Parameters.
// There are few concrete types
// - Required
// - Optional
// - Repeated
type Parameter interface {
	isParam()

	name() ParamID
	getShortDoc() string
	parse(string) (any, error)
	minCount() int
	maxCount() int
}

// Positional is a parameter that can be used as a positional argument
type Positional interface {
	Parameter

	usagePositional() string
}

// Flag is a parameter which can be specified on the command line
type Flag interface {
	Parameter

	usageFlag() string
}

// Required is a required parameter
// If it is not satisfied, then the command will error
type Required[T any] struct {
	// ID uniquely identifies the parameter.
	ID    ParamID
	Parse Parser[T]

	ShortDoc string
}

func (p Required[T]) Load(c Context) T {
	panicIfNotHas(p.ID, c)
	return c.Values[p.ID][0].(T)
}

func (p Required[T]) name() ParamID {
	return p.ID
}

func (p Required[T]) parse(x string) (any, error) {
	return p.Parse(x)
}

func (p Required[T]) isParam() {}

func (p Required[T]) usagePositional() string {
	return fmt.Sprintf("<%v>", p.ID)
}

func (p Required[T]) usageFlag() string {
	return "(required)"
}

func (p Required[T]) minCount() int {
	return 1
}

func (p Required[T]) maxCount() int {
	return 1
}

func (p Required[T]) getShortDoc() string {
	return p.ShortDoc
}

var _ Parameter = Optional[struct{}]{}

// Optional is an optional parameter, it can be provided once, or not at all.
type Optional[T any] struct {
	ID    ParamID
	Parse Parser[T]

	// ShortDoc is a short description of the parameter, used in the help text.
	// It should be less than a single line of text.
	ShortDoc string
}

// Load loads the value for an optional parameter
func (opt Optional[T]) LoadOpt(c Context) (T, bool) {
	panicIfNotHas(opt.ID, c)
	vals := c.Values[opt.ID]
	if len(vals) == 0 {
		var zero T
		return zero, false
	}
	return vals[0].(T), true
}

func (opt Optional[T]) name() ParamID {
	return opt.ID
}

func (p Optional[T]) parse(x string) (any, error) {
	return p.Parse(x)
}

func (p Optional[T]) getShortDoc() string {
	return p.ShortDoc
}

func (opt Optional[T]) usagePositional() string {
	return fmt.Sprintf("[%v]", opt.ID)
}

func (opt Optional[T]) usageFlag() string {
	return "(required)"
}

func (opt Optional[T]) minCount() int {
	return 0
}

func (opt Optional[T]) maxCount() int {
	return 1
}

func (opt Optional[T]) isParam() {}

// Repeated is a parameter that can be passed as a flag multiple times.
type Repeated[T any] struct {
	ID    ParamID
	Parse Parser[T]
	Min   int

	ShortDoc string
}

func (r Repeated[T]) Load(c Context) []T {
	panicIfNotHas(r.ID, c)
	vals := c.Values[r.ID]
	return slices2.Map(vals, func(x any) T {
		return x.(T)
	})
}

func (p Repeated[T]) name() ParamID {
	return p.ID
}

func (p Repeated[T]) parse(x string) (any, error) {
	return p.Parse(x)
}

func (p Repeated[T]) getShortDoc() string {
	return p.ShortDoc
}

func (r Repeated[T]) usagePositional() string {
	return fmt.Sprintf("[%s ...]", r.ID)
}

func (r Repeated[T]) usageFlag() string {
	return "(repeated)"
}

func (r Repeated[T]) minCount() int {
	return 0
}

func (r Repeated[T]) maxCount() int {
	return math.MaxInt
}

func (r Repeated[T]) isParam() {}

// Boolean is a Parameter that either exists or doesn't
type Boolean struct {
	Name ParamID
}

func panicIfNotHas(name ParamID, c Context) {
	if !c.self.HasParam(name) {
		panic(fmt.Sprintf("Command does not take param %q", name))
	}
}
