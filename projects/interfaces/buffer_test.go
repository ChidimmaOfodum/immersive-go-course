package interfaces

import (
	"reflect"
	"testing"
)

func TestBuffer(t *testing.T) {
	input := []byte("Hello world!")

	t.Run("return same bytes as start byte", func(t *testing.T) {
		var b OurByteBuffer
		b.Write(input)
		got := b.Bytes()
		want := input

		assertCorrectMessage(t, got, want)
	})

	t.Run("write extra bytes to buffer", func(t *testing.T) {
		var b OurByteBuffer
		firstInput := "Hello "
		secondInput := "world"
		b.Write([]byte(firstInput))
		b.Write([]byte(secondInput))
		got := b.Bytes()
		want := []byte(firstInput + secondInput)
		assertCorrectMessage(t, got, want)

	})

	t.Run("read from buffer", func(t *testing.T) {
		var b OurByteBuffer
		readbuf := make([]byte, len(input))
		b.Write(input)
		bytesRead, err := b.Read(readbuf)

		if err != nil {
			t.Errorf("expect error to be nil, got %v\n", err)
		}
		expectedNum := len(input)

		if bytesRead != expectedNum {
			t.Errorf("got %v but expected %v\n", bytesRead, expectedNum)
		}
		assertCorrectMessage(t, readbuf, input)

	})

	t.Run("read: returns error when buffer is empty", func(t *testing.T) {
		var b OurByteBuffer
		readbuf := make([]byte, len(input))
		bytesRead, err := b.Read(readbuf)

		if err.Error() != "EOF" {
			t.Errorf("expect error to be EOF, got %v\n", err.Error())
		}
		expectedNum := 0

		if bytesRead != expectedNum {
			t.Errorf("got %v but expected %v\n", bytesRead, expectedNum)
		}

	})

	t.Run("read: no error when buffer is empty and input length is 0", func(t *testing.T) {
		var b OurByteBuffer
		readbuf := make([]byte, 0)
		bytesRead, err := b.Read(readbuf)

		if err != nil {
			t.Errorf("expect error to be nil, got %v\n", err)
		}
		expectedNum := 0

		if bytesRead != expectedNum {
			t.Errorf("got %v but expected %v\n", bytesRead, expectedNum)
		}

	})

	t.Run("can read from buffer in bits with slice smaller than buffer content", func(t *testing.T) {
		var b OurByteBuffer

		size := 2
		readbuf := make([]byte, size)

		b.Write(input)

		//first call
		readCount, err := b.Read(readbuf)
		expectedReadByte := []byte{'H', 'e'} //first two letters
		if err != nil {
			t.Errorf("expect error to be nil, got %v\n", err)
		}
		if readCount != size {
			t.Errorf("got %v but expected %v\n", readCount, size)
		}
		assertCorrectMessage(t, readbuf, expectedReadByte)

		//second call

		readCount, err = b.Read(readbuf)
		expectedReadByte = []byte{'l', 'l'} //second 2 letters

		if err != nil {
			t.Errorf("expect error to be nil, got %v\n", err)
		}
		if readCount != size {
			t.Errorf("got %v but expected %v\n", readCount, size)
		}

		assertCorrectMessage(t, readbuf, expectedReadByte)
	})

}

func assertCorrectMessage(t testing.TB, got, want []byte) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v\n", got, want)
	}
}
