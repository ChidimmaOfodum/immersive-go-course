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
	var result []byte
	for _, v:= range input{
		if !(v >= 48 && v<= 57) {
			result = append(result, v)
		}
	}
	return result
}

