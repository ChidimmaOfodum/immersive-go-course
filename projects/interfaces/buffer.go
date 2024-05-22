package interfaces

import (
	"slices"
	"errors"
)

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


