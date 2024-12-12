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

type Region struct {
    Plants []Plant
}

func (r Region) Area() int {
    return len(r.Plants)
}

func (r Region) Perimeter() int {
    basis := 4
    sum := 0
    for _,p := range r.Plants {
        sum += basis - len(p.Neighbours)
    }
    return sum
}

type Plant struct {
    Type rune
    Pos Vec2
    Neighbours []Plant
}

func (p Plant) ToString() string {
    return string(p.Type) + ":" + p.Pos.ToString()
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

func findNextRegion(game [][]rune) Region {

    w := len(game)
    plants := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
    var r rune
    region := Region{}
    outer:
    for y := 0; y < w; y++ {
        for x := 0; x < w; x++ {
            r = getRuneAt(game, Vec2{x,y})
            valid := false
            for _,p := range plants {
                if r == p {
                    valid = true
                    break
                }
            }

            if valid {
                region = Region{[]Plant{{Type:r, Pos:Vec2{x,y}}}}
                break outer
            }
        }
    }

    for i := 0; i < len(region.Plants); i++ {
        
        p := region.Plants[i]
        directions := []Vec2{
            {0,1},
            {0,-1},
            {1,0},
            {-1,0},
        }
        for _,d := range directions {
            exists := false
            var neighbour Plant
            pnPos := p.Pos.Add(d)
            if pnPos.WithinBounds(w) {
                pnType := getRuneAt(game,pnPos)
                if pnType == p.Type {
                    neighbour = Plant{Type:pnType,Pos:pnPos,Neighbours:[]Plant{p}}
                    exists := false
                    for _,n := range p.Neighbours {
                        if n.Pos.Equals(neighbour.Pos) {
                            exists = true
                        }
                    }
                    if !exists {
                        p.Neighbours = append(p.Neighbours, neighbour)
                        region.Plants[i] = p
                    }
                }
                for _,p := range region.Plants {
                    if p.Pos.Equals(neighbour.Pos) {
                        exists = true
                    }
                }
                if !exists {
                    region.Plants = append(region.Plants, neighbour)
                }
            }
        }

    }

    region.Plants = slices.DeleteFunc(region.Plants, func(p Plant)(bool){
        return p.Type != r
    })

    return region

}

func removeRegion(game [][]rune, region Region) [][]rune {

    for _,p := range region.Plants {
        game = setRuneAt(game, p.Pos, '.')
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
    
    region := findNextRegion(game)
    game = removeRegion(game, region)
    sum := region.Area()*region.Perimeter()

    for len(region.Plants) > 0 {
        
        region = findNextRegion(game)
        game = removeRegion(game, region)
        regionCost := region.Area()*region.Perimeter()
        sum += regionCost

    }

    fmt.Println(sum)

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

    fmt.Printf("%vx%v\n", len(game[0]), len(game))
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
