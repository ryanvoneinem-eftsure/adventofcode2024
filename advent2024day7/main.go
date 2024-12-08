package main

import (
	"bufio"
	"fmt"
	"math"
	"slices"

	// "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Problem struct {
    Vars []int
    Answer int
}

var wg sync.WaitGroup

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    problems := collectProblems(scanner)
    ch := make(chan int, 850)

    // Part1
    // operators := []rune{'+','*'}

    // Part2
    operators := []string{"+","*","||"}

    for _, p := range problems {
        fmt.Println(p)
        go solve(*p, operators, ch)
        wg.Add(1)
    }

    wg.Wait()
    close(ch)

    sum := 0
    for o := range ch {

        // fmt.Println(o)
        sum += o

    }

    fmt.Println("sum:",sum)

}

func solve(problem Problem, operators []string, ch chan int) {
    defer wg.Done()

    // id := rand.Int()

    var i int64
    for i = 0; i < int64(math.Pow(float64(len(operators)), float64(len(problem.Vars)-1))); i++ {
        places := strconv.FormatInt(i, len(operators))
        operands := []string{}
        if len(places) < len(problem.Vars)-1 {
            places = strings.Repeat("0", len(problem.Vars)-1-len(places)) + places
        }
        // fmt.Println(problem.Vars, places)
        for _, p := range places {
            pInt, _ := strconv.Atoi(string(p))
            operands = append(operands, operators[pInt])
        }

        // fmt.Println(id, "variables:",problem.Vars)
        // fmt.Println(id, "operands:",operands)

        variables := slices.Clone(problem.Vars)
        // fmt.Println(id, "variables:", variables)

        answer := variables[0]
        for j := 1; j < len(variables); j++ {
            switch operands[j-1] {
            case "+":
                answer += variables[j]
            case "*":
                answer *= variables[j]
            case "||":
                // fmt.Println(id, "combining", variables)
                newVar, _ := strconv.Atoi(strconv.Itoa(answer) + strconv.Itoa(variables[j]))
                if len(variables) == 2 {
                    variables = []int{newVar}
                    operands = []string{}
                    answer = variables[0]
                } else {
                    variables = slices.Delete(variables, j-1, j+1)
                    variables = slices.Insert(variables, j-1, newVar)
                    operands = slices.Delete(operands, j-1, j)
                    answer = variables[j-1]
                    j--
                }
                // fmt.Println(id, "combined", variables)
            }
        }
        if answer == problem.Answer {
            ch <- answer
            return
        }
    }
    
    ch <- 0

}

func collectProblems(scanner *bufio.Scanner) []*Problem {
    problems := []*Problem{}
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > 0 {
            problem := Problem{}
            sides := strings.Split(line, ":")
            answer, _ := strconv.Atoi(sides[0])
            problem.Answer = answer

            variables := strings.Split(strings.TrimSpace(sides[1]), " ")
            for _, v := range variables {
                intV, _ := strconv.Atoi(v)
                problem.Vars = append(problem.Vars, intV)
            }
            problems = append(problems, &problem)
        }

    }
    
    return problems
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
