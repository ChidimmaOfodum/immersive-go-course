package interfaces

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFilteringPipe(t *testing.T) {
	t.Run("initializes correctly", func(t *testing.T) {
		var writer bytes.Buffer
		myFilteringPipe := NewFilteringPipe(&writer)

		if myFilteringPipe.PipeWriter == nil {
			t.Errorf("expected filtering pipe to be initialized but got nil")
		}

		if myFilteringPipe.PipeWriter != &writer {
			t.Errorf("Expected writer to  be %v but got %v\n", writer, myFilteringPipe.PipeWriter)
		}
	})

	t.Run("calls filterNumber", func(t *testing.T) {
		var writer bytes.Buffer
		myFilteringPipe := NewFilteringPipe(&writer)

		called := 0
		originalFilterNumber := filterNumber
		testFilterNumber := func(input []byte) []byte {
			called += 1
			return originalFilterNumber(input)
		}
		filterNumber = testFilterNumber
		myFilteringPipe.Write([]byte("start=1, end=10"))

		if called != 1 {
			t.Errorf("filterNumber was not called")
		}
	})
}

func TestFilterNumbers(t *testing.T) {

	tests := map[string]struct {
		input []byte
		want  []byte
	}{
		"without numbers": {input: []byte("I am a string without numbers"), want: []byte("I am a string without numbers")},
		"with numbers":    {input: []byte("start=1, end=10"), want: []byte("start=, end=")},
		"empty": {input: []byte{}, want: []byte{}},
		"only numbers" : {input: []byte("1233455566"), want: []byte{}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := filterNumber(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected %v but got %v\n", string(tc.want), string(got))
			}
		})
	}

}
