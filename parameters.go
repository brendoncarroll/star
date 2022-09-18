package star

import (
	"fmt"
	"reflect"
	"strings"
)

type Type = reflect.Type

func TypeOf[T any]() Type {
	var zero T
	return reflect.TypeOf(zero)
}

// Param is a value that needs to be passed to a program
type Param struct {
	Type       Type
	IsRequired bool
	Default    any
	Check      func(any) error
}

func (p Param) String() string {
	sb := &strings.Builder{}
	sb.WriteString("{")
	sb.WriteString(p.Type.String())
	fmt.Fprintf(sb, " required=%v", p.IsRequired)
	if p.Default != nil {
		fmt.Fprintf(sb, " default=%v", p.Default)
	}
	sb.WriteString("}")
	return sb.String()
}

// Pos is a specification for a positional argument
type Pos struct {
	Name  string
	Short string

	IsRepeated bool
	Param
}

func (p Pos) WithDefault(x any) Pos {
	p.Default = x
	return p
}

// NewPos returns a new positional parameter
func NewPos[T any](name, shortDoc string, isRequired bool) Pos {
	ty := TypeOf[T]()
	isRepeated := ty.Kind() == reflect.Slice
	return Pos{
		Name:       name,
		Short:      shortDoc,
		IsRepeated: isRepeated,
		Param: Param{
			Type:       ty,
			IsRequired: isRequired,
		},
	}
}

// Flag is a the specification of an argument
type Flag struct {
	Name  string
	Char  byte
	Short string

	IsRepeated bool
	Param
}

func (f Flag) WithDefault(x any) Flag {
	f.Default = x
	return f
}

// NewFlag constructs a new Flag parameter
func NewFlag[T any](name, shortDoc string, isRequired bool) Flag {
	ty := TypeOf[T]()
	isRepeated := ty.Kind() == reflect.Slice
	return Flag{
		Name:       name,
		Short:      shortDoc,
		IsRepeated: isRepeated,
		Param: Param{
			Type:       ty,
			IsRequired: isRequired,
		},
	}
}
