package cmd

import (
	"os"
	"bufio"
	"fmt"
)

func Execute() (error) {
	args := os.Args[1:]
	for _, fileName := range args {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			fmt.Printf("%v\n", scanner.Text())	
		}
	}
	return nil
}
