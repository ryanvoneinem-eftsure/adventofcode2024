package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

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

func (v Vec2) Reverse() Vec2 {
    return Vec2{X:-v.X, Y:-v.Y}
}

func (v Vec2) WithinBounds(bounds Bounds) bool {
    return v.X >= bounds.TL.X && v.X < bounds.BR.X && v.Y >= bounds.TL.Y && v.Y < bounds.BR.Y
}

func getPlayerMovement(player rune) Vec2 {
    switch player {
    case '>':
        return Vec2{X:1,Y:0}
    case 'v':
        return Vec2{X:0,Y:1}
    case '<':
        return Vec2{X:-1,Y:0}
    default:
        return Vec2{X:0,Y:-1}
    }
}

func scanColForCoords(searchChars []rune, game [][]rune, x int) (locations []Vec2) {

    for y := 0; y < len(game); y++ {
        for _, c := range searchChars {

            if c == game[y][x] {
                locations = append(locations, Vec2{x, y})
            }

        }
    }
    return

}

func scanRowForCoords(searchChars []rune, game [][]rune, y int) (locations []Vec2) {

    for x := 0; x < len(game[0]); x++ {
        for _, c := range searchChars {

            if c == game[y][x] {
                locations = append(locations, Vec2{x, y})
            }

        }
    }
    return

}

func scanColForCoordsBounded(searchChars []rune, game [][]rune, x int, b Bounds) (locations []Vec2) {

    for y := b.TL.Y; y < b.BR.Y; y++ {
        for _, c := range searchChars {

            if c == game[y][x] {
                locations = append(locations, Vec2{x, y})
            }

        }
    }
    return

}

func scanRowForCoordsBounded(searchChars []rune, game [][]rune, y int, b Bounds) (locations []Vec2) {

    for x := b.TL.X; x < b.BR.X; x++ {
        for _, c := range searchChars {

            if c == game[y][x] {
                locations = append(locations, Vec2{x, y})
            }

        }
    }
    return

}

func searchPath(loc, dir Vec2, game [][]rune) Vec2 {

    if dir.X == 0 {

        obstructions := scanColForCoords([]rune{'.','#'}, game, loc.X)
        seeker := Vec2{Y:loc.Y}
        for seeker.Y >= 0 && seeker.Y < len(game) {
            for _, obs := range obstructions {
                if obs.Y == seeker.Y {
                    return obs
                }
            }
            seeker.Y += dir.Y
        }

    } else {

        obstructions := scanRowForCoords([]rune{'.','#'}, game, loc.Y)
        seeker := Vec2{X:loc.X}
        for seeker.X >= 0 && seeker.X < len(game[0]) {
            for _, obs := range obstructions {
                if obs.X == seeker.X {
                    return obs
                }
            }
            seeker.X += dir.X
        }

    }

    return Vec2{-1,-1}

}

func getRuneAt(game [][]rune, coord Vec2) rune {
    return game[coord.Y][coord.X]
}

func setRuneAt(game [][]rune, coord Vec2, r rune) [][]rune {
    game[coord.Y][coord.X] = r
    return game
}

func scanForCoords(searchChars []rune, game [][]rune) (locations [][]Vec2) {

    locations = make([][]Vec2, len(searchChars))
    for x := 0; x < len(game[0]); x++ {
        for y := 0; y < len(game); y++ {
            for i, c := range searchChars {
                if c == game[y][x] {
                    locations[i] = append(locations[i], Vec2{x, y})
                }
            }
        }
    }
    return

}

type Bounds struct {
    TL Vec2
    BR Vec2
}

type Box Bounds

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }
    steps := -1
    if len(os.Args) > 2 {
        s := os.Args[2]
        steps, _ = strconv.Atoi(s)
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    board, moves := loadBoardAndMovesPart2(scanner)
    // printBoard(board)

    robot := scanForCoords([]rune{'@'}, board)[0][0]

    for m := 0; m < len(moves); m++ {
        // fmt.Println("move",m)
        board, robot = moveRobotPart2(board, robot, moves[m])
        // printBoard(board)
        if m == steps {
            panic(m)
        }
    }

    sum := solveForPart2(board)

    fmt.Println("sum",sum)

}

