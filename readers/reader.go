package readers

import (
	"strings"

	liberrors "github.com/tomefile/lib-errors"
)

type RuneReader interface {
	ReadRune() (rune, int, error)
	UnreadRune() error
	Peek(int) ([]byte, error)
}

type Reader struct {
	Inner RuneReader

	buffers          []*Buffer
	ContextBufferIdx int
	SegmentBufferIdx int

	Offset, Col, Row, PrevCol, PrevRow uint
}

func New(reader RuneReader) *Reader {
	return &Reader{
		Inner:            reader,
		buffers:          []*Buffer{},
		ContextBufferIdx: 0,
		SegmentBufferIdx: 0,
		Offset:           0,
		Col:              1,
		Row:              1,
		PrevCol:          1,
		PrevRow:          1,
	}
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

	switch len(reader.buffers) {

	case 0:
		return

	case 1:
		buffer := reader.buffers[0]
		if buffer.IsEmpty() {
			return
		}
		buffer.TrimRight(1)
		return
	}

	index := len(reader.buffers) - 1
	if reader.buffers[index].IsEmpty() {
		reader.buffers = reader.buffers[:index]
		index--
	}

	reader.buffers[index].TrimRight(1)
}

func (reader *Reader) Read() (rune, error) {
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

	if len(reader.buffers) == 0 {
		reader.buffers = append(reader.buffers, &Buffer{})
	}
	reader.buffers[len(reader.buffers)-1].Write(char)

	return char, nil
}

func (reader *Reader) MarkContext() {
	if len(reader.buffers) == 0 {
		return
	}

	reader.buffers = append(reader.buffers, &Buffer{})
	reader.ContextBufferIdx = len(reader.buffers) - 1
	reader.SegmentBufferIdx = reader.ContextBufferIdx
}

func (reader *Reader) MarkSegment() {
	if len(reader.buffers) == 0 {
		return
	}

	reader.buffers = append(reader.buffers, &Buffer{})
	reader.SegmentBufferIdx = len(reader.buffers) - 1
}

func (reader *Reader) Context() liberrors.Context {
	// TODO: Add more context if needed (sometimes the context is literally 1 line, not really that helpful)
	var ctx_buffer strings.Builder
	var ctx_highlighted strings.Builder

	if len(reader.buffers) != 0 {
		for _, buffer := range reader.buffers[reader.ContextBufferIdx:reader.SegmentBufferIdx] {
			ctx_buffer.WriteString(buffer.String())
		}
		for _, buffer := range reader.buffers[reader.SegmentBufferIdx:] {
			ctx_highlighted.WriteString(buffer.String())
		}
	}

	ctx_buffer_height := uint(strings.Count(ctx_buffer.String(), "\n"))
	ctx_highlighted_height := uint(strings.Count(ctx_highlighted.String(), "\n"))

	return liberrors.Context{
		FirstLine:   reader.PrevRow - ctx_buffer_height - ctx_highlighted_height + 1,
		Buffer:      ctx_buffer.String(),
		Highlighted: ctx_highlighted.String(),
	}
}
