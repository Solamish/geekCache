package DistributedCache

type ByteView struct {
	b []byte
}

func(v ByteView) Len() int {
	return len(v.b)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	 return string(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c,b)
	return c
}
