package main

import (
	"bufio"
	"fmt"
	"os"

	"errors"
	"unicode/utf8"
    "strconv"
)

type Function struct {
	keyword    string
	paramCount int
	paramTypes []string
}

type FunctionInstance struct {
	function Function
	params   []Param
}

type Param struct {
	typeStr string
	intVal  string
}

func (fI *FunctionInstance) ToString() string {
    str := fI.function.keyword + "("
    for i, v := range fI.params {
        if i == 0 {
            str += v.intVal
        } else {
            str += ","+v.intVal
        }
    }
    return str + ")"
}

func compareRuneWithString(r rune, s string, i int) bool {

	check, _ := utf8.DecodeRuneInString(s[i:])
	return r == check

}

func findParamsFollowingFunctionCall(reader *bufio.Reader, funcDef Function) (*FunctionInstance, error) {
	allowedChars := "(),1234567890"
	paramChars := "1234567890"
	params := []Param{}

    for i, c := range paramChars {
        fmt.Printf("allowed chars: %v | %v", i, c)
    }

	runeVal, _, readErr := reader.ReadRune()
    if readErr != nil {
        return nil, readErr
    }

	if !compareRuneWithString(runeVal, allowedChars, 0) {
		return nil, errors.New("malformed")
	}

	runeVal, _, readErr = reader.ReadRune()
    if readErr != nil {
        return nil, readErr
    }

	currentParam := Param{"int", ""}

	for readErr == nil {
		fmt.Println(runeVal)

		if compareRuneWithString(runeVal, allowedChars, 1) {
            fmt.Println(") found")
            if len(currentParam.intVal) == 0 || len(params) < funcDef.paramCount-1 {
				return nil, errors.New("malformed")
            }

            params = append(params, currentParam)

            funcCall := FunctionInstance{function: funcDef, params: params}
			return &funcCall, nil
		}

		if compareRuneWithString(runeVal, allowedChars, 2) {
            fmt.Println(", found")
			if len(currentParam.intVal) > 0 {
				params = append(params, currentParam)
                fmt.Printf("len:%v\n", len(params))
                fmt.Printf("max:%v\n", funcDef.paramCount)
				if len(params) > funcDef.paramCount {
					return nil, errors.New("malformed")
				}
				currentParam = Param{"int", ""}

                runeVal, _, readErr = reader.ReadRune()
                continue
			} else {
                return nil, errors.New("malformed")
            }
		}

		found := false
		for i, c := range paramChars {
			if runeVal == rune(c) {
				found = true
                fmt.Printf("%v == %v\n", runeVal, c)
				currentParam.intVal += paramChars[i : i+1]
			}
		}
		
        if !found && compareRuneWithString(runeVal, allowedChars, 2) {
            fmt.Println(", found")
			if len(currentParam.intVal) > 0 {
				params = append(params, currentParam)
                fmt.Printf("len:%v\n", len(params))
                fmt.Printf("max:%v\n", funcDef.paramCount)
				if len(params) > funcDef.paramCount {
					return nil, errors.New("malformed")
				}
				currentParam = Param{"int", ""}

			} else {
                return nil, errors.New("malformed")
            }
		} else if !found {
            fmt.Println("something weird found")
			return nil, errors.New("malformed")
		}

		runeVal, _, readErr = reader.ReadRune()
	}

    return nil, readErr
}

func findNextFunctionCall(reader *bufio.Reader, funcDef Function) (*FunctionInstance, error) {
	runeVal, _, readErr := reader.ReadRune()
    if readErr != nil {
        return nil, readErr
    }

	functionNameLen := len(funcDef.keyword)
	functionCharsFound := 0

	for readErr == nil {
		fmt.Println(runeVal)
		if compareRuneWithString(runeVal, funcDef.keyword, functionCharsFound) {
			fmt.Printf("%#U found\n", runeVal)
			functionCharsFound++
			if functionCharsFound == functionNameLen {
				fmt.Println("'mul' found")
				functionCharsFound = 0
				return findParamsFollowingFunctionCall(reader, funcDef)
			}
		}

		runeVal, _, readErr = reader.ReadRune()
	}

	return nil, readErr
}

func mulDefinition(funcI *FunctionInstance) int {
    answer := 0
    for i, v := range funcI.params {
        if i == 0 {
            answer, _ = strconv.Atoi(v.intVal)
        } else {
            iVal, _ := strconv.Atoi(v.intVal)
            answer = answer * iVal
        }
    }
    return answer
}

func main() {

	filename := "./input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	file, reader := createFileReader(filename)
	defer file.Close()

	mulFunc := Function{"mul", 2, []string{"int", "int"}}

    callables := []*FunctionInstance{}
    call, err := findNextFunctionCall(reader, mulFunc)
    if err == nil {
        callables = append(callables, call)
    }

    for err == nil {
        call, err = findNextFunctionCall(reader, mulFunc)
        if call != nil {
            fmt.Printf("created %v\n", call.ToString())
        }
        if err != nil && err.Error() == "malformed" {
            err = nil
            err = reader.UnreadRune()
        } else if err == nil {
            fmt.Printf("appending %v\n", call.ToString())
            callables = append(callables, call)
        }
    }

    if err != nil && err.Error() != "EOF" && err.Error() != "malformed" {
        check(err)
    } else if err != nil {
        fmt.Printf("success: %v\n", err)
    }

    finalAnswer := 0
    fmt.Printf("callables: %v\n", len(callables))
    for _, callable := range callables {
        fmt.Println(callable.ToString())
        finalAnswer += mulDefinition(callable)
    }

    fmt.Printf("final answer: %v\n", finalAnswer)

}

func createFileReader(filename string) (file *os.File, reader *bufio.Reader) {

	file, fileErr := os.Open(filename)
	check(fileErr)

	fmt.Println("file opened")

	reader = bufio.NewReader(file)
	return

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
