package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stone struct {
    Mark string
}

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }
    blinks := 1
    if len(os.Args) > 2 {
        blinks , _ = strconv.Atoi(os.Args[2])
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    scanner.Scan()
    line := scanner.Text()
    fmt.Println(line)
    stones := make(map[Stone]int)
    for _, s := range strings.Split(line, " ") {
        stones[Stone{s}] += 1
    }


    for b := 0; b < blinks; b++ {
        newStones := make(map[Stone]int)
        for s,c := range stones {
            newS := s.updateStone()
            if newS.Mark != "" {
                newStones[newS] += c
            }
            newStones[s] += c
        }
        stones = newStones

        // for s, c := range stones {
        //     fmt.Println(b, s.Mark, c)
        // }
    }

    sum := 0
    for _, c := range stones {
        // fmt.Println(s.Mark, c)
        sum += c
    }
    fmt.Println("sum:",sum)

}

func (s *Stone) updateStone() Stone {
    // fmt.Println("stone",s)

    var sPlusOne Stone
    if s.Mark == "0" {

        s.Mark = "1"

    } else if len(s.Mark) % 2 == 0 {

        l := len(s.Mark)
        markRunes := []rune(s.Mark)
        s.Mark = string(markRunes[:(l/2)])
        sPlusOne = Stone{string(markRunes[(l/2):])}

        if rune(sPlusOne.Mark[0]) == '0' {
            sPlusOne.Mark = strings.TrimLeft(sPlusOne.Mark, "0")
            if sPlusOne.Mark == "" {
                sPlusOne.Mark = "0"
            }
        }

    } else {
        mInt, _ := strconv.Atoi(s.Mark)
        mInt *= 2024
        s.Mark = strconv.Itoa(mInt)
    }

    return sPlusOne

    // fmt.Println("calculated", result)
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
