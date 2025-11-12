package libparser

import (
	"bufio"
	"os"
)

type CanClose interface {
	Close() error
}

var OpenedSources []CanClose

func OpenNew(path string) (*Parser, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	OpenedSources = append(OpenedSources, file)

	return New(path, bufio.NewReader(file)), nil
}

// Closes all files opened using [OpenNew()].
//
// Recommended to use `defer parser.Close()` in your main function.
func (parser *Parser) Close() {
	for _, file := range OpenedSources {
		file.Close()
	}
	OpenedSources = []CanClose{}
}
