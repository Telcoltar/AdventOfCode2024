package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type equation struct {
	result       int
	inputNumbers []int
}

func readData(fileName string) []equation {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	equations := make([]equation, 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		splitLine := strings.Split(strings.TrimSpace(line), ":")
		result, err := strconv.Atoi(strings.TrimSpace(splitLine[0]))
		if err != nil {
			panic(err)
		}
		inputNumbers := strings.Split(strings.TrimSpace(splitLine[1]), " ")
		equationNumbers := make([]int, len(inputNumbers))
		for _, inputNumber := range inputNumbers {
			parsedNumber, err := strconv.Atoi(strings.TrimSpace(inputNumber))
			if err != nil {
				panic(err)
			}
			equationNumbers = append(equationNumbers, parsedNumber)
		}
		equations = append(equations, equation{result, equationNumbers})
	}
	return equations
}

func testEquation(eq equation, extraOP bool) bool {
	currentPossibleValues := map[int]struct{}{eq.inputNumbers[0]: {}}
	for _, number := range eq.inputNumbers[1:] {
		nextPossibleValues := make(map[int]struct{})
		for currentPossibility := range currentPossibleValues {
			if currentPossibility > eq.result {
				continue
			}
			nextPossibleValues[currentPossibility*number] = struct{}{}
			nextPossibleValues[currentPossibility+number] = struct{}{}
			if extraOP {
				opResultStr := fmt.Sprintf("%d%d", currentPossibility, number)
				opResult, err := strconv.Atoi(opResultStr)
				if err != nil {
					panic(err)
				}
				nextPossibleValues[opResult] = struct{}{}
			}
		}
		currentPossibleValues = nextPossibleValues
	}
	for possibility := range currentPossibleValues {
		if possibility == eq.result {
			return true
		}
	}
	return false
}

func part1(equations []equation) {
	sum := 0
	for _, eq := range equations {
		if testEquation(eq, false) {
			sum += eq.result
		}
	}
	fmt.Println(sum)
}

func part2(equations []equation) {
	sum := 0
	for _, eq := range equations {
		if testEquation(eq, true) {
			sum += eq.result
		}
	}
	fmt.Println(sum)
}

func main() {
	part1(readData("input.txt"))
	part2(readData("input.txt"))
}
