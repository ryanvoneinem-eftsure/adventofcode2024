package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	filename := "./input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	file, scanner := createFileScanner(filename)
	defer file.Close()

	rules := collectOrderingRules(scanner)
	fmt.Println(rules)

}

func collectOrderingRules(scanner *bufio.Scanner) (rules []string) {
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			return
		}
		rules = append(rules, line)
	}
	return

}

func createFileScanner(filename string) (file *os.File, scanner *bufio.Scanner) {

	file, fileErr := os.Open(filename)
	check(fileErr)

	fmt.Println("file opened")

	scanner = bufio.NewScanner(file)
	return

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
