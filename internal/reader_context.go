package internal

import (
	liberrors "github.com/tomefile/lib-errors"
)

func (reader *SourceCodeReader) ContextReset() {
	reader.context.Reset()
	reader.buffer.Reset()
}

func (reader *SourceCodeReader) ContextBookmark() {
	reader.context.WriteString(reader.buffer.String())
	reader.buffer.Reset()
}

func (reader *SourceCodeReader) ErrorContext() liberrors.Context {
	return liberrors.Context{
		FirstLine:   reader.PrevRow,
		Buffer:      reader.context.String(),
		Highlighted: reader.buffer.String(),
	}
}
