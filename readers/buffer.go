package readers

type Buffer struct {
	inner []rune
}

func (buffer *Buffer) IsEmpty() bool {
	return len(buffer.inner) == 0
}

func (buffer *Buffer) TrimRight(amount int) {
	buffer.inner = buffer.inner[:len(buffer.inner)-amount]
}

func (buffer *Buffer) Write(char rune) {
	buffer.inner = append(buffer.inner, char)
}

func (buffer *Buffer) String() string {
	return string(buffer.inner)
}
