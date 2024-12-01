package main

import (
    "os"
    "fmt"
    "bufio"

    "strings"
    "strconv"
    // "sort"
)

func collectNumbers(scanner *bufio.Scanner) (l []int, r map[int]int) {
    
    l = []int{}
    r = make(map[int]int)
    for scanner.Scan() {
        line := scanner.Text()
        nums := strings.Split(line, "   ")
        lN, _ := strconv.Atoi(nums[0])
        rN, _ := strconv.Atoi(nums[1])
        l = append(l, lN)
        r[rN] += 1
    }

    return
}

func calcSimilarity(nums []int, counts map[int]int) (score int) {
    for _, v := range nums {
        score += v * counts[v]
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

    l, r := collectNumbers(scanner)
	
    answer := calcSimilarity(l, r)

    fmt.Printf("answer: %v\n", answer)

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
