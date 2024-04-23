package cmd

import (
	"fmt"
	"os"
	"log"
	"flag"
	"strings"
)
func Execute() {
	help := flag.Bool("help", false, "describes how to use go-ls command")
    flag.Parse()

    if *help {
        flag.PrintDefaults()
        fmt.Println("Usage: go-ls [directories]")
        return
    }
	//args := os.Args[1:]
	args := flag.Args()
	
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
