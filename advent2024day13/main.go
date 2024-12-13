package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func (v Vec2) StepTo(w Vec2) Vec2 {
    return Vec2{X: w.X-v.X, Y:w.Y-v.Y}
}

func (v Vec2) Equals(w Vec2) bool {
    return v.X == w.X && v.Y == w.Y
}

func (v Vec2) Add(w Vec2) Vec2 {
    return Vec2{X:v.X+w.X, Y:v.Y+w.Y}
}

func (v Vec2) Multiply(m int) Vec2 {
    if m == 0 {
        return Vec2{0,0}
    }
    for i := 1; i < m; i++ {
        v = v.Add(v)
    }
    return v
}

func (v Vec2) ToString() string {
    return strconv.FormatFloat(v.X, 'd', 9, 64) + "," + strconv.FormatFloat(v.Y, 'd', 9, 64)
}

type Vec2 struct {
    X float64
    Y float64
}

type Button struct {
    Move Vec2
    Cost int
}

type Machine struct {
    Id int
    ButtonA Button
    ButtonB Button
    Prize Vec2
}

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    machines := loadMachines(scanner)

    sum := 0
    for _,m := range machines {
        cost := calculateCost(m)
        if cost != -1 {
            sum += cost
        }
    }
    fmt.Println("sum:",sum)
}

func calculateCost(m Machine) int {
    
    a1 := m.ButtonA.Move.X
    a2 := m.ButtonA.Move.Y
    b1 := m.ButtonB.Move.X
    b2 := m.ButtonB.Move.Y
    c1 := m.Prize.X + 10000000000000
    c2 := m.Prize.Y + 10000000000000

    denom := (a1 * b2) - (b1 * a2)
    aPress := (c1 * b2) - (b1*c2)
    aPress = aPress / denom
    bPress := (a1 * c2) - (c1 * a2)
    bPress = bPress / denom

    _, aFrac := math.Modf(aPress)
    _, bFrac := math.Modf(bPress)
    if aFrac == 0.0 && bFrac == 0.0 {
        return (int(aPress) * m.ButtonA.Cost) + (int(bPress) * m.ButtonB.Cost)
    }
    
    return -1

}

func loadMachines(scanner *bufio.Scanner) []Machine {

    machines := []Machine{}
    counter := 1
    m := Machine{Id:counter}
    for scanner.Scan() {
        line := scanner.Text()

        if len(line) > 0 {
            bInfo := strings.Split(line, ":")
            values := strings.Split(bInfo[1], ",")
            xStr := strings.Trim(values[0], "X+= ")
            xInt, _ := strconv.Atoi(xStr)
            yStr := strings.Trim(values[1], "Y+= ")
            yInt, _ := strconv.Atoi(yStr)
            if bInfo[0] == "Button A" {
                m.ButtonA = Button{Cost:3, Move:Vec2{float64(xInt),float64(yInt)}}
            } else if bInfo[0] == "Button B" {
                m.ButtonB = Button{Cost:1, Move:Vec2{float64(xInt),float64(yInt)}}
            } else {
                m.Prize = Vec2{float64(xInt),float64(yInt)}
            }
        } else {
            machines = append(machines, m)
            counter++
            m = Machine{Id:counter}
        }
    }

    if m.Id > 0 {
            machines = append(machines, m)
    }

    return machines

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
