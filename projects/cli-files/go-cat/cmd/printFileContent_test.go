package cmd

import (
	"bufio"
	"strings"
	"testing"
)

func TestPrintFileContent(t *testing.T) {
	t.Run("Without line numbers", func (t *testing.T) {
		scanner := bufio.NewScanner(strings.NewReader("Hello\nWorld"))
		got := printFileContent(false, scanner)
		expected := 2
		if expected != got {
			t.Errorf("got %v expected %v", got, expected)
		}
	})
}
