package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Trail struct {
    Steps []Vec2
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

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    game := loadGameBoard(scanner)
    trailheads := scanForCoords([]rune{'0'}, game)[0]
    fmt.Println(trailheads)
    trails := make([][]Trail, len(trailheads))
    for i, th := range trailheads {
        trails[i] = append(trails[i], exploreTrails(th, trails[i], game)...)
    }

    sum := 0
    for i := range trailheads {
        counter := 0
        for j, t := range trails[i] {
            fmt.Println("trailhead", i, "trail", j, t)
            counter++
        }
        fmt.Println("count:", counter)
        sum += counter
    }

    fmt.Println("sum:", sum)
}

func exploreTrails(trailhead Vec2, trails []Trail, game [][]rune) []Trail {

    completeTrails := []Trail{}
    trails = slices.DeleteFunc(trails, func(t Trail)(bool){
        return !slices.Contains(t.Steps, trailhead)
    })

    h, _ := strconv.Atoi(string(getRuneAt(game, trailhead)))

    directions := []Vec2{
        {0,-1},
        {1,0},
        {0,1},
        {-1,0},
    }

    nextSteps := []Vec2{}
    for _, d := range directions {

        nextStep := trailhead.Add(d)
        if nextStep.WithinBounds(len(game)) {
            h2, _ := strconv.Atoi(string(getRuneAt(game, nextStep)))
            if h2 == h+1 {
                nextSteps = append(nextSteps, nextStep)
                trails = append(trails, Trail{[]Vec2{trailhead}})
            }
        }
    }

    for i, ns := range nextSteps {

        h2, _ := strconv.Atoi(string(getRuneAt(game, ns)))
        trails[i].Steps = append(trails[i].Steps, ns)

        if h2 == 9 {

            completeTrails = append(completeTrails, trails[i])

        } else {

            completeTrails = append(completeTrails, exploreTrails(ns, trails, game)...)

        }
    }

    // uncomment deduplication logic for Part 1

    // peaksSeen := []Vec2{}
    // dupes := []int{}
    // for i, c := range completeTrails {
    //     if slices.Contains(peaksSeen, c.Steps[len(c.Steps)-1]) {
    //         dupes = append(dupes, i)
    //     } else {
    //         peaksSeen = append(peaksSeen, c.Steps[len(c.Steps)-1])
    //     }
    // }
    // for i := len(dupes)-1; i >= 0; i-- {
    //
    //     completeTrails = slices.Delete(completeTrails, dupes[i], dupes[i]+1)
    //
    // }

    return completeTrails
    
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
