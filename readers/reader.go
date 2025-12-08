package readers

import (
	"strings"

	liberrors "github.com/tomefile/lib-errors"
)

const NUM_OF_CONTEXT_LINES = 3

type RuneReader interface {
	ReadRune() (rune, int, error)
	UnreadRune() error
	Peek(int) ([]byte, error)
}

type Reader struct {
	Inner RuneReader

	buffer []rune

	Offset, Col, Row, PrevCol, PrevRow uint
}

func New(reader RuneReader) *Reader {
	return &Reader{
		Inner:   reader,
		buffer:  make([]rune, 0, 1_000),
		Offset:  0,
		Col:     1,
		Row:     1,
		PrevCol: 1,
		PrevRow: 1,
	}
}

func (reader *Reader) Previous() rune {
	if len(reader.buffer) <= 1 {
		return 0
	}
	return reader.buffer[len(reader.buffer)-2]
}

// Returns the next byte without advancing the reader
func (reader *Reader) Peek() (byte, error) {
	data, err := reader.Inner.Peek(1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (reader *Reader) Unread() {
	reader.Inner.UnreadRune()
	reader.Offset--
	reader.buffer = reader.buffer[:len(reader.buffer)-1]
}

func (reader *Reader) Read() (rune, error) {
	char, size, err := reader.Inner.ReadRune()
	if err != nil {
		return 0, err
	}

	reader.Offset++
	reader.PrevCol = reader.Col
	reader.PrevRow = reader.Row

	if char == '\n' {
		reader.Row++
		reader.Col = 0
	} else {
		reader.Col += uint(size)
	}

	reader.buffer = append(reader.buffer, char)

	return char, nil
}

func (reader *Reader) Context(at uint) liberrors.Context {
	if len(reader.buffer) == 0 {
		return liberrors.Context{
			FirstLine:   1,
			Buffer:      "",
			Highlighted: "",
		}
	}

	at = min(uint(len(reader.buffer)-1), at)
	lines := strings.Split(string(reader.buffer[:at]), "\n")
	delim := max(0, len(lines)-1-NUM_OF_CONTEXT_LINES)
	ctx_lines := lines[delim:]

	highlighted := string(reader.buffer[at:])
	if len(strings.TrimSpace(highlighted)) == 0 {
		highlighted = ""
	}

	return liberrors.Context{
		FirstLine:   uint(len(lines)-len(ctx_lines)) + 1,
		Buffer:      strings.Join(ctx_lines, "\n"),
		Highlighted: highlighted,
	}
}