func solveForPart1(board [][]rune) int {

    boxes := scanForCoords([]rune{'O'}, board)[0]
    sum := 0
    for _,b := range boxes {
        sum += (100*b.Y)+b.X
    }
    return sum

}

func solveForPart2(board [][]rune) int {

    boxes := scanForCoords([]rune{'['}, board)[0]
    sum := 0
    for _,b := range boxes {
        sum += (100*b.Y)+b.X
    }
    return sum

}

func moveRobotPart1(board [][]rune, robot Vec2, moveR rune) ([][]rune, Vec2) {

    move := getPlayerMovement(moveR)
    path := searchPath(robot, move, board)
    // fmt.Println(robot, move, path)

    if !path.Equals(Vec2{-1,-1}) && getRuneAt(board, path) != '#' {

        backwards := move.Reverse()
        for !path.Equals(robot) {

            board = setRuneAt(board, path, getRuneAt(board, path.Add(backwards)))
            path = path.Add(backwards)

        }
        
        board = setRuneAt(board, robot, '.')
        robot = robot.Add(move)

    }
    

    return board, robot
}

func moveRobotPart2(board [][]rune, robot Vec2, moveR rune) ([][]rune, Vec2) {

    move := getPlayerMovement(moveR)
    path := searchPath(robot, move, board)
    // fmt.Println(robot, move, path)

    if !path.Equals(Vec2{-1,-1}) && getRuneAt(board, path) != '#' {

        if move.Y == 0 || getRuneAt(board, robot.Add(move)) == '.' {

            backwards := move.Reverse()
            for !path.Equals(robot) {

                board = setRuneAt(board, path, getRuneAt(board, path.Add(backwards)))
                path = path.Add(backwards)

            }
            
            board = setRuneAt(board, robot, '.')
            robot = robot.Add(move)

        } else {

            searchBounds := Bounds{}
            if move.Y > 0 {
                searchBounds.TL.Y = robot.Y
                searchBounds.BR.Y = path.Y
            } else {
                searchBounds.TL.Y = path.Y
                searchBounds.BR.Y = robot.Y
            }
            boxSearch := scanColForCoordsBounded([]rune{'[',']'}, board, robot.X, searchBounds)
            boxSearch = slices.DeleteFunc(boxSearch, func(v Vec2)(bool){
                if move.Y > 0 {
                    return v.Y <= robot.Y
                } else {
                    return v.Y >= robot.Y
                }
            })
            expandedPath := make(map[Vec2]int)
            boxes := []Box{}

            for bi := 0; bi < len(boxSearch); bi++ {
                b := boxSearch[bi]

                box := Box{}
                box.TL.Y = b.Y
                box.BR.Y = b.Y
                pathSearch := Vec2{-1,-1}
                if getRuneAt(board, b) == '[' {
                    box.TL.X = b.X
                    box.BR.X = b.X+1
                    pathSearch = searchPath(box.BR, move, board)
                    if pathSearch.Equals(Vec2{-1,-1}) || getRuneAt(board, pathSearch) == '#' {
                        return board, robot
                    }
                    if move.Y > 0 {
                        searchBounds.TL.Y = box.TL.Y
                        searchBounds.BR.Y = pathSearch.Y
                    } else {
                        searchBounds.TL.Y = pathSearch.Y
                        searchBounds.BR.Y = box.BR.Y
                    }
                    bs := scanColForCoordsBounded([]rune{'[',']'}, board, box.BR.X, searchBounds)
                    bs = slices.DeleteFunc(bs, func(v Vec2)(bool){
                        if move.Y > 0 {
                            return v.Y <= box.TL.Y
                        } else {
                            return v.Y >= box.TL.Y
                        }
                    })
                    boxSearch = append(boxSearch, bs...)
                } else {
                    box.TL.X = b.X-1
                    box.BR.X = b.X
                    pathSearch = searchPath(box.TL, move, board)
                    if pathSearch.Equals(Vec2{-1,-1}) || getRuneAt(board, pathSearch) == '#' {
                        return board, robot
                    }
                    if move.Y > 0 {
                        searchBounds.TL.Y = box.TL.Y
                        searchBounds.BR.Y = pathSearch.Y
                    } else {
                        searchBounds.TL.Y = pathSearch.Y
                        searchBounds.BR.Y = box.BR.Y
                    }
                    bs := scanColForCoordsBounded([]rune{'[',']'}, board, box.TL.X, searchBounds)
                    bs = slices.DeleteFunc(bs, func(v Vec2)(bool){
                        if move.Y > 0 {
                            return v.Y <= box.TL.Y
                        } else {
                            return v.Y >= box.TL.Y
                        }
                    })
                    boxSearch = append(boxSearch, bs...)
                }

                if pathSearch.Equals(Vec2{-1,-1}) || getRuneAt(board, pathSearch) == '#' {
                    return board, robot
                }

                if expandedPath[pathSearch] == 0 {
                    expandedPath[pathSearch] = box.BR.Y
                } else if move.Y > 0 && expandedPath[pathSearch] > box.BR.Y {
                    expandedPath[pathSearch] = box.BR.Y
                } else if move.Y < 0 && expandedPath[pathSearch] < box.BR.Y {
                    expandedPath[pathSearch] = box.BR.Y
                }
                // fmt.Println("es",expandedPath)
                boxes = append(boxes, box)

            }

            backwards := move.Reverse()

            for pv,s := range expandedPath {

                // fmt.Println("bvpv",pv)
                stop := func(m Vec2, p int, s int)(bool){
                    if m.Y > 0 {
                        return p < s
                    } else {
                        return p > s
                    }
                }
                for !pv.Equals(path) && pv.Y != s {

                    // fmt.Println("moving",pv.Add(backwards),"to",pv)
                    replacement := getRuneAt(board, pv.Add(backwards))
                    if stop(move, pv.Add(backwards).Y, s) {
                        replacement = '.'
                    }
                    board = setRuneAt(board, pv, replacement)
                    // printBoard(board)
                    pv = pv.Add(backwards)

                }

                board = setRuneAt(board, pv, '.')

            }

            for !path.Equals(robot) {

                board = setRuneAt(board, path, getRuneAt(board, path.Add(backwards)))
                path = path.Add(backwards)

            }
            
            board = setRuneAt(board, robot, '.')
            robot = robot.Add(move)

        }

    }
    

    return board, robot
}

func printBoard(board [][]rune) {
    for y := 0; y < len(board); y++ {
        s := string(board[y])
        fmt.Println(s)
    }
}

func loadBoardAndMovesPart1(scanner *bufio.Scanner) ([][]rune, []rune) {

    board := [][]rune{}
    moves := []rune{}
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            // do nothing
        } else if line[0] == '#' {
            board = append(board, []rune{})
            i := len(board)-1
            board[i] = append(board[i], []rune(line)...)
        } else {
            moves = append(moves, []rune(line)...)
        }
    }

    return board, moves

}

func loadBoardAndMovesPart2(scanner *bufio.Scanner) ([][]rune, []rune) {

    board := [][]rune{}
    moves := []rune{}
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            // do nothing
        } else if line[0] == '#' {
            board = append(board, []rune{})
            j := len(board)-1
            for i := 0; i < len(line); i++ {
                r := rune(line[i])
                if r == 'O' {
                    board[j] = append(board[j], '[', ']')
                } else if r == '@' {
                    board[j] = append(board[j], '@', '.')
                } else {
                    board[j] = append(board[j], r, r)
                }
            }
        } else {
            moves = append(moves, []rune(line)...)
        }
    }

    return board, moves

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

