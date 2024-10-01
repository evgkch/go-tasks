//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	args := os.Args[1:]
	m := make(map[string]uint)

	for _, path := range args {
		file, err := os.Open(path)

		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			m[line] += 1
		}

		file.Close()
	}

	for key, value := range m {
		if value > 1 {
			fmt.Printf("%d\t%s\n", value, key)
		}
	}
}
