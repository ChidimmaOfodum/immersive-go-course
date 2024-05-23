package interfaces

import (
	"slices"
	"errors"
)

//OurByteBuffer - Limitations
// - the buffer grows and shrinks as data is added and read from it. The maximum data stored by the data is dependent on the machine in which the code is run. The buffer frees up all memory that has been successfully read from the buffer.
//Both read and write operations modifies the buffer, hence it is unsafe to perform them concurrently. Retrieving bytes is also unsafe to run concurrently with either read or write operations.

type OurByteBuffer struct {
	buffer []byte
}


func (b *OurByteBuffer) Write(data []byte) (int, error) {
	count := 0
	for _, v := range data {
		result := append(b.buffer, v)
		b.buffer = result
		count++
	}
	return count, nil
}

func (b *OurByteBuffer) Bytes() []byte {
	return b.buffer
}

func (b *OurByteBuffer) Read(dest []byte) (int, error) {
	var err error
	if b.buffer == nil && len(dest) != 0 {
		err = errors.New("EOF")
	}
	byteCopied := copy(dest, b.buffer)
	newBuffer := slices.Delete(b.buffer, 0, byteCopied)
	b.buffer = newBuffer
	return byteCopied, err
}


