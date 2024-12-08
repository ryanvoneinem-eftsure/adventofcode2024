package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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

func (v Vec2) WithinBounds(width int) bool {
    return v.X >= 0 && v.X < width && v.Y >= 0 && v.Y < width
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
    for x := 0; x < len(game); x++ {
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

func findAntinodesInBoundsPart1(locations []Vec2, width int) []Vec2 {
    antinodes := []Vec2{}
    for i := 0; i < len(locations)-1; i++ {
        for j := i+1; j < len(locations); j++ {
            step := locations[i].StepTo(locations[j])
            pAnti1 := locations[j].Add(step)
            if pAnti1.WithinBounds(width) {
                antinodes = append(antinodes, pAnti1)
            }
            pAnti2 := locations[i].Add(step.Reverse())
            if pAnti2.WithinBounds(width) {
                antinodes = append(antinodes, pAnti2)
            }
        }
    }
    return antinodes
}

func createAntinodeMap(width int, locations []Vec2) [][]rune {

    game := make([][]rune, width)
    for i := 0; i < len(game); i++ {
        game[i] = []rune(strings.Repeat(".", width))
    }
    for _, loc := range locations {
        game = setRuneAt(game, loc, '#')
    }

    return game

}

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    game := loadGameBoard(scanner)
    possibleChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

    locations := scanForCoords(possibleChars, game)
    locations = slices.DeleteFunc(locations, func(lSl []Vec2)(bool){
        return len(lSl) == 0
    })
    fmt.Println(locations)

    antinodes := []Vec2{}
    for _, lSl := range locations {
        antinodes = slices.Concat(antinodes, findAntinodesInBoundsPart1(lSl, len(game)))
    }

    // fmt.Println(antinodes)
    // fmt.Println(len(antinodes))

    antinodeMap := createAntinodeMap(len(game), antinodes)

    for _, l := range antinodeMap {
        fmt.Println(string(l))
    }

    count := len(scanForCoords([]rune{'#'}, antinodeMap)[0])
    fmt.Println(count)

}

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

    // fmt.Printf("%vx%v\n", len(game[0]), len(game))
    return

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

