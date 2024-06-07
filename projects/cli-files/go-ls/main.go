package main

import (
	"go-ls/cmd"
	"log"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
