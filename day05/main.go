package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readData(filename string) ([][]int, [][]int) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	idx := 0
	orderRules := make([][]int, 0)
	for {
		line := strings.TrimSpace(lines[idx])
		if line == "" {
			break
		}
		split := strings.Split(line, "|")
		left, err := strconv.Atoi(split[0])
		if err != nil {
			panic(err)
		}
		right, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}
		orderRules = append(orderRules, []int{left, right})
		idx++
	}
	idx++
	pagesToProduce := make([][]int, 0)
	for idx < len(lines) {
		line := strings.TrimSpace(lines[idx])
		split := strings.Split(line, ",")
		pages := make([]int, len(split))
		for i, page := range split {
			parsed, err := strconv.Atoi(page)
			if err != nil {
				panic(err)
			}
			pages[i] = parsed
		}
		pagesToProduce = append(pagesToProduce, pages)
		idx++
	}
	return orderRules, pagesToProduce
}

func checkRules(rules [][]int, page map[int]int) bool {
	for _, rule := range rules {
		firstIdx, ok := page[rule[0]]
		if !ok {
			continue
		}
		secondIdx, ok := page[rule[1]]
		if !ok {
			continue
		}
		if firstIdx > secondIdx {
			return false
		}
	}
	return true
}

func part1(filename string) {
	rules, pagesToOrder := readData(filename)
	sum := 0
	for _, pages := range pagesToOrder {
		pagesMap := make(map[int]int)
		for i, page := range pages {
			pagesMap[page] = i
		}
		correct := checkRules(rules, pagesMap)
		if correct {
			middle := len(pages) / 2
			sum += pages[middle]
		}
	}
	fmt.Println(sum)
}

type Pair struct {
	left, right int
}

func buildRuleHashMap(rules [][]int) map[Pair]struct{} {
	hashMap := make(map[Pair]struct{})
	for _, rule := range rules {
		hashMap[Pair{rule[0], rule[1]}] = struct{}{}
	}
	return hashMap
}

func naiveSort(pages []int, ruleHashMap map[Pair]struct{}) {
	for i := 0; i < len(pages); i++ {
		for j := 0; j < len(pages)-1; j++ {
			if _, ok := ruleHashMap[Pair{pages[j], pages[j+1]}]; !ok {
				tmp := pages[j]
				pages[j] = pages[j+1]
				pages[j+1] = tmp
			}
		}
	}
}

func part2(filename string) {
	rules, pagesToOrder := readData(filename)
	sum := 0
	for _, pages := range pagesToOrder {
		pagesMap := make(map[int]int)
		for i, page := range pages {
			pagesMap[page] = i
		}
		if !checkRules(rules, pagesMap) {
			hashMap := buildRuleHashMap(rules)
			naiveSort(pages, hashMap)
			middle := len(pages) / 2
			sum += pages[middle]
		}
	}
	fmt.Println(sum)
}

func main() {
	part2("input.txt")
}
