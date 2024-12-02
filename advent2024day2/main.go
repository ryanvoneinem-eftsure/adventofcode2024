package main

import (
    "os"
    "fmt"
    "bufio"

    "strconv"
    "strings"
)

func collectReports(scanner *bufio.Scanner) (reports [][]int) {
    for i := 0; scanner.Scan(); i++ {
        line := scanner.Text()
        numbers := strings.Split(line, " ")
        reports = append(reports, make([]int, len(numbers)))
        for j, v := range numbers {
            num, numErr := strconv.Atoi(v)
            check(numErr)
            reports[i][j] = num
        }
    }
    return
}

func reportIsSafe(report []int) (safe bool) {
    increase := 0
    decrease := 0

    for i := 1; i < len(report); i++ {
        diff := report[i] - report[i-1]

        if diff == 0 {
            return false
        } else if diff > 0 {
            increase++
            if diff > 3 {
                return false
            }
        } else if diff < 0 {
            decrease++
            if diff < -3 {
                return false
            }
        }
        
        if increase > 0 && decrease > 0 {
            return false
        }
    }

    return true
}

func safetyScore(reports [][]int) (score int) {
    for _, sl := range reports {
        if reportIsSafe(sl) {
            score++
        }
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

    reports := collectReports(scanner)

    // fmt.Println(reports)

    fmt.Printf("safety score: %v\n", safetyScore(reports))

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
