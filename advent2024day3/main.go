package main

import (
	"bufio"
	"fmt"
	"os"

	"errors"
	// "unicode/utf8"
    "strconv"
)

type Function struct {
	keyword    string
	paramCount int
	paramTypes []string
}

type FunctionInstance struct {
	function Function
	params   []*Param
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

type TokeState int

const (
    OPEN_PAREN TokeState = iota
    LESS_PARAMS
    FINAL_PARAM
    CLOSE_PAREN
)

func getAllowedChars(st TokeState) string {
    switch st {
    case OPEN_PAREN:
        return "("
    case LESS_PARAMS:
        return "1234567890,"
    case FINAL_PARAM:
        return "1234567890)"
    case CLOSE_PAREN:
        return ")"
    default:
        return ""
    }

}

func findParamsFollowingFunctionCall(reader *bufio.Reader, funcDef Function) (*FunctionInstance, error) {
	params := []*Param{}
    currentParam := new(Param)
    currentParam.typeStr = "int"
    state := OPEN_PAREN

	runeVal, _, readErr := reader.ReadRune()
    for readErr == nil {
        fmt.Println(currentParam)
        allowed := getAllowedChars(state)
        fmt.Println(allowed)
        fmt.Printf("evaluating %#U\n", runeVal)
        malformed := true
        for _, v := range allowed {
            if runeVal == rune(v) {
                malformed = false
            }
        }

        if malformed {
            return nil, errors.New("malformed")
        }

        fmt.Println("entering switch")

        switch state {
        case OPEN_PAREN:
            fmt.Println("( found")
            if funcDef.paramCount > 0 {
                state = LESS_PARAMS
            } else {
                state = CLOSE_PAREN
            }

        case LESS_PARAMS, FINAL_PARAM:

            if state == FINAL_PARAM && runeVal == rune(")"[0]) {
                funcCall := FunctionInstance{function: funcDef, params: params}
                return &funcCall, nil
            }

            if currentParam.intVal == "" && runeVal == rune(","[0]) {
                fmt.Println("unexpected ,")
                return nil, errors.New("malformed")
            } else if runeVal == rune(","[0]) {
                fmt.Println(", found")
                if len(params) == funcDef.paramCount - 1 {
                    state = FINAL_PARAM
                }
                currentParam = new(Param)
                currentParam.typeStr = "int"
            } else if currentParam.intVal == "" {
                currentParam.intVal += string(runeVal)
                params = append(params, currentParam)
            } else {
                currentParam.intVal += string(runeVal)
            }

        case CLOSE_PAREN:
            funcCall := FunctionInstance{function: funcDef, params: params}
			return &funcCall, nil
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
		if runeVal == rune(funcDef.keyword[functionCharsFound]) {
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
