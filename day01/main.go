package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func getInputData() ([]int, []int) {
	fileContent, err := os.ReadFile("example_input.txt")
	if err != nil {
		panic(err)
	}
	leftNumbers := make([]int, 0)
	rightNumbers := make([]int, 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		leftRight := slices.DeleteFunc(
			strings.Split(line, " "),
			func(a string) bool {
				return a == ""
			},
		)
		left, err := strconv.ParseInt(leftRight[0], 10, 16)
		if err != nil {
			panic(err)
		}
		leftNumbers = append(leftNumbers, int(left))
		right, err := strconv.ParseInt(leftRight[1], 10, 16)
		if err != nil {
			panic(err)
		}
		rightNumbers = append(rightNumbers, int(right))
	}
	return leftNumbers, rightNumbers
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	first, second := getInputData()
	slices.Sort(first)
	slices.Sort(second)
	if len(first) != len(second) {
		panic("lenght of lists is different")
	}
	difference := 0
	for i := 0; i < len(first); i++ {
		difference += absInt(first[i] - second[i])
	}
	fmt.Printf("%d\n", difference)
}
