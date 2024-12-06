package main

import (
    "os"
    "fmt"
    "bufio"
)

func loadGameBoard(scanner *bufio.Scanner) (game [][]rune) {

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
        lineSl := make([]rune, len(line))
        for i, c := range line {
            lineSl[i] = rune(c)
        }
        game = append(game, lineSl)
    }

    fmt.Printf("%vx%v\n", len(game[0]), len(game))
    return

}

type Coords struct {
    X int
    Y int
}

type Bounds struct {
    TopL Coords
    TopR Coords
    BotL Coords
    BotR Coords
}

func findPossibleXmasBounds(game [][]rune) (candidates []Bounds) {

    coords := scanHorizontalForCoords([]rune("A")[0], game, 1)
    // fmt.Println(len(coords))
    for _, c := range coords {
        candidate := Bounds{
            TopL:Coords{c.X-1,c.Y-1},
            TopR:Coords{c.X+1,c.Y-1},
            BotL:Coords{c.X-1,c.Y+1},
            BotR:Coords{c.X+1,c.Y+1},
        }
        // fmt.Println(candidate)
        candidates = append(candidates, candidate)
    }

    return

}

func scanHorizontalForCoords(searchChar rune, game [][]rune, buffer int) (locations []Coords) {

    for x := buffer; x < len(game)-buffer; x++ {
        for y := buffer; y < len(game)-buffer; y++ {

            if searchChar == game[y][x] {
                locations = append(locations, Coords{x, y})
            }

        }
    }
    return

}

func getGameRegion(game [][]rune, bounds Bounds) (region [][]rune) {
    
    for y := bounds.TopL.Y; y <= bounds.BotL.Y; y++ {
        line := []rune{}
        for x := bounds.TopL.X; x <= bounds.TopR.X; x++ {
            line = append(line, game[y][x])
        }
        // fmt.Println(line)
        region = append(region, line)
    }
    return

}

func scanHorizontal(candidates []string, game [][]rune) (counter int) {

    registers := make([]rune, len(candidates[0]))

    for x := 0; x < len(game); x++ {
        for y := 0; y < len(game); y++ {
            walk:
            for w := 0; w < len(registers); w++ {
                xw := x+w
                if xw == len(game) {
                    break walk
                }
                if w == len(registers)-1 {
                    registers[w] = game[y][xw]
                    matches := matches(candidates, registers)

                    // fmt.Printf("%v, %v-%v\n", x, y, yw)
                    // fmt.Println(matches)
                    for _, m := range matches {
                        if m {
                            counter++
                        }
                    }
                }

                registers[w] = game[y][xw]

            }
        }
    }
    return

}

func scanVertical(candidates []string, game [][]rune) (counter int) {

    registers := make([]rune, len(candidates[0]))

    for x := 0; x < len(game); x++ {
        for y := 0; y < len(game); y++ {
            walk:
            for w := 0; w < len(registers); w++ {
                yw := y+w
                if yw == len(game) {
                    break walk
                }
                if w == len(registers)-1 {
                    registers[w] = game[yw][x]
                    matches := matches(candidates, registers)

                    // fmt.Printf("%v, %v-%v\n", x, y, yw)
                    // fmt.Println(matches)
                    for _, m := range matches {
                        if m {
                            counter++
                        }
                    }
                }

                registers[w] = game[yw][x]

            }
        }
    }
    return

}

func scanDiagonalForward(candidates []string, game [][]rune) (counter int) {

    registers := make([]rune, len(candidates[0]))

    for x := len(game)-1; x >= 0; x-- {
        for y := 0; y < len(game); y++ {
            walk:
            for w := 0; w < len(registers); w++ {
                xw := x-w
                yw := y+w
                if xw < 0 || yw == len(game)  {
                    break walk
                }
                if w == len(registers)-1 {
                    registers[w] = game[yw][xw]
                    matches := matches(candidates, registers)

                    // fmt.Printf("%v, %v-%v\n", x, y, yw)
                    // fmt.Println(matches)
                    for _, m := range matches {
                        if m {
                            counter++
                        }
                    }
                }

                registers[w] = game[yw][xw]

            }
        }
    }
    return

}

func scanDiagonalBackward(candidates []string, game [][]rune) (counter int) {

    registers := make([]rune, len(candidates[0]))

    for x := 0; x < len(game); x++ {
        for y := 0; y < len(game); y++ {
            walk:
            for w := 0; w < len(registers); w++ {
                xw := x+w
                yw := y+w
                if xw == len(game) || yw == len(game)  {
                    break walk
                }
                if w == len(registers)-1 {
                    registers[w] = game[yw][xw]
                    matches := matches(candidates, registers)

                    // fmt.Printf("%v, %v-%v\n", x, y, yw)
                    // fmt.Println(matches)
                    for _, m := range matches {
                        if m {
                            counter++
                        }
                    }
                }

                registers[w] = game[yw][xw]

            }
        }
    }
    return

}

func matches(candidates []string, registers []rune) []bool {

    found := make([]bool, len(candidates))
    for i, _ := range found {
        found[i] = true
    }

    for x, c := range candidates {
        inner:
        for y, r := range registers {
            // fmt.Printf("Comparing %#U and %#U\n", rune(c[y]), r)
            // if rune(c[y]) == r {
            //     fmt.Println("match")
            // }
            if found[x] && rune(c[y]) != r {
                // fmt.Println("not a match")
                found[x] = false
            }
            if !found[x] {
                break inner
            }
        }
    }

    return found
}

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    game := loadGameBoard(scanner)

    // Part 1:
    // targets := []string{"XMAS","SAMX"}
    // count := 0
    // count += scanVertical(targets, game)
    // count += scanHorizontal(targets, game)
    // count += scanDiagonalForward(targets, game)
    // count += scanDiagonalBackward(targets, game)

    // Part 2:

    candidates := findPossibleXmasBounds(game)
    targets := []string{"MAS","SAM"}
    count := 0
    for _, c := range candidates {
        box := getGameRegion(game, c)
        tempCount := scanDiagonalForward(targets, box)
        tempCount += scanDiagonalBackward(targets, box)
        if tempCount == 2 {
            count++
        }
    }
    fmt.Println(count)

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

