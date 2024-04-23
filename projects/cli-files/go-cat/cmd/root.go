package cmd

import (
	"log"
	"os"
)

func Execute() {
	args := os.Args[1:]
	for _, v := range args {
		data, err := os.ReadFile(v)

		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(data)
	}
}
