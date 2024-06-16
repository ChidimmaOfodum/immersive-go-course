package interfaces

import (
	"io"
)

type FilteringPipe struct {
	PipeWriter io.Writer
}

func NewFilteringPipe(w io.Writer) *FilteringPipe {
	return &FilteringPipe{PipeWriter: w}

}

func (f *FilteringPipe) Write(b []byte) (int, error) {
	result := filterNumber(b)
	return f.PipeWriter.Write(result)
}

var filterNumber = func (input []byte) ([]byte) {
	left := 0
	for i, v:= range input {
		if !(v >= 48 && v <= 57) {
			input[left], input[i] = input[i], input[left]
			left++
		}
	}
	return input[:left]
}
