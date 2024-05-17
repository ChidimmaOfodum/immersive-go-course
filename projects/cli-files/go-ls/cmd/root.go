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

	//flags
	help := flag.Bool("h", false, "describes how to use go-ls command")
	comma:= flag.Bool("m", false, "prints with a delimiter")
	flag.Parse()

	delimiter := "\t"

	if *help {
		flag.PrintDefaults()
		fmt.Println("Usage: go-ls [directories]")
		return nil
	}
	if *comma {
		delimiter = ","
	}

	args := flag.Args()

	if len(args) == 0 {
		args = []string{"."}
	}

	for _, filePath := range args {

		fileInfo, err = os.Stat(filePath)

		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			files, err := os.ReadDir(filePath)

			if err != nil {
				return err
			}
			listFiles(files, delimiter)
		} else {
			fmt.Printf("%s\n", fileInfo.Name()) // print filename if not directory
		}
	}
	return nil
}

func listFiles(fileNames []fs.DirEntry, delim string) {
	for index, value := range fileNames {
		fmt.Printf("%s", value.Name())
		if index != len(fileNames) - 1 {
			fmt.Printf("%s ", delim)
		}
	}
}
