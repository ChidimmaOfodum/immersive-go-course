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

		i, err1 := os.Stat(y)

		if err1 != nil {
			log.Panic(err1)
		}
		if i.IsDir() {
			listFiles(y)
		} else {
			fmt.Printf("%s\n", i.Name())
		}

	}

}

func listFiles(name string) {
	files, err := os.ReadDir(name)

	if err != nil {
		log.Panic(err)
	}

	for _, v := range files {
		fmt.Printf("%s\n", v.Name())
	}
}
