package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type EmptyFile struct {
    Blocks []Block
}

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

    files, emptyFiles, width := readDiskMap(scanner)
    // fmt.Println(files, width)

    // Part 1
    // diskMap := uncompressedDiskMap(files, width)
    // fmt.Println(diskMap.ToString())
    // diskMap = readressRToL(diskMap)
    // fmt.Println(diskMap.ToString())

    // Part 2
    diskMap := uncompressedDiskMap(files, width)
    // fmt.Println(diskMap.ToString())
    files, emptyFiles = readressContiguous(files, emptyFiles)
    diskMap = uncompressedDiskMap(files, width)
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

func readressContiguous(files []File, emptyFiles []EmptyFile) ([]File, []EmptyFile) {
    // fmt.Println("files", files)
    // fmt.Println("emptyFiles", emptyFiles)

    for i := len(files)-1; i >= 0; i-- {

        f := files[i]
        w := len(f.Blocks)
        for j := 0; j < len(emptyFiles); j++ {

            ef := emptyFiles[j]

            if len(ef.Blocks) < w {
                continue
            }
            if ef.Blocks[0].Id > f.Blocks[0].Id {
                break
            }

            for k := 0; k < len(f.Blocks); k++ {
                f.Blocks[k].Id = ef.Blocks[k].Id
            }
            
            emptyFiles[j].Blocks = slices.Delete(ef.Blocks, 0, w)
        }
        
    }

    // fmt.Println("files", files)
    // fmt.Println("emptyFiles", emptyFiles)

    return files, emptyFiles

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

func readDiskMap(scanner *bufio.Scanner) ([]File, []EmptyFile, int) {
    
    files := []File{}
    emptyFiles := []EmptyFile{}

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
                if w > 0 {
                    ef := EmptyFile{Blocks:make([]Block,w)}
                    for i := 0; i < len(ef.Blocks); i++ {
                        ef.Blocks[i].Id = blockCount
                        blockCount++
                    }
                    emptyFiles = append(emptyFiles, ef)
                }
                readingFile = true
            }
        }
    }

    return files, emptyFiles, blockCount

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
