package internal

import (
	"fmt"
	"io"
	"strings"

	libescapes "github.com/bbfh-dev/lib-ansi-escapes"
)

func (reader *SourceCodeReader) ContextReset() {
	reader.context.Reset()
	reader.buffer.Reset()
}

func (reader *SourceCodeReader) ContextBookmark() {
	reader.context.WriteString(reader.buffer.String())
	reader.buffer.Reset()
}

func (reader *SourceCodeReader) PrintContext(writer io.Writer) {
	writer.Write([]byte(libescapes.TextColorWhite))

	var i uint
	for line := range strings.SplitSeq(reader.context.String(), "\n") {
		if i != 0 {
			writer.Write([]byte{'\n'})
		}
		fmt.Fprintf(writer, "%5d |  %s", reader.PrevRow+i, line)
		i++
	}

	writer.Write([]byte(libescapes.TextColorBrightRed))

	if reader.buffer.Len() == 0 {
		writer.Write([]byte("<empty string>"))
	} else {
		i = 0
		buffer := strings.TrimSuffix(reader.buffer.String(), "\n")
		for line := range strings.SplitSeq(buffer, "\n") {
			if i == 0 {
				writer.Write([]byte(line))
			} else {
				fmt.Fprintf(writer, "\n%5d |  %s", reader.PrevRow+i, line)
			}
			i++
		}

	}

	writer.Write([]byte(libescapes.ColorReset))
}

func (reader *SourceCodeReader) GetPrintedContext() string {
	var builder strings.Builder
	reader.PrintContext(&builder)
	return builder.String()
}
