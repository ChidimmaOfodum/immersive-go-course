package cmd

import (
	"bytes"
	"testing"
)

func TestListFiles(t *testing.T) {
	var b bytes.Buffer

	tests := map[string]struct {
		input []string
		want  string
		delim string
	}{
		"with delim":    {input: []string{"file1", "file2"}, want: "file1,file2", delim: ","},
		"without delim": {input: []string{"file1", "file2"}, want: "file1file2", delim: ""},
		"tab":           {input: []string{"file1", "file2"}, want: "file1\tfile2", delim: "\t"},
	}

	for _, tc := range tests {
		printFileNames(tc.input, tc.delim, &b)
		readBuf := make([]byte, b.Len())
		b.Read(readBuf)
		got := string(readBuf)

		if got != tc.want {
			t.Errorf("got %v but want %v\n", got, tc.want)
		}
	}
}
