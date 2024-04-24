package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func Execute() {
	help := flag.Bool("h", false, "describes how to use go-ls command")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		fmt.Println("Usage: go-ls [directories]")
		return
	}
	args := flag.Args()

	if len(args) == 0 {
		args = []string{"."}
	}

	for _, y := range args {

		fileInfo, err := os.Stat(y)

		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			listFiles(y)
		} else {
			fmt.Printf("%s\n", fileInfo.Name())
		}
	}

}

func listFiles(name string) {
	files, err := os.ReadDir(name)

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range files {
		fmt.Printf("%s\n", v.Name())
	}
}
