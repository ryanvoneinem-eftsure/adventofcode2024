package main

import (
    "os"
    "fmt"
    "bufio"
)

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    scanner := createFileScanner(filename)


}

func createFileScanner(filename string) scanner bufio.Scanner {
    
    file, fileErr := os.Open(filename)
    check(err)
    defer file.Close()

    fmt.Println("file opened")

    scanner := bufio.NewScanner(file)
    return

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
