package cmd

import (
	"bytes"
	"testing"
)

func TestPrintFileNames(t *testing.T) {
	

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
		var b bytes.Buffer
		printFileNames(tc.input, tc.delim, &b)
		got := b.String()
		if got != tc.want {
			t.Errorf("got %v but want %v\n", got, tc.want)
		}
	}
}
