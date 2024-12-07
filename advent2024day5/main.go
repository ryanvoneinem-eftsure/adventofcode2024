package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type PageRules struct {
	Page        string
	SucceededBy []string
	PreceededBy []string
}

func (pr *PageRules) ToString() string {
	return pr.Page + ": SucceededBy[" + strings.Join(pr.SucceededBy, ",") + "] PreceededBy[" + strings.Join(pr.PreceededBy, ",") + "]"
}

func main() {

	filename := "./input.txt"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	file, scanner := createFileScanner(filename)
	defer file.Close()

	rules := collectOrderingRules(scanner)
	fmt.Println("Rules:")
	for _, r := range rules {
		fmt.Println(r.ToString())
	}

	// Part 1
	// updates := findOrderedUpdates(scanner, rules)

	// Part 2
	updates := findUnorderedUpdates(scanner, rules)
	updates = sortUpdates(updates, rules)

	sumMedian := 0
	for _, u := range updates {
		fmt.Println(u)

		i := (len(u) / 2)
		m, _ := strconv.Atoi(u[i])
		sumMedian += m
	}

	fmt.Println("sum of median page:", sumMedian)

}

func sortUpdates(updates [][]string, rules map[string]*PageRules) [][]string {

	rulesUpdates := [][]*PageRules{}
	for _, u := range updates {
		rulesUpdate := []*PageRules{}
		for _, p := range u {
			rulesUpdate = append(rulesUpdate, rules[p])
		}
		rulesUpdates = append(rulesUpdates, rulesUpdate)
	}
	ordered := [][]string{}
	for _, rulesUpdate := range rulesUpdates {
		slices.SortFunc(rulesUpdate, func(a *PageRules, b *PageRules) int {
			if len(a.PreceededBy) == 0 && len(b.PreceededBy) == 0 {
				return 0
			} else if len(a.PreceededBy) == 0 {
				return -1
			} else if len(b.PreceededBy) == 0 {
				return 1
			}

			if slices.Contains(a.SucceededBy, b.Page) || slices.Contains(b.PreceededBy, a.Page) {
				return -1
			} else if slices.Contains(a.PreceededBy, b.Page) || slices.Contains(b.SucceededBy, a.Page) {
				return 1
			}
			return 0
		})
		ou := []string{}
		for _, ru := range rulesUpdate {
			ou = append(ou, ru.Page)
		}
		ordered = append(ordered, ou)
	}

	return ordered

}

func findUnorderedUpdates(scanner *bufio.Scanner, rules map[string]*PageRules) [][]string {
	ordered := [][]string{}

	for scanner.Scan() {
		line := scanner.Text()
		update := strings.Split(line, ",")
		applicable := []*PageRules{}
		isOrdered := true
		for _, p := range update {
			applicable = append(applicable, rules[p])
		}
		for _, r := range applicable {
			pIdx := slices.Index(update, r.Page)
			for i := 0; i < pIdx; i++ {
				if slices.Contains(r.SucceededBy, update[i]) {
					isOrdered = false
					break
				}
			}
			for i := pIdx + 1; i < len(update); i++ {
				if slices.Contains(r.PreceededBy, update[i]) {
					isOrdered = false
					break
				}
			}
		}
		if !isOrdered {
			ordered = append(ordered, update)
		}

	}

	return ordered
}

func findOrderedUpdates(scanner *bufio.Scanner, rules map[string]*PageRules) [][]string {
	ordered := [][]string{}

	for scanner.Scan() {
		line := scanner.Text()
		update := strings.Split(line, ",")
		applicable := []*PageRules{}
		isOrdered := true
		for _, p := range update {
			applicable = append(applicable, rules[p])
		}
		for _, r := range applicable {
			pIdx := slices.Index(update, r.Page)
			for i := 0; i < pIdx; i++ {
				if slices.Contains(r.SucceededBy, update[i]) {
					isOrdered = false
					break
				}
			}
			for i := pIdx + 1; i < len(update); i++ {
				if slices.Contains(r.PreceededBy, update[i]) {
					isOrdered = false
					break
				}
			}
		}
		if isOrdered {
			ordered = append(ordered, update)
		}

	}

	return ordered
}

func collectOrderingRules(scanner *bufio.Scanner) map[string]*PageRules {
	rules := map[string]*PageRules{}

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			return rules
		}

		pages := strings.Split(line, "|")

		ruleL := rules[pages[0]]
		ruleR := rules[pages[1]]

		if ruleL == nil {
			ruleL = &PageRules{Page: pages[0]}
			rules[pages[0]] = ruleL
		}
		ruleL.SucceededBy = append(ruleL.SucceededBy, pages[1])

		if ruleR == nil {
			ruleR = &PageRules{Page: pages[1]}
			rules[pages[1]] = ruleR
		}
		ruleR.PreceededBy = append(ruleR.PreceededBy, pages[0])

	}

	return rules

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
