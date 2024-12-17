package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Instruction [2]int

var A int
var B int
var C int
var ptr int
var output []string

func combop(combo int) (int, error) {
    switch combo {
    case 0,1,2,3:
        return combo, nil
    case 4:
        return A, nil
    case 5:
        return B, nil
    case 6:
        return C, nil
    default:
        return -1, errors.New("unrecognised combo")

    }
}

func adv(combo int) {
    // fmt.Println("executing 0 adv:",combo)
    combo, comboErr := combop(combo)
    if comboErr != nil {
        fmt.Println(comboErr.Error())
        return
    }
    d := int(math.Pow(2, float64(combo)))
    A = A / d
}

func bxl(literal int) {
    // fmt.Println("executing 1 bxl:",literal)
    B = B ^ literal
}

func bst(combo int) {
    // fmt.Println("executing 2 bst:",combo)
    combo, comboErr := combop(combo)
    if comboErr != nil {
        fmt.Println(comboErr.Error())
        return
    }
    B = combo % 8
}

func jnz(literal int) {
    // fmt.Println("executing 3 jnz:",literal)
    if A == 0 {
        return
    }
    if literal > 0 {
        literal /= 2
    }
    ptr = literal-1
}

func bxc(ignore int) {
    // fmt.Println("executing 4 bxc:",ignore)
    // fmt.Println("ignoring",ignore)
    B = B ^ C
}

func out(combo int) {
    // fmt.Println("executing 5 out:",combo)
    combo, comboErr := combop(combo)
    if comboErr != nil {
        fmt.Println(comboErr.Error())
        return
    }
    output = append(output, strconv.Itoa(combo % 8))
}

func bdv(combo int) {
    // fmt.Println("executing 6 bdv:",combo)
    combo, comboErr := combop(combo)
    if comboErr != nil {
        fmt.Println(comboErr.Error())
        return
    }
    d := int(math.Pow(2, float64(combo)))
    B = A / d
}

func cdv(combo int) {
    // fmt.Println("executing 7 cdv:",combo)
    combo, comboErr := combop(combo)
    if comboErr != nil {
        fmt.Println(comboErr.Error())
        return
    }
    d := int(math.Pow(2, float64(combo)))
    C = A / d
}

var functionMap map[int]func(int)

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()
    
    scanner.Scan()
    aLine := scanner.Text()
    A, _ = strconv.Atoi(strings.TrimSpace(strings.Split(aLine, ":")[1]))
    scanner.Scan()
    bLine := scanner.Text()
    B, _ = strconv.Atoi(strings.TrimSpace(strings.Split(bLine, ":")[1]))
    scanner.Scan()
    cLine := scanner.Text()
    C, _ = strconv.Atoi(strings.TrimSpace(strings.Split(cLine, ":")[1]))

    scanner.Scan()
    scanner.Scan()
    instructionLine := strings.TrimSpace(strings.Split(scanner.Text(), ":")[1])
    instructions := loadInstructions(instructionLine)

    fmt.Println("A:",A,"\nB:",B,"\nC:",C,"\n",instructions)

    functionMap = make(map[int]func(int))
    functionMap[0] = adv
    functionMap[1] = bxl
    functionMap[2] = bst
    functionMap[3] = jnz
    functionMap[4] = bxc
    functionMap[5] = out
    functionMap[6] = bdv
    functionMap[7] = cdv

    ptr = 0
    for ptr < len(instructions) {
        i := instructions[ptr]
        functionMap[i[0]](i[1])
        ptr++
        // fmt.Println("output so far:",strings.Join(output, ","))
    }

    fmt.Println("output:",strings.Join(output, ","))

}

func executeInstructions(instructions []Instruction) {


}

func loadInstructions(line string) []Instruction {

    out := []Instruction{}
    ints := strings.Split(line,",")
    for i := 0; i < len(ints); i += 2 {
        operation, _ := strconv.Atoi(ints[i])
        operand, _ := strconv.Atoi(ints[i+1])
        
        out = append(out, Instruction{operation,operand})
    }

    return out

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
