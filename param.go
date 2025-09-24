package star

import (
	"fmt"
	"math"

	"go.brendoncarroll.net/exp/slices2"
)

// Symbol is the name of a parameter
type Symbol string

type Parser[T any] = func(string) (T, error)

// Parameter is the common interface implemented by all Parameters.
// There are few concrete types
// - Required
// - Optional
// - Repeated
type Parameter interface {
	isParam()

	name() Symbol
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
	Name  Symbol
	Parse Parser[T]
}

func (p Required[T]) Load(c Context) T {
	panicIfNotHas(p.Name, c)
	return c.Params[p.Name][0].(T)
}

func (p Required[T]) name() Symbol {
	return p.Name
}

func (p Required[T]) parse(x string) (any, error) {
	return p.Parse(x)
}

func (p Required[T]) isParam() {}

func (p Required[T]) isSatisfied(m map[Symbol][]string) bool {
	_, yes := m[p.Name]
	return yes
}

func (p Required[T]) usagePositional() string {
	return fmt.Sprintf("<%v>", p.Name)
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

var _ Parameter = Optional[struct{}]{}

// Optional is an optional parameter, it can be provided once, or not at all.
type Optional[T any] struct {
	Name  Symbol
	Parse Parser[T]
}

// Load loads the value for an optional parameter
func (opt Optional[T]) LoadOpt(c Context) (T, bool) {
	panicIfNotHas(opt.Name, c)
	vals := c.Params[opt.Name]
	if len(vals) == 0 {
		var zero T
		return zero, false
	}
	return vals[0].(T), true
}

func (opt Optional[T]) name() Symbol {
	return opt.Name
}

func (p Optional[T]) parse(x string) (any, error) {
	return p.Parse(x)
}

func (opt Optional[T]) usagePositional() string {
	return fmt.Sprintf("[%v]", opt.Name)
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
	Name  Symbol
	Parse Parser[T]
	Min   int
}

func (r Repeated[T]) Load(c Context) []T {
	panicIfNotHas(r.Name, c)
	vals := c.Params[r.Name]
	return slices2.Map(vals, func(x any) T {
		return x.(T)
	})
}

func (p Repeated[T]) name() Symbol {
	return p.Name
}

func (p Repeated[T]) parse(x string) (any, error) {
	return p.Parse(x)
}

func (r Repeated[T]) usagePositional() string {
	return fmt.Sprintf("[%s ...]", r.Name)
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
	Name Symbol
}

func panicIfNotHas(name Symbol, c Context) {
	if !c.self.HasParam(name) {
		panic(fmt.Sprintf("Command does not take param %q", name))
	}
}
