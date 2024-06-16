package interfaces

import (
	"io"
)

//OurByteBuffer - Limitations
// - the buffer grows and shrinks as data is added and read from it. The maximum data stored by the buffer is dependent on the machine in which the code is run.
//Both read and write operations modifies the buffer, hence it is unsafe to perform them concurrently. Retrieving bytes is also unsafe to run concurrently with either read or write operations.

type OurByteBuffer struct {
	buffer []byte
}


func (b *OurByteBuffer) Write(data []byte) (int, error) {
	b.buffer = append(b.buffer, data...)
	return len(data), nil
}

func (b *OurByteBuffer) Bytes() []byte {
	return b.buffer
}

//Read returns io.EOF error (unless len(dest) is zero) when buffer has no data
func (b *OurByteBuffer) Read(dest []byte) (int, error) {
	var err error
	if len(b.buffer) == 0 && len(dest) != 0 {
		err = io.EOF
	}
	byteCopied := copy(dest, b.buffer)
	b.buffer = b.buffer[byteCopied:]
	return byteCopied, err
}
