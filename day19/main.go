package main

import (
	"fmt"
	"os"
	"strings"
)

func readData(filename string) ([]string, []string) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(fileContent), "\n")
	split := strings.Split(strings.TrimSpace(lines[0]), ",")
	var towels []string
	for _, towel := range split {
		towels = append(towels, strings.TrimSpace(towel))
	}
	var designs []string
	for _, s := range lines[2:] {
		designs = append(designs, strings.TrimSpace(s))
	}
	return designs, towels
}

func getPossibilities(designs []string, towels []string) []int {
	possibilities := make([]int, len(designs))
	for j, design := range designs {
		count := make([]int, len(design)+1)
		count[0] = 1
		for i := range design {
			for _, towel := range towels {
				if len(towel)+i <= len(design) {
					if towel == design[i:i+len(towel)] {
						count[i+len(towel)] += count[i]
					}
				}
			}
		}
		possibilities[j] = count[len(count)-1]
	}
	return possibilities
}

func part1(filename string) {
	designs, towels := readData(filename)
	impossible := 0
	for _, possibility := range getPossibilities(designs, towels) {
		if possibility > 0 {
			impossible++
		}
	}
	fmt.Printf("Part 1: %d\n", impossible)
}

func part2(filename string) {
	designs, towels := readData(filename)
	possibilities := 0
	for _, possibility := range getPossibilities(designs, towels) {
		possibilities += possibility
	}
	fmt.Printf("Part 2: %d\n", possibilities)
}

func main() {
	part1("input.txt")
	part2("input.txt")
}
