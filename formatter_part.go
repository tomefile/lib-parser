package libparser

import (
	"fmt"
)

type Scope map[string]string

type FormatPart interface {
	Eval(Scope) (string, error)
}

type LiteralFormat struct {
	Literal string
}

func (format LiteralFormat) Eval(_ Scope) (string, error) {
	return format.Literal, nil
}

type VariableFormat struct {
	Name     string
	Modifier FormatModifier
}

func (format VariableFormat) Eval(scope Scope) (string, error) {
	value, exists := scope[format.Name]
	if !exists {
		return format.Name, fmt.Errorf("variable %q is not defined", value)
	}

	if format.Modifier == nil {
		return value, nil
	}

	return format.Modifier(value), nil
}
