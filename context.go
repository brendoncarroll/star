package star

import (
	"context"
	"fmt"
	"io"
)

// Context is passed to Funcs
type Context struct {
	context.Context

	// In is stdin
	In io.Reader
	// Out is stdout
	Out io.Writer
	// Err is stderr
	Err io.Writer

	CalledAs string
	Args     []string

	values map[string]value
}

func (c *Context) Uint64(x string) uint64 {
	return getValue[uint64](c.values, x)
}

func (c *Context) Int64(x string) int64 {
	return getValue[int64](c.values, x)
}

func (c *Context) String(x string) string {
	return getValue[string](c.values, x)
}

func (c *Context) Time(x string) string {
	return getValue[string](c.values, x)
}

func (c *Context) StringSlice(x string) []string {
	return getValue[[]string](c.values, x)
}

func (c *Context) Bytes(x string) []byte {
	return getValue[[]byte](c.values, x)
}

func (c *Context) MaybeUint64(x string) (uint64, bool) {
	return getMaybeValue[uint64](c.values, x)
}

func (c *Context) MaybeInt64(x string) (int64, bool) {
	return getMaybeValue[int64](c.values, x)
}

func (c *Context) MaybeString(x string) (string, bool) {
	return getMaybeValue[string](c.values, x)
}

func (c *Context) WithValue(k, v interface{}) *Context {
	c2 := *c
	c2.Context = context.WithValue(c.Context, k, v)
	return &c2
}

func (c *Context) env() Env {
	return Env{
		In:  c.In,
		Out: c.Out,
		Err: c.Err,
	}
}

// V is a generic function which retrieves values from a Context
// The V stands for Value
func V[T any](ctx *Context, x string) T {
	return getValue[T](ctx.values, x)
}

// M is a generic function which retrieve maybe values from a Context
// The M stands for Maybe.
func M[T any](ctx *Context, x string) (T, bool) {
	return getMaybeValue[T](ctx.values, x)
}

type value struct {
	Type        Type
	IsRequired  bool
	WasProvided bool
	X           interface{}
}

func getValue[T any](vs map[string]value, x string) T {
	v, exists := vs[x]
	if !exists {
		panic(fmt.Sprintf("asked star.Context for a value which does not exist. This is a bug. You must specify a Flag or Arg for %q", x))
	}
	if err := checkType(v.Type, v.Type); err != nil {
		panic(err.Error())
	}
	return v.X.(T)
}

func getMaybeValue[T any](vs map[string]value, x string) (T, bool) {
	var zero T
	v, exists := vs[x]
	if !exists {
		return zero, false
	}
	if err := checkType(v.Type, v.Type); err != nil {
		panic(err.Error())
	}
	if v.X == nil {
		// no default
		return zero, false
	}
	if v.IsRequired {
		panic(fmt.Sprintf("don't use Maybe* methods for required parameter %q", x))
	}
	return v.X.(T), v.WasProvided
}

func checkType(pty, rty Type) error {
	if !pty.AssignableTo(rty) {
		return errTypeMismatch(pty, rty)
	}
	return nil
}

func errTypeMismatch(pty, rty Type) error {
	return fmt.Errorf("requested data of type %T but parameter is type %v", rty, pty)
}
