package libparser

import (
	"fmt"
	"strings"
)

type NodeString struct {
	Segments SegmentedString
	NodeContext
}

func (node *NodeString) Context() NodeContext {
	return node.NodeContext
}

func (node *NodeString) String() string {
	value := node.Segments.String()
	if ShouldStringBeQuoted(value) {
		return fmt.Sprintf("\"%s\"", node.Segments.String())
	}
	return value
}

func (node *NodeString) Eval(locals Locals) (string, error) {
	var builder strings.Builder

	for _, segment := range node.Segments {
		part, err := segment.Eval(locals)
		if err != nil {
			return "", err
		}
		builder.WriteString(part)
	}

	return builder.String(), nil
}

func NewSimpleNodeString(contents string) *NodeString {
	return &NodeString{
		Segments: SegmentedString{
			&LiteralStringSegment{
				Contents: contents,
			},
		},
	}
}

// ————————————————————————————————

type StringSegment interface {
	Segment() string
	Eval(Locals) (string, error)
}

// ————————————————————————————————

type SegmentedString []StringSegment

func (segments SegmentedString) String() string {
	var builder strings.Builder
	for _, segment := range segments {
		builder.WriteString(segment.Segment())
	}
	return builder.String()
}

// ————————————————————————————————

type LiteralStringSegment struct {
	Contents string
}

func (segment *LiteralStringSegment) Segment() string {
	return segment.Contents
}

func (segment *LiteralStringSegment) Eval(_ Locals) (string, error) {
	return segment.Contents, nil
}

// ————————————————————————————————

type VariableStringSegment struct {
	Name       string
	Modifiers  []StringModifier
	IsOptional bool
}

func (segment *VariableStringSegment) Segment() string {
	if len(segment.Modifiers) == 0 {
		if segment.IsOptional {
			return fmt.Sprintf("${%s?}", segment.Name)
		}
		return fmt.Sprintf("$%s", segment.Name)
	}

	name := segment.Name
	if segment.IsOptional {
		name += "?"
	}

	var builder strings.Builder
	for _, modifier := range segment.Modifiers {
		builder.WriteString(":" + modifier.String())
	}

	return fmt.Sprintf("${%s%s}", segment.Name, builder.String())
}

func (segment *VariableStringSegment) Eval(locals Locals) (string, error) {
	value, exists := locals[segment.Name]
	if !exists {
		if segment.IsOptional {
			return "", nil
		}
		return segment.Name, fmt.Errorf(
			"variable %q is not defined in the current scope",
			segment.Name,
		)
	}

	for _, modifier := range segment.Modifiers {
		value = modifier.Call(locals, value)
	}

	return value, nil
}

// ————————————————————————————————

func ShouldStringBeQuoted(value string) bool {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		return false
	}
	return strings.ContainsAny(value, " \"'`[]{}")
}
