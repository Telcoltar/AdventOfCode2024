package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func readData(filename string) [][]int {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	numbers := make([][]int, 0)
	for _, line := range strings.Split(string(content), "\n") {
		lineNumbers := make([]int, 0)
		for _, number := range strings.Split(line, " ") {
			if parsedNumber, err := strconv.Atoi(number); err == nil {
				lineNumbers = append(lineNumbers, parsedNumber)
			} else {
				panic(err)
			}
		}
		numbers = append(numbers, lineNumbers)
	}
	return numbers
}

func validReport(numbers []int) int {
	diff := numbers[0] - numbers[1]
	if diff == 0 || diff > 3 || diff < -3 {
		return 1
	}
	descending := diff > 0
	for i := 1; i < len(numbers); i++ {
		diff := numbers[i-1] - numbers[i]
		if diff == 0 || diff > 3 || diff < -3 {
			return i
		}

		if descending {
			if diff < 0 {
				return i
			}
		} else {
			if diff > 0 {
				return i
			}
		}
	}
	return -1
}

func part1(reports [][]int) int {
	count := 0
	for _, report := range reports {
		if validReport(report) == -1 {
			count++
		}
	}
	return count
}

func removeElement(numbers []int, i int) []int {
	if i == len(numbers)-1 {
		return slices.Clone(numbers[:i])
	}
	if i == 0 {
		return slices.Clone(numbers[1:])
	}

	return append(slices.Clone(numbers[:i]), numbers[i+1:]...)
}

func part2(reports [][]int) int {
	count := 0
	for idx, report := range reports {
		faultyIndex := validReport(report)
		if faultyIndex == -1 {
			count++
		} else {
			modifiedReport := removeElement(report, faultyIndex)
			// fmt.Printf("%d: %d -> %d: %d\n", idx, faultyIndex, modifiedReport, validReport(modifiedReport))
			if validReport(modifiedReport) == -1 {
				count++
				continue
			}
			modifiedReport = removeElement(report, faultyIndex-1)
			if validReport(modifiedReport) == -1 {
				fmt.Printf("%d: %d -> %#v\n", idx, faultyIndex, report)
				count++
				continue
			}
		}
	}
	return count
}

func main() {
	reports := readData("input.txt")
	fmt.Printf("Part 1: %d\n", part1(reports))
	fmt.Printf("Part 2: %d\n", part2(reports))
}
