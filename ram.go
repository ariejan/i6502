package i6502

/*
Random Access Memory, read/write, can be of any size.
*/
type Ram struct {
	data []byte
}

// Create a new Ram component of the given size.
func NewRam(size int) (*Ram, error) {
	return &Ram{data: make([]byte, size)}, nil
}

func (r *Ram) Size() uint16 {
	return uint16(len(r.data))
}

func (r *Ram) Read(address uint16) byte {
	return r.data[address]
}

func (r *Ram) Write(address uint16, data byte) {
	r.data[address] = data
}
