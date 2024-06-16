package interfaces

import (
	"errors"
	"io"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {
	input := []byte("Hello world!")

	t.Run("return same bytes as start byte", func(t *testing.T) {
		var b OurByteBuffer
		b.Write(input)
		got := b.Bytes()
		want := input

		require.Equal(t, want, got)
	})

	t.Run("write extra bytes to buffer", func(t *testing.T) {
		var b OurByteBuffer
		firstInput := "Hello "
		secondInput := "world"
		b.Write([]byte(firstInput))
		b.Write([]byte(secondInput))
		got := b.Bytes()
		want := []byte(firstInput + secondInput)
		require.Equal(t, want, got)

	})

	t.Run("read from buffer", func(t *testing.T) {
		var b OurByteBuffer
		readbuf := make([]byte, len(input))
		b.Write(input)
		want := len(input)
		
		got, err := b.Read(readbuf)
		require.Nil(t, err)
		require.Equal(t, want, got)

	})

	t.Run("read: returns EOF error when buffer is empty", func(t *testing.T) {
		var b OurByteBuffer
		readbuf := make([]byte, len(input))
		bytesRead, err := b.Read(readbuf)
		require.Equal(t, err, io.EOF)

		if !errors.Is(err, io.EOF) {
			t.Errorf("expect error to be EOF, got %v\n", err)
		}
		require.Equal(t, 0, bytesRead)

	})

	t.Run("read: no error when buffer is empty and input length is 0", func(t *testing.T) {
		var b OurByteBuffer
		readbuf := make([]byte, 0)
		bytesRead, err := b.Read(readbuf)
		require.Nil(t, err)
		require.Equal(t, 0, bytesRead)
	})

	t.Run("can read from buffer in bits with slice smaller than buffer content", func(t *testing.T) {
		var b OurByteBuffer

		size := 2
		readbuf := make([]byte, size)

		b.Write(input)

		//first call
		readCount, err := b.Read(readbuf)
		expectedReadByte := []byte{'H', 'e'} //first two letters
		require.Nil(t, err)
		require.Equal(t, size, readCount)
		require.Equal(t, expectedReadByte, readbuf)
		
		//second call

		readCount, err = b.Read(readbuf)
		expectedReadByte = []byte{'l', 'l'} //second 2 letters
		require.Nil(t, err)
		require.Equal(t, size, readCount)
		require.Equal(t, expectedReadByte, readbuf)
	})

	t.Run("can read from buffer with slice bigger than buffer content", func(t *testing.T) {
		var b OurByteBuffer
		readbuf := make([]byte, len(input) * 2)
		b.Write(input)
		bytesRead, err := b.Read(readbuf)
		require.Nil(t, err)
		require.Equal(t, len(input), bytesRead)
		require.Equal(t, input, readbuf[:len(input)])
	})
}
