package cmd

import (
	"testing"
	"os"
)

func TestListFiles(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")

	if err != nil {
		return
	}

	defer os.RemoveAll(dir)

	
}
