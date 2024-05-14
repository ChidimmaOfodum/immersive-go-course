package cmd

import (
	"os"
	"bufio"
	"fmt"
	"flag"
)

func Execute() (error) {
	lineNum := flag.Bool("n", false, "Number the output lines, starting at 1")
	flag.Parse()

	args := flag.Args()
	for _, fileName := range args {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		printFileContent(*lineNum, scanner)
	}
	return nil
}

func printFileContent(displayLineNum bool, scanner *bufio.Scanner) (int) {
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if displayLineNum {
			fmt.Printf("%d  ", lineNumber)
		}
		fmt.Printf("%v\n", scanner.Text())
	}
	return lineNumber
}
