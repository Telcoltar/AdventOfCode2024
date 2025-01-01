package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func getInputData(inputDataPath string) ([]int, []int) {
	fileContent, err := os.ReadFile(inputDataPath)
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
		left, err := strconv.ParseInt(leftRight[0], 10, 32)
		if err != nil {
			panic(err)
		}
		leftNumbers = append(leftNumbers, int(left))
		right, err := strconv.ParseInt(leftRight[1], 10, 32)
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

func part1() {
	leftList, rightList := getInputData("inputData.txt")
	slices.Sort(leftList)
	slices.Sort(rightList)
	if len(leftList) != len(rightList) {
		panic("length of lists is different")
	}
	difference := 0
	for i := 0; i < len(leftList); i++ {
		difference += absInt(leftList[i] - rightList[i])
	}
	fmt.Printf("%d\n", difference)
}

func part2() {
	leftList, rightList := getInputData("inputData.txt")
	count := make(map[int]int)
	for i := 0; i < len(rightList); i++ {
		count[rightList[i]]++
	}
	total := 0
	for i := 0; i < len(leftList); i++ {
		total += count[leftList[i]] * leftList[i]
	}
	fmt.Printf("%d\n", total)
}

func main() {
	part1()
	part2()
}
