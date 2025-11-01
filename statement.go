package libparser

type StatementKind byte

const (
	SK_NULL StatementKind = iota
	SK_COMMENT
	SK_READ_ERROR
	SK_EOF_ERROR
	SK_SYNTAX_ERROR
	SK_DIRECTIVE
	SK_MACRO
	SK_EXEC

	ErrMin = SK_READ_ERROR
	ErrMax = SK_SYNTAX_ERROR
)

var NullStatement = Statement{Kind: SK_NULL}

type Statement struct {
	Kind    StatementKind
	Literal string
	// (Optional) statement arguments.
	//
	// WARN: Will be `nil` if not appropriate for the kind
	Args []string

	// The offset in bytes where this statement begins
	OffsetStart uint
	// The offset in bytes where this statement ends
	OffsetEnd uint

	// (Optional) statements inside of a block surrounded by curly-braces.
	//
	// WARN: Will be `nil` if not appropriate for the kind
	Children []Statement
}

func (statement Statement) IsConsumable() bool {
	return statement.Kind != SK_NULL && statement.Kind != SK_EOF_ERROR
}

func (statement Statement) IsError() bool {
	return statement.Kind >= ErrMin && statement.Kind <= ErrMax
}
