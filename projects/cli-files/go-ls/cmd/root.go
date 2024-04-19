package cmd

import (
	"fmt"
	"os"
	"log"
	"strings"
)
func Execute() {
	args := os.Args[1:]
	
	if len(args) == 0 {
		args = []string{"."}
	}
	

	for _, y:= range args {
		files, err := os.ReadDir(y)

		if err != nil {
			errMessage := err.Error()
			if (strings.Contains(errMessage, "fdopendir")) {
				fmt.Println(y)
				continue
			} else {
				log.Fatal(err)
			}
			
		}

		for _, file := range files {
			fmt.Printf("%s\n", file.Name())
		}
	}

}
