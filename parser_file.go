package libparser

import (
	"io"
	"os"
)

var OpenedFiles = []File{}

// The parts of *os.File that parser cares about
type File interface {
	io.Reader
	Name() string
	Close() error
}

// Recommended way of opening files for parsing.
//
// Use `defer CloseAll()` in [main] to make sure no files are closed before parsing is finished.
func OpenFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	OpenedFiles = append(OpenedFiles, file)
	return file, nil
}

// Close all [OpenedFiles].
//
// Recommended to be defered in [main], e.g. `defer libparser.CloseAll()`.
func CloseAll() {
	for _, file := range OpenedFiles {
		file.Close()
	}
	OpenedFiles = []File{}
}
