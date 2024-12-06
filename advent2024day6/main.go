package main

import (
    "os"
    "fmt"
    "bufio"
)

type Player struct {
    Loc Vec2
    Dir Vec2
    Rune rune
}

type Vec2 struct {
    X int
    Y int
}

func getPlayer(game [][]rune) (player Player) {
    player.Loc = scanForCoords([]rune{rune("^"[0]), rune(">"[0]), rune("v"[0]), rune("<"[0])}, game)[0]
    player.Rune = getRuneAt(game, player.Loc)
    player.Dir = getPlayerMovement(player.Rune)
    return
}

func getRuneAt(game [][]rune, coord Vec2) rune {
    return game[coord.Y][coord.X]
}

func rotateVec2Right(v Vec2) Vec2 {
    if v.X == 1 {
        return Vec2{X:0,Y:1}

    } else if v.Y == 1 {
        return Vec2{X:-1,Y:0}

    } else if v.X == -1 {
        return Vec2{X:0,Y:-1}

    } else {
        return Vec2{X:1,Y:0}
    }
}

func getMovementPlayer(v Vec2) rune {
    if v.X == 1 {
        return '>'
    } else if v.Y == 1 {
        return 'v'
    } else if v.X == -1 {
        return '<'
    } else {
        return '^'
    }
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

func searchPath(loc, dir Vec2, game [][]rune) Vec2 {

    if dir.X == 0 {

        obstructions := scanColForCoords([]rune{'#'}, game, loc.X)
        seeker := Vec2{loc.X, loc.Y}
        for seeker.Y >= 0 && seeker.Y < len(game) {
            for _, obs := range obstructions {
                if obs.Y == seeker.Y {
                    return obs
                }
            }
            seeker.Y += dir.Y
        }

    } else {

        obstructions := scanRowForCoords([]rune{'#'}, game, loc.Y)
        seeker := Vec2{loc.X, loc.Y}
        for seeker.X >= 0 && seeker.X < len(game) {
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

func walkPlayer(game *[][]rune, player Player, obs Vec2) Player {

    if player.Dir.X == 0 {

        for player.Loc.Y != obs.Y-1 && player.Loc.Y != obs.Y+1 {
            (*game)[player.Loc.Y][player.Loc.X] = 'X'
            player.Loc.Y += player.Dir.Y
        }

    } else {

        for player.Loc.X != obs.X-1 && player.Loc.X != obs.X+1 {
            (*game)[player.Loc.Y][player.Loc.X] = 'X'
            player.Loc.X += player.Dir.X
        }

    }

    player.Dir = rotateVec2Right(player.Dir)
    (*game)[player.Loc.Y][player.Loc.X] = getMovementPlayer(player.Dir)

    return player
    
}

func scanColForCoordsBounded(searchChars []rune, game [][]rune, x int, startY int, endY int) (locations []Vec2) {

    for y := startY; y < len(game) && y <= endY; y++ {
        for _, c := range searchChars {

            if c == game[y][x] {
                locations = append(locations, Vec2{x, y})
            }

        }
    }
    return
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

func scanRowForCoordsBounded(searchChars []rune, game [][]rune, y int, startX int, endX int) (locations []Vec2) {

    for x := startX; x < len(game) && x <= endX; x++ {
        for _, c := range searchChars {

            if c == game[y][x] {
                locations = append(locations, Vec2{x, y})
            }

        }
    }
    return
}

func scanRowForCoords(searchChars []rune, game [][]rune, y int) (locations []Vec2) {

    for x := 0; x < len(game); x++ {
        for _, c := range searchChars {

            if c == game[y][x] {
                locations = append(locations, Vec2{x, y})
            }

        }
    }
    return

}

func scanForCoords(searchChars []rune, game [][]rune) (locations []Vec2) {

    for x := 0; x < len(game); x++ {
        for y := 0; y < len(game); y++ {
            for _, c := range searchChars {

                if c == game[y][x] {
                    locations = append(locations, Vec2{x, y})
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
    player := getPlayer(game)

    fmt.Println(player)

    uniqueLocations := 0

    for player.Loc.X != -1 {
        obstruction := searchPath(player.Loc, player.Dir, game)
        fmt.Println(obstruction)

        if obstruction.X != -1 {
            player = walkPlayer(&game, player, obstruction)
        } else {
            remainingUniqueSteps := 0
            if player.Dir.X == 0 {
                if player.Dir.Y == 1 {
                    remainingUniqueSteps = len(scanColForCoordsBounded([]rune{'.'}, game, player.Loc.X, player.Loc.Y, len(game)))
                } else {
                    remainingUniqueSteps = len(scanColForCoordsBounded([]rune{'.'}, game, player.Loc.X, 0, player.Loc.Y))
                }
            } else {
                if player.Dir.X == 1 {
                    remainingUniqueSteps = len(scanRowForCoordsBounded([]rune{'.'}, game, player.Loc.Y, player.Loc.X, len(game)))
                } else {
                    remainingUniqueSteps = len(scanRowForCoordsBounded([]rune{'.'}, game, player.Loc.Y, 0, player.Loc.X))
                }
            }
            uniqueLocations += remainingUniqueSteps
            game[player.Loc.Y][player.Loc.X] = 'X'
            player.Loc.X = -1
        }

    }

    uniqueLocations += len(scanForCoords([]rune{'X'}, game))
    fmt.Println(uniqueLocations)

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
