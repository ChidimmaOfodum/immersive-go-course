package cmd

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
)

func Execute()(error) {
	var fileInfo fs.FileInfo
	var err error

	help := flag.Bool("h", false, "describes how to use go-ls command")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		fmt.Println("Usage: go-ls [directories]")
		return nil
	}
	args := flag.Args()

	if len(args) == 0 {
		args = []string{"."}
	}

	for _, y := range args {

		fileInfo, err = os.Stat(y)

		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if err := listFiles(y); err != nil {
				return err
			}
		} else {
			fmt.Printf("%s\n", fileInfo.Name())
		}
	}
	return nil
}

func listFiles(name string) (error) {
	files, err := os.ReadDir(name)
	for _, v := range files {
		fmt.Printf("%s\n", v.Name())
	}
	return err
}
