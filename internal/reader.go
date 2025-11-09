package internal

import (
	"bufio"
	"strings"
)

type SourceCodeReader struct {
	Inner *bufio.Reader

	context strings.Builder
	buffer  strings.Builder

	Offset, Col, Row, PrevCol, PrevRow uint
}

func NewSourceCodeReader(reader *bufio.Reader) *SourceCodeReader {
	return &SourceCodeReader{
		Inner:   reader,
		context: strings.Builder{},
		buffer:  strings.Builder{},
		Offset:  0,
		Col:     1,
		Row:     1,
		PrevCol: 1,
		PrevRow: 1,
	}
}

func (reader *SourceCodeReader) Peek() (byte, error) {
	data, err := reader.Inner.Peek(1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (reader *SourceCodeReader) Read() (rune, error) {
	char, size, err := reader.Inner.ReadRune()
	if err != nil {
		return 0, err
	}

	reader.Offset += uint(size)

	reader.PrevCol = reader.Col
	reader.PrevRow = reader.Row

	if char == '\n' {
		reader.Row++
		reader.Col = 0
	} else {
		reader.Col += uint(size)
	}

	reader.buffer.WriteRune(char)
	return char, nil
}
