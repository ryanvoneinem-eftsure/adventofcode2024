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

type Dir struct {
    T Vec2
    B Vec2
    L Vec2
    R Vec2
    TL Vec2
    TR Vec2
    BL Vec2
    BR Vec2
}

func dir() Dir {
    return Dir{
        Vec2{0,-1},
        Vec2{0,1},
        Vec2{-1,0},
        Vec2{1,0},
        Vec2{-1,-1},
        Vec2{1,-1},
        Vec2{-1,1},
        Vec2{1,1},
    }
}

type Region struct {
    Plants []Plant
}

func (r Region) Area() int {
    return len(r.Plants)
}

func (r Region) Sides(game [][]rune) int {
    return r.Corners(game)
}

func (r Region) Corners(game [][]rune) int {
    if len(r.Plants) == 0 {
        return 0
    }
    if len(r.Plants) <= 2 {
        return 4
    }

    w := len(game)
    sides := 0

    for _, p := range r.Plants {

        if len(p.Neighbours) == 1 {
            // fmt.Println(p.Pos,"peninsula")
            sides += 2
            continue
        }
        if len(p.Neighbours) == 2 {
            if !p.Pos.Linear(p.Neighbours[0].Pos,p.Neighbours[1].Pos) {

                // fmt.Println(p.Pos,"corner")
                sides += 1
                // check possible interior corner
                check := p.Pos.StepTo(p.Neighbours[0].Pos).Add(p.Pos.StepTo(p.Neighbours[1].Pos))
                if getRuneAt(game, p.Pos.Add(check)) != p.Type {
                    // fmt.Println(p.Pos,"interiorcorner")
                    sides += 1
                }
            }
            continue
        }
        if len(p.Neighbours) == 3 {
            var nonLinear Plant
            linear := []Plant{}
            if p.Pos.Linear(p.Neighbours[0].Pos,p.Neighbours[1].Pos) {
                nonLinear = p.Neighbours[2]
                linear = []Plant{p.Neighbours[0],p.Neighbours[1]}
            } else if p.Pos.Linear(p.Neighbours[1].Pos,p.Neighbours[2].Pos) {
                nonLinear = p.Neighbours[0]
                linear = []Plant{p.Neighbours[1],p.Neighbours[2]}
            } else {
                nonLinear = p.Neighbours[1]
                linear = []Plant{p.Neighbours[0],p.Neighbours[2]}
            }
            checks := []Vec2{
                p.Pos.StepTo(linear[0].Pos).Add(p.Pos.StepTo(nonLinear.Pos)),
                p.Pos.StepTo(linear[1].Pos).Add(p.Pos.StepTo(nonLinear.Pos)),
            }
            for _, check := range checks {
                // fmt.Println(p.Pos,"checking",p.Pos.Add(check))
                if p.Pos.Add(check).WithinBounds(w) && getRuneAt(game, p.Pos.Add(check)) != p.Type {
                    // fmt.Println(string(getRuneAt(game, p.Pos.Add(check))),p.Pos.Add(check),"corner")
                    sides += 1
                }
            }
            continue
        }
        
        dirs := dir()
        checks := []Vec2{
            dirs.TL,
            dirs.TR,
            dirs.BL,
            dirs.BR,
        }
        for _, check := range checks {
            if p.Pos.Add(check).WithinBounds(w) && getRuneAt(game, p.Pos.Add(check)) != p.Type {
                // fmt.Println(string(getRuneAt(game, p.Pos.Add(check))),p.Pos.Add(check),"corner")
                sides += 1
            }
        }
    }

    return sides
}

type Bounds struct {
    Top int
    Bot int
    Lef int
    Rig int
}

func (r Region) Bounds() Bounds {

    var top, bot, lef, rig int
    top = 999999
    lef = 999999

    for _,pl := range r.Plants {
        if pl.Pos.Y < top {
            top = pl.Pos.Y
        }
        if pl.Pos.Y > bot {
            bot = pl.Pos.Y
        }
        if pl.Pos.X < lef {
            lef = pl.Pos.X
        }
        if pl.Pos.X > rig {
            rig = pl.Pos.X
        }
    }
    return Bounds{top, bot, lef, rig}
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

func (v Vec2) Linear(w1 Vec2, w2 Vec2) bool {
    return (v.X == w1.X && w1.X == w2.X) || (v.Y == w1.Y && w1.Y == w2.Y)
}

func getRuneAt(game [][]rune, coord Vec2) rune {
    return game[coord.Y][coord.X]
}

func setRuneAt(game [][]rune, coord Vec2, r rune) [][]rune {
    game[coord.Y][coord.X] = r
    return game
}

func scanForCoordsWithinBounds(searchChars []rune, game [][]rune, bounds Bounds) (locations [][]Vec2) {

    locations = make([][]Vec2, len(searchChars))
    for x := bounds.Lef; x <= bounds.Rig; x++ {
        for y := bounds.Top; y <= bounds.Bot; y++ {
            for i, c := range searchChars {
                if c == game[y][x] {
                    locations[i] = append(locations[i], Vec2{x, y})
                }
            }
        }
    }
    return

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
                region = Region{[]Plant{{Type:r, Pos:Vec2{x,y},Neighbours: []Plant{}}}}
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
            var neighbour Plant
            pnPos := p.Pos.Add(d)
            if pnPos.WithinBounds(w) {
                pnType := getRuneAt(game,pnPos)
                if pnType == p.Type {
                    neighbour = Plant{Type:pnType,Pos:pnPos,Neighbours:[]Plant{}}
                    neighbour.Neighbours = append(neighbour.Neighbours, p)
                    exists := false
                    for _,ns := range p.Neighbours {
                        if ns.Pos.Equals(neighbour.Pos) {
                            exists = true
                        }
                    }
                    if !exists {
                        p.Neighbours = append(p.Neighbours, neighbour)
                    }
                }
                exists := false
                for _,ep := range region.Plants {
                    if ep.Pos.Equals(neighbour.Pos) {
                        exists = true
                    }
                }
                if !exists {
                    region.Plants = append(region.Plants, neighbour)
                }
            }
        }

        region.Plants[i] = p
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
    sum := region.Area()*region.Sides(game)

    game = removeRegion(game, region)

    region = findNextRegion(game)
    for len(region.Plants) > 0 {
        
        regionCost := region.Area()*region.Sides(game)
        sum += regionCost
        game = removeRegion(game, region)
        region = findNextRegion(game)

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
