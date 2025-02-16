package main

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	X, Y int
}

var numpad = [][]string{
	{"7", "8", "9"},
	{"4", "5", "6"},
	{"1", "2", "3"},
	{"f", "0", "A"},
}

var robotPad = [][]string{
	{"f", "^", "A"},
	{"<", "v", ">"},
}

func convertToNumpadMap(numpad [][]string) map[string]Point {
	m := make(map[string]Point)
	for i, row := range numpad {
		for j, col := range row {
			m[col] = Point{j, i}
		}
	}
	return m
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func addXPath(pathString string, diff Point) string {
	if diff.X > 0 {
		return strings.Repeat(">", abs(diff.X)) + pathString
	} else if diff.X < 0 {
		return strings.Repeat("<", abs(diff.X)) + pathString
	}
	return pathString
}

func addYPath(pathString string, diff Point) string {
	if diff.Y > 0 {
		return strings.Repeat("v", abs(diff.Y)) + pathString
	} else if diff.Y < 0 {
		return strings.Repeat("^", abs(diff.Y)) + pathString
	}
	return pathString
}

func getPathSimple(numpadMap map[string]Point, startPoint string, endPoint string, forbidenPoint string) []string {
	startPointPoint := numpadMap[startPoint]
	endPointPoint := numpadMap[endPoint]
	forbidenPointPoint := numpadMap[forbidenPoint]

	diff := Point{endPointPoint.X - startPointPoint.X, endPointPoint.Y - startPointPoint.Y}
	allPaths := []string{}
	if forbidenPointPoint.X != startPointPoint.X || forbidenPointPoint.Y != endPointPoint.Y {
		pathString := ""
		pathString = addXPath(pathString, diff)
		pathString = addYPath(pathString, diff)
		allPaths = append(allPaths, pathString)
	}
	if forbidenPointPoint.Y != startPointPoint.Y || forbidenPointPoint.X != endPointPoint.X {
		pathString := ""
		pathString = addYPath(pathString, diff)
		pathString = addXPath(pathString, diff)
		allPaths = append(allPaths, pathString)
	}

	// check if the path is valid
	validPaths := []string{}

	// check if any of the permutations is a valid path
	for _, permutation := range allPaths {
		if isValidPath(startPointPoint, endPointPoint, permutation, forbidenPointPoint) {
			validPaths = append(validPaths, permutation)
		}
	}

	if len(validPaths) == 0 {
		fmt.Printf("no valid paths found for %s to %s\n", startPoint, endPoint)
		fmt.Println(allPaths)
		fmt.Println(diff)
		fmt.Println(startPointPoint)
		fmt.Println(endPointPoint)
		fmt.Println(forbidenPointPoint)
	}

	// if there are 2 valid paths, return the one that ends with a '<'
	if len(validPaths) == 2 {
		if validPaths[0] == validPaths[1] {
			return []string{validPaths[0]}
		}
		if slices.Contains([]string{"v<", ">v", "A^"}, validPaths[0][len(validPaths[0])-2:]) {
			return []string{validPaths[1]}
		} else if slices.Contains([]string{"v<", ">v", "A^"}, validPaths[1][len(validPaths[1])-2:]) {
			return []string{validPaths[0]}
		}

	}

	return validPaths
}

func isValidPath(start Point, end Point, path string, forbidden Point) bool {
	current := start

	for _, move := range path {
		next := current
		switch move {
		case '^':
			next.Y--
		case 'v':
			next.Y++
		case '<':
			next.X--
		case '>':
			next.X++
		}

		if next == forbidden {
			return false
		}
		current = next
	}

	return current == end
}

func constructMapOfAllPaths(numpadMap map[string]Point, getPathFunc func(numpadMap map[string]Point, startPoint string, endPoint string, forbidenPoint string) []string) map[string][]string {
	allPaths := make(map[string][]string)

	// Get all pairs of keys
	keys := make([]string, 0, len(numpadMap))
	for k := range numpadMap {
		if k != "f" { // Skip forbidden point
			keys = append(keys, k)
		}
	}

	// For each pair of keys
	for _, start := range keys {
		for _, end := range keys {
			if start != end {
				// Get all valid paths between these points
				paths := getPathFunc(numpadMap, start, end, "f")
				if len(paths) > 0 {
					// Store the paths in the map using a composite key
					key := start + end
					allPaths[key] = paths
				}
			} else {
				allPaths[start+end] = []string{""}
			}
		}
	}

	return allPaths

}

func expandString(allPaths map[string][]string, stringToExpand string) []string {
	start := "A"
	allExpandedPaths := map[string]struct{}{"": {}}

	for _, char := range stringToExpand {
		nextChar := string(char)

		newExpandedPaths := make(map[string]struct{})

		for path := range allExpandedPaths {
			for _, newPath := range allPaths[start+nextChar] {
				newExpandedPaths[path+newPath+"A"] = struct{}{}
			}
		}
		allExpandedPaths = newExpandedPaths
		start = nextChar
	}

	return slices.Collect(maps.Keys(allExpandedPaths))
}

func expandPaths(allPaths map[string][]string, paths []string, printPaths bool) []string {
	expandedPaths := []string{}
	for _, path := range paths {
		if printPaths {
			fmt.Printf("expanding path: %s: %d\n", path, len(path))
		}
		expandedPaths = append(expandedPaths, expandString(allPaths, path)...)
		if printPaths {
			for _, expandedPath := range expandString(allPaths, path) {
				fmt.Printf("expanded path: %s: %d\n", expandedPath, len(expandedPath))
			}
			fmt.Println("---")
		}

	}
	return expandedPaths
}

func getShortetsPathLength(paths []string) int {
	shortestPath := paths[0]
	for _, path := range paths {
		if len(path) < len(shortestPath) {
			shortestPath = path
		}
	}
	return len(shortestPath)
}

func numericPart(code string) int {
	numericPartStr := code[:len(code)-1]
	numericPart, err := strconv.Atoi(numericPartStr)
	if err != nil {
		return 0
	}
	return numericPart
}

func expandPathsN(allPaths map[string][]string, paths []string, n int) []string {
	for i := 0; i < n; i++ {
		paths = expandPaths(allPaths, paths, false)
	}
	return paths
}

func filterPaths(allPaths map[string][]string, robotPadAllPaths map[string][]string) map[string][]string {
	filteredPaths := map[string][]string{}
	for key, paths := range allPaths {
		if len(paths) == 2 {
			// fmt.Printf("filtering paths: %s\n", key)
			// fmt.Printf("paths: %v\n", paths)
			shortestPathLength1 := getShortetsPathLength(expandPathsN(robotPadAllPaths, []string{paths[0] + "A"}, 4))
			shortestPathLength2 := getShortetsPathLength(expandPathsN(robotPadAllPaths, []string{paths[1] + "A"}, 4))
			//fmt.Printf("shortestPathLength1: %d, shortestPathLength2: %d\n", shortestPathLength1, shortestPathLength2)
			if shortestPathLength1 < shortestPathLength2 {
				filteredPaths[key] = []string{paths[0]}
			} else if shortestPathLength1 > shortestPathLength2 {
				filteredPaths[key] = []string{paths[1]}
			} else {
				fmt.Printf("equal paths: %s\n", key)
				filteredPaths[key] = []string{paths[0]}
			}
		} else {
			filteredPaths[key] = paths
		}
	}
	return filteredPaths
}

func createSeqMap(allPaths map[string][]string) (map[string]int, map[int]string) {
	seqMap := map[string]int{}
	seqMapReverse := map[int]string{}
	index := 0
	for _, paths := range allPaths {
		if _, ok := seqMap[paths[0]]; !ok {
			seqMap[paths[0]] = index
			seqMapReverse[index] = paths[0]
			index++
		}
	}
	return seqMap, seqMapReverse
}

func createPathExpandMap(allPaths map[string][]string, seqMap map[string]int) map[int][]int {
	pathExpandMap := map[int][]int{}
	for _, paths := range allPaths {
		path := paths[0]
		seq := seqMap[path]
		expandedPath := expandPaths(allPaths, []string{path + "A"}, false)[0]
		pathExpandMap[seq] = pathToSeq(seqMap, expandedPath)
	}
	return pathExpandMap
}

func createPathExpandMapN(allPaths map[string][]string, seqMap map[string]int, pathExpandMap map[int][]int, n int) map[int][]int {
	pathExpandMapN := map[int][]int{}
	for _, paths := range allPaths {
		path := paths[0]
		seq := seqMap[path]
		expandedSeq := pathExpandMap[seq]
		for i := 0; i < n-1; i++ {
			expandedSeq = expandSeq(pathExpandMap, expandedSeq)
		}
		pathExpandMapN[seq] = expandedSeq
	}
	return pathExpandMapN
}

func expandSeq(pathExpandMap map[int][]int, seq []int) []int {
	expandedSeq := []int{}
	for _, s := range seq {
		expandedSeq = append(expandedSeq, pathExpandMap[s]...)
	}
	return expandedSeq
}

func pathToSeq(seqMap map[string]int, path string) []int {
	// split between A
	parts := strings.Split(path[:len(path)-1], "A")
	seq := []int{}
	for _, part := range parts {
		seq = append(seq, seqMap[part])
	}
	return seq
}

func seqToLenMap(seqMapReverse map[int]string) map[int]int {
	lenMap := map[int]int{}
	for seqIndex, s := range seqMapReverse {
		lenMap[seqIndex] = len(s) + 1
	}
	return lenMap
}

func seqToLen(seqLenMap map[int]int, seq []int) int {
	len := 0
	for _, s := range seq {
		len += seqLenMap[s]
	}
	return len
}

func createPathLengthMap(allPaths map[string][]string, seqMap map[string]int, seqLenMap map[int]int, pathExpandMap map[int][]int) map[int]int {
	pathLengthMap := map[int]int{}
	for _, paths := range allPaths {
		path := paths[0]
		seq := seqMap[path]
		expandedSeq := pathExpandMap[seq]
		pathLengthMap[seq] = seqToLen(seqLenMap, expandedSeq)
	}
	return pathLengthMap
}

func main() {
	startTime := time.Now()
	numpadMap := convertToNumpadMap(numpad)
	robotPadMap := convertToNumpadMap(robotPad)

	// stringToExpand := "029A"
	/*stringToExpand := []string{
		"029A",
		"980A",
		"179A",
		"456A",
		"379A",
	}*/

	stringToExpand := []string{
		"286A",
		"974A",
		"189A",
		"802A",
		"805A",
	}

	numPadAllPaths := constructMapOfAllPaths(numpadMap, getPathSimple)
	robotPadAllPaths := constructMapOfAllPaths(robotPadMap, getPathSimple)
	robotPadAllPaths = filterPaths(robotPadAllPaths, robotPadAllPaths)
	numPadAllPaths = filterPaths(numPadAllPaths, robotPadAllPaths)

	seqMap, seqMapReverse := createSeqMap(robotPadAllPaths)
	seqLenMap := seqToLenMap(seqMapReverse)
	pathExpandMap := createPathExpandMap(robotPadAllPaths, seqMap)
	pathExpandMapN := createPathExpandMapN(robotPadAllPaths, seqMap, pathExpandMap, 13)
	seqLengthMapN := createPathLengthMap(robotPadAllPaths, seqMap, seqLenMap, pathExpandMapN)

	totalComplexity := 0

	for _, code := range stringToExpand {
		expanded := expandPaths(numPadAllPaths, []string{code}, false)
		expanded = expandPaths(robotPadAllPaths, expanded, false)

		seq := pathToSeq(seqMap, expanded[0])

		// expandedSeq := expandSeq(pathExpandMap, seq)
		expandedSeqAlt := expandSeq(pathExpandMap, seq)
		/*for i := 0; i < 12; i++ {
			expandedSeq = expandSeq(pathExpandMap, expandedSeq)
		}*/

		for i := 0; i < 10; i++ {
			expandedSeqAlt = expandSeq(pathExpandMap, expandedSeqAlt)
		}

		pathLen := 0
		for _, s := range expandedSeqAlt {
			pathLen += seqLengthMapN[s]
		}

		/*for i := 0; i < 17; i++ {
			expandedSeqAlt = expandSeq(pathExpandMap, expandedSeqAlt)
		}*/

		// fmt.Printf("expanded: %v\n", seqToLen(seqLenMap, expandedSeqAlt))
		numPart := numericPart(code)
		complexity := pathLen * numPart
		fmt.Printf("code: %s, complexity: %d\n", code, complexity)
		totalComplexity += complexity
	}

	fmt.Printf("time: %s\n", time.Since(startTime))
	fmt.Printf("total complexity: %d\n", totalComplexity)
}
