package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

    // Part 1
    // steps := 100
    // if len(os.Args) > 2 {
    //     steps, _ = strconv.Atoi(os.Args[2])
    // }
    
    maxSteps := 100
    if len(os.Args) > 2 {
        maxSteps, _ = strconv.Atoi(os.Args[2])
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    bounds := Bounds{Vec2{},Vec2{101,103}}
    if strings.Contains(filename, "test") {
        bounds.BR = Vec2{11,7}
    }

    robots := loadRobots(scanner)
    fmt.Println(len(robots))

    // Part 1
    // var q1,q2,q3,q4 int
    // for i := 0; i < len(robots); i++ {
    //     robots[i] = robots[i].Move(steps, bounds)
    //     if robots[i].Pos.WithinBounds(bounds.GetQuadrant(1)) {
    //         q1++
    //     } else if robots[i].Pos.WithinBounds(bounds.GetQuadrant(2)) {
    //         q2++
    //     } else if robots[i].Pos.WithinBounds(bounds.GetQuadrant(3)) {
    //         q3++
    //     } else if robots[i].Pos.WithinBounds(bounds.GetQuadrant(4)) {
    //         q4++
    //     }
    // }
    // fmt.Println("q1",q1,"q2",q2,"q3",q3,"q4",q4)
    // fmt.Println("safety factor",q1*q2*q3*q4)

    board := make([][]bool, bounds.BR.Y)
    for i := 0; i < len(board); i++ {
        board[i] = make([]bool, bounds.BR.X)
    }

    // Part 2
    steps := 0
    for i := 0; i < maxSteps; i++ {
        for j := 0; j < len(robots); j++ {
            robots[j] = robots[j].Move(1, bounds)
            board[robots[j].Pos.Y][robots[j].Pos.X] = true
        }
        // fmt.Println("step",i+1)
        // printBoard(board)
        if findFrame(3, board) {
            steps = i+1
            printBoard(board)
            break
        }
        board = resetBoard(board)
    }

    fmt.Println("steps",steps)

}

func resetBoard(board [][]bool) [][]bool {

    for y := 0; y < len(board); y++ {
        for x := 0; x < len(board[0]); x++ {
            board[y][x] = false
        }
    }
    return board
}

func printBoard(board [][]bool) {
    for y := 0; y < len(board); y++ {
        string := ""
        for x := 0; x < len(board[0]); x++ {
            if board[y][x] {
                string += "#"
            } else {
                string += " "
            }
        }
        fmt.Println(string)
    }
}

func findFrame(heuristic int, board [][]bool) bool {

    // fmt.Println(board)

    frameCorners := [8]*Vec2{}
    for y := 0; y < len(board); y++ {
        count := 0;
        var start *Vec2
        for x := 0; x < len(board[0]); x++ {
            if board[y][x] {
                if count == 0 {
                    start = &Vec2{x,y}
                }
                count++
            } else {
                if count >= heuristic {
                    if frameCorners[0] == nil {
                        frameCorners[0] = start
                        frameCorners[1] = &Vec2{x-1,y}
                    } else {
                        frameCorners[2] = start
                        frameCorners[3] = &Vec2{x-1,y}
                    }
                    // fmt.Println("framesides",frameCorners)
                }
                count = 0
            }
        }
    }

    for x := 0; x < len(board[0]); x++ {
        count := 0;
        var start *Vec2
        for y := 0; y < len(board); y++ {
            if board[y][x] {
                if count == 0 {
                    start = &Vec2{x,y}
                }
                count++
            } else {
                if count >= heuristic {
                    if frameCorners[4] == nil {
                        frameCorners[4] = start
                        frameCorners[5] = &Vec2{x,y-1}
                    } else {
                        frameCorners[6] = start
                        frameCorners[7] = &Vec2{x,y-1}
                    }
                    // fmt.Println("framesides",frameCorners)
                }
                count = 0
            }
        }
    }

    return makesRect(frameCorners[0:8])

}

func makesRect(points []*Vec2) bool {
    
    // fmt.Println(points)
    if slices.IndexFunc(points, func(p *Vec2)(bool){return p == nil}) != -1 {
        return false
    }
    // fmt.Println(points)
    return points[0].Equals(*points[4]) && points[1].Equals(*points[6]) && points[2].Equals(*points[5])

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
