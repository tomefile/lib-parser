package libparser

import (
	"bufio"

	liberrors "github.com/tomefile/lib-errors"
	"github.com/tomefile/lib-parser/readers"
)

type Parser struct {
	parent          *Parser
	Name            string
	reader          *readers.Reader
	root            *NodeTree
	endOfSectionErr *liberrors.DetailedError // FIXME: There should be a better way
	PostProcessors  []PostProcessor
}

func New(file File) *Parser {
	return &Parser{
		parent: nil,
		Name:   file.Name(),
		reader: readers.New(bufio.NewReader(file)),
		root: &NodeTree{
			Tomes:        map[string]Node{},
			NodeChildren: NodeChildren{},
		},
		endOfSectionErr: nil,
		PostProcessors:  []PostProcessor{},
	}
}

// Appends the [PostProcessor] to be applied to every single node before it gets appended to the tree.
//
// NOTE: Order matters (sequentially from first to last)
func (parser *Parser) With(processor PostProcessor) *Parser {
	parser.PostProcessors = append(parser.PostProcessors, processor)
	return parser
}

// Used for error tracing
func (parser *Parser) SetParent(parent *Parser) *Parser {
	parser.parent = parent
	return parser
}

func (parser *Parser) Result() *NodeTree {
	return parser.root
}

func (parser *Parser) Parse() *liberrors.DetailedError {
	parser.endOfSectionErr = parser.failSyntax("unexpected '}' with no matching '{' pair")

	for {
		err := parser.next(&parser.root.NodeChildren)
		if err != nil {
			if err == EOF {
				break
			}
			return err
		}
	}

	return nil
}
