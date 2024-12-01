package main

import (
    "os"
    "fmt"
    "bufio"

    "strings"
    "strconv"
    "sort"
)

func collectLists(scanner *bufio.Scanner) (l, r []int) {
    
    l = []int{}
    r = []int{}
    for scanner.Scan() {
        line := scanner.Text()
        nums := strings.Split(line, "   ")
        lN, _ := strconv.Atoi(nums[0])
        rN, _ := strconv.Atoi(nums[1])
        l = append(l, lN)
        r = append(r, rN)
    }

    return
}

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    l, r := collectLists(scanner)
    // fmt.Printf("list l = %v\n", l)
    // fmt.Printf("list r = %v\n", r)
	sort.Ints(l)
	sort.Ints(r)
	
	totalDiffs := 0
	for i, v := range l {
        diff := v - r[i]
        if diff < 0 {
            diff = -diff
        }
        totalDiffs += diff
	}

    fmt.Printf("answer: %v\n", totalDiffs)

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
