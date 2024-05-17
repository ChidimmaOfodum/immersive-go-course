package interfaces_test

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBuffer(t *testing.T) {
	t.Run("return same bytes as start byte", func(t *testing.T) {
		input := "Hello World!"
		b := bytes.NewBufferString(input)
		got := b.Bytes()
		expected := []byte(input)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v but expected %v\n", got, expected)
		}
	})

	t.Run("can write extra bytes to buffer", func(t *testing.T) {
		firstInput := "Hello "
		secondInput := "world"
		b := bytes.NewBufferString(firstInput)
		b.Write([]byte(secondInput))
		got := b.Bytes()
		expected := []byte(firstInput + secondInput)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v but expected %v\n", got, expected)
		}

	})

	t.Run("can read from buffer", func(t *testing.T) {
		input := "Hello"
		readbuf := make([]byte, len(input))
		b := bytes.NewBufferString(input)
		bytesRead, err := b.Read(readbuf)
		
		if err != nil {
			t.Errorf("expect error to be nil, got %v\n", err)
		}
		contentRead := string(readbuf)
		expectedNum := len(input)

		

		if bytesRead != expectedNum {
			t.Errorf("got %v but expected %v\n", bytesRead, expectedNum)
		}

		if input != contentRead {
			t.Errorf("got %v but expected %v\n", contentRead, input)
		}

		
	})

	t.Run("can read from buffer in bits with slice smaller than buffer content", func(t *testing.T) {
		size := 2
		readbuf := make([]byte, size)
		input := "Hello world"
		b := bytes.NewBufferString(input)

		//first call
		bytesRead, err := b.Read(readbuf)
		contentRead := string(readbuf)
		if err != nil {
			t.Errorf("expect error to be nil, got %v\n", err)
		}
		if bytesRead != size {
			t.Errorf("got %v but expected %v\n", bytesRead, size)
		}

		if contentRead != "He" {
			t.Errorf("got %v but expected %v\n", contentRead, "He")
		}

		//second call

		bytesRead, err = b.Read(readbuf)
		contentRead = string(readbuf)
		if err != nil {
			t.Errorf("expect error to be nil, got %v\n", err)
		}
		if bytesRead != size {
			t.Errorf("got %v but expected %v\n", bytesRead, size)
		}

		if contentRead != "ll" {
			t.Errorf("got %v but expected %v\n", contentRead, "He")
		}

	})

}
