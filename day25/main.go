package main

import (
	"fmt"
	"os"
	"strings"
)

func runesEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func mirrorGrid(grid [][]rune) [][]rune {
	mirroredGrid := make([][]rune, len(grid))
	for i := range (len(grid) / 2) + 1 {
		mirroredGrid[i] = grid[len(grid)-i-1]
		mirroredGrid[len(grid)-i-1] = grid[i]
	}
	return mirroredGrid
}

func printKeyLock(representation [][]rune) {
	for _, line := range representation {
		fmt.Println(string(line))
	}
	fmt.Println()
}

func parseKeyLock(repr [][]rune) ([5]int, bool) {
	// Parse key lock representation
	var key [5]int = [5]int{}
	isLock := runesEqual(repr[0], []rune("#####"))
	if !isLock {
		repr = mirrorGrid(repr)
	}
	// printKeyLock(repr)
	for i := 0; i < 5; i++ {
		currentIndex := 0
		currentField := repr[currentIndex+1][i]
		for currentField == '#' {
			currentIndex++
			currentField = repr[currentIndex+1][i]
		}
		key[i] = currentIndex
	}

	return key, isLock
}

func readData(filepath string) ([][5]int, [][5]int, error) {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, nil, err
	}

	keys := make([][5]int, 0)
	locks := make([][5]int, 0)
	for _, strRepr := range strings.Split(string(fileContent), "\n\n") {
		// Parse key lock representation
		repr := make([][]rune, 0)
		for _, line := range strings.Split(strRepr, "\n") {
			repr = append(repr, []rune(line))
		}
		key, isLock := parseKeyLock(repr)
		if isLock {
			locks = append(locks, key)
		} else {
			keys = append(keys, key)
		}
	}

	return keys, locks, nil
}

func solutionPart1(keys, locks [][5]int) int {
	hits := 0
	// try any key on any lock
	for _, key := range keys {
		for _, lock := range locks {
			fits := true
			for i := 0; i < 5; i++ {
				if key[i]+lock[i] > 5 {
					fits = false
					break
				}
			}
			if fits {
				hits++
			}
		}
	}
	return hits
}

func main() {
	keys, locks, err := readData("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Part 1 Solution:", solutionPart1(keys, locks))
}
