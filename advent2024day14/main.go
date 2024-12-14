package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Bounds struct {
    TL Vec2
    BR Vec2
}

func (b Bounds) GetQuadrant(q int) Bounds {
    quad := Bounds{b.TL,b.BR}
    switch q {
    case 1:
        quad.BR = Vec2{b.BR.X/2,b.BR.Y/2}
    case 2:
        quad.TL = Vec2{b.BR.X/2+1,b.TL.Y}
        quad.BR = Vec2{b.BR.X,b.BR.Y/2}
    case 3:
        quad.TL = Vec2{b.TL.X,b.BR.Y/2+1}
        quad.BR = Vec2{b.BR.X/2,b.BR.Y}
    case 4:
        quad.TL = Vec2{b.BR.X/2+1,b.BR.Y/2+1}
    }
    return quad
}

type Vec2 struct {
    X int
    Y int
}

func (v Vec2) StepTo(w Vec2) Vec2 {
    return Vec2{X: w.X-v.X, Y:w.Y-v.Y}
}

func (v Vec2) Equals(w Vec2) bool {
    return v.X == w.X && v.Y == w.Y
}

func (v Vec2) Add(w Vec2) Vec2 {
    return Vec2{X:v.X+w.X, Y:v.Y+w.Y}
}

func (v Vec2) ToString() string {
    return strconv.Itoa(v.X) + "," + strconv.Itoa(v.Y)
}

func NewVec2(s string) Vec2 {
    xy := strings.Split(s, ",")
    x, _ := strconv.Atoi(xy[0])
    y, _ := strconv.Atoi(xy[1])
    return Vec2{x,y}
}

func (v Vec2) Reverse() Vec2 {
    return Vec2{X:-v.X, Y:-v.Y}
}

func (v Vec2) WithinBounds(bounds Bounds) bool {
    return v.X >= bounds.TL.X && v.X < bounds.BR.X && v.Y >= bounds.TL.Y && v.Y < bounds.BR.Y
}

type Robot struct {
    Pos Vec2
    Vec Vec2
}

func (r Robot) Move(steps int, bounds Bounds) Robot {
    
    for s := 0; s < steps; s++ {
        r.Pos = r.Pos.Add(r.Vec)
        if !r.Pos.WithinBounds(bounds) {
            if r.Pos.X < 0 {
                r.Pos.X = bounds.BR.X + r.Pos.X
            } else  {
                r.Pos.X = r.Pos.X % bounds.BR.X
            }
            if r.Pos.Y < 0 {
                r.Pos.Y = bounds.BR.Y + r.Pos.Y
            } else {
                r.Pos.Y = r.Pos.Y % bounds.BR.Y
            }
        }
    }
    return r

}

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }
    steps := 100
    if len(os.Args) > 2 {
        steps, _ = strconv.Atoi(os.Args[2])
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    bounds := Bounds{Vec2{},Vec2{101,103}}
    if strings.Contains(filename, "test") {
        bounds.BR = Vec2{11,7}
    }

    robots := loadRobots(scanner)
    // fmt.Println(robots)

    var q1,q2,q3,q4 int
    for i := 0; i < len(robots); i++ {
        robots[i] = robots[i].Move(steps, bounds)
        if robots[i].Pos.WithinBounds(bounds.GetQuadrant(1)) {
            q1++
        } else if robots[i].Pos.WithinBounds(bounds.GetQuadrant(2)) {
            q2++
        } else if robots[i].Pos.WithinBounds(bounds.GetQuadrant(3)) {
            q3++
        } else if robots[i].Pos.WithinBounds(bounds.GetQuadrant(4)) {
            q4++
        }
    }

    // fmt.Println(robots)
    fmt.Println("q1",q1,"q2",q2,"q3",q3,"q4",q4)
    fmt.Println("safety factor",q1*q2*q3*q4)

}

func loadRobots(scanner *bufio.Scanner) []Robot {

    robots := []Robot{}
    for scanner.Scan() {
        line := scanner.Text()
        rInfo := strings.Split(line, " ")
        pos := NewVec2(strings.Trim(rInfo[0], "p= "))
        vec := NewVec2(strings.Trim(rInfo[1], "v= "))
        robots = append(robots, Robot{Pos: pos, Vec: vec})
    }

    return robots
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
