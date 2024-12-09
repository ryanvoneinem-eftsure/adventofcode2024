package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type File struct {
    Id int
    Blocks []Block
}

type Block struct {
    Id int
}

type DiskMap []string

func (dm DiskMap) ToString() string {
    diskString := ""
    for _, i := range dm {
        if i == "" {
            diskString += "."
        } else  {
            diskString += i
        }
    }
    return diskString
}

func main() {

    filename := "./input.txt"
    if len(os.Args) > 1 {
        filename = os.Args[1]
    }

    file, scanner := createFileScanner(filename)
	defer file.Close()

    files, width := readDiskMap(scanner)
    // fmt.Println(files, width)
    diskMap := uncompressedDiskMap(files, width)
    // fmt.Println(diskMap.ToString())
    diskMap = readressRToL(diskMap)
    // fmt.Println(diskMap.ToString())

    fmt.Println("sum: ", checksum(diskMap))
}

func checksum(diskMap DiskMap) int {
    var sum int
    for i, a := range diskMap {
        aInt, _ := strconv.Atoi(string(a))
        sum += aInt * i
    }
    return sum
}

func readressRToL(diskMap DiskMap) DiskMap {

    emptySpaces := []int{}
    for i, a := range diskMap {
        if a == "" {
            emptySpaces = append(emptySpaces, i)
        }
    }
    e := len(diskMap) - 1
    for len(emptySpaces) > 0 && e >= 0 && e > emptySpaces[0] {
        if diskMap[e] == "" {
            e--
            continue
        }

        diskMap[emptySpaces[0]] = diskMap[e]
        diskMap[e] = ""

        emptySpaces = slices.Delete(emptySpaces, 0, 1)
        e--

    }
    return diskMap
}

func uncompressedDiskMap(files []File, width int) DiskMap {
    diskMap := make([]string, width)
    for _, f := range files {
        for _, b := range f.Blocks {
            diskMap[b.Id] = strconv.Itoa(f.Id)
        }
    }
    return diskMap
}

func readDiskMap(scanner *bufio.Scanner) ([]File, int) {
    
    files := []File{}
    fileCount := 0
    blockCount := 0
    readingFile := true
    for scanner.Scan() {
        line := scanner.Text()
        for _, r := range []rune(line) {
            w, _ := strconv.Atoi(string(r))
            if readingFile {
                f := File{Id:fileCount, Blocks:make([]Block, w)}
                fileCount++
                
                for i := 0; i < len(f.Blocks); i++ {
                    f.Blocks[i].Id = blockCount
                    blockCount++
                }

                files = append(files, f)
                readingFile = false
            } else {
                blockCount += w
                readingFile = true
            }
        }
    }

    return files, blockCount

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
