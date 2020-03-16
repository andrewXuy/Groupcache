package cache

// hold an immutable view of bytes
// it can support any type data e.t. Image, string...
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return clone(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func clone(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
