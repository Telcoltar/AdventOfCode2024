package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

func readData(filepath string) ([]int, error) {
	fileContent, readFileErr := os.ReadFile(filepath)
	if readFileErr != nil {
		return nil, readFileErr
	}
	numbers := make([]int, 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		parsedNumber, parsingErr := strconv.Atoi(strings.TrimSpace(line))
		if parsingErr != nil {
			return nil, parsingErr
		}
		numbers = append(numbers, parsedNumber)
	}
	return numbers, nil
}

func calculateSecretNumbers(init int) []int {
	secretNumbers := make([]int, 0)
	// modulo 16777216
	mask := (1 << 24) - 1
	multMask := (1 << 18) - 1
	secretNumbers = append(secretNumbers, init)
	number := init & mask
	for i := 0; i < 2000; i++ {
		number = number ^ ((number & multMask) << 6)
		number = number & mask
		number = number ^ (number >> 5)
		number = number & mask
		number = number ^ (number << 11)
		number = number & mask
		secretNumbers = append(secretNumbers, number)
	}
	return secretNumbers
}

func solutionPart1(numbers []int) int {
	total := 0
	for _, number := range numbers {
		secertNumbers := calculateSecretNumbers(number)
		total += secertNumbers[len(secertNumbers)-1]
	}
	return total
}

func calculateChangesAndPrices(secretNumbers []int) ([]int, []int) {
	changes := make([]int, 0)
	prices := make([]int, 0)
	prevLastDigit := secretNumbers[0] % 10
	for i := 1; i < len(secretNumbers); i++ {
		lastDigit := secretNumbers[i] % 10
		prices = append(prices, lastDigit)
		changes = append(changes, lastDigit-prevLastDigit)
		prevLastDigit = lastDigit
	}
	return changes, prices
}

func solutionPart2(numbers []int) int {
	bannanasByKey := make(map[[4]int]int)
	for _, number := range numbers {
		secretNumbers := calculateSecretNumbers(number)
		changes, prices := calculateChangesAndPrices(secretNumbers)
		seenKeys := make(map[[4]int]bool)
		for i := 3; i < len(changes); i++ {
			key := [4]int{changes[i-3], changes[i-2], changes[i-1], changes[i]}
			if _, ok := seenKeys[key]; ok {
				continue
			}
			seenKeys[key] = true
			bannanasByKey[key] += prices[i]
		}
	}
	// find squence of four changes maximize price
	maxBannanas := slices.Max(slices.Collect(maps.Values(bannanasByKey)))

	fmt.Printf("Max bannanas: %d\n", maxBannanas)
	return maxBannanas
}

func main() {
	numbers, readDataErr := readData("input.txt")
	if readDataErr != nil {
		panic(readDataErr)
	}
	println(solutionPart1(numbers))
	println(solutionPart2(numbers))
}
