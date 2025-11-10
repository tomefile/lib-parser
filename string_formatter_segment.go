package libparser

import "fmt"

type Locals map[string]string

type Segment interface {
	Eval(Locals) (string, error)
}

type VariableSegment struct {
	Name       string
	Modifier   StringModifier
	IsOptional bool
}

func (segment *VariableSegment) Eval(locals Locals) (string, error) {
	value, exists := locals[segment.Name]
	if !exists {
		if segment.IsOptional {
			return "", nil
		}
		return segment.Name, fmt.Errorf("variable %q is not defined in the current scope", value)
	}

	if segment.Modifier == nil {
		return value, nil
	}

	return segment.Modifier(value), nil
}
