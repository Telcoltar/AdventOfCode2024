package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type grid = [][]int

func readData(filename string) grid {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	hikingMap := make(grid, 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		row := make([]int, 0)
		for _, char := range line {
			var parsed int
			if char == '.' {
				parsed = -2
			} else {
				parsed, err = strconv.Atoi(string(char))
				if err != nil {
					panic(err)
				}
			}
			row = append(row, parsed)
		}
		hikingMap = append(hikingMap, row)
	}
	return hikingMap
}

type Point struct {
	x, y int
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getNeighbours(p *Point, hikingMap grid) []Point {
	currentElevation := hikingMap[p.y][p.x]
	neighbours := make([]Point, 0)
	if p.y-1 >= 0 && hikingMap[p.y-1][p.x]-currentElevation == -1 {
		neighbours = append(neighbours, Point{p.x, p.y - 1})
	}
	if p.x-1 >= 0 && hikingMap[p.y][p.x-1]-currentElevation == -1 {
		neighbours = append(neighbours, Point{p.x - 1, p.y})
	}
	if p.y+1 < len(hikingMap) && hikingMap[p.y+1][p.x]-currentElevation == -1 {
		neighbours = append(neighbours, Point{p.x, p.y + 1})
	}
	if p.x+1 < len(hikingMap[0]) && hikingMap[p.y][p.x+1]-currentElevation == -1 {
		neighbours = append(neighbours, Point{p.x + 1, p.y})
	}
	return neighbours
}

func fillGrid(hikingMap grid, start Point, valueMap grid) {
	queue := make([]Point, 0)
	queue = append(queue, start)
	visited := make(map[Point]struct{})
	var current Point
	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]
		if _, ok := visited[current]; ok {
			continue
		}
		valueMap[current.y][current.x]++
		visited[current] = struct{}{}
		queue = append(queue, getNeighbours(&current, hikingMap)...)
	}
}

func fillGrid2(hikingMap grid, start Point, valueMap grid) {
	queue := make([]Point, 0)
	queue = append(queue, start)
	var current Point
	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]
		valueMap[current.y][current.x]++
		queue = append(queue, getNeighbours(&current, hikingMap)...)
	}
}

func printGrid(hikingMap grid) {
	for _, row := range hikingMap {
		for _, cell := range row {
			fmt.Print(cell, " ")
		}
		fmt.Println()
	}
}

func cellsWithValue(hikingMap grid, value int) []Point {
	cells := make([]Point, 0)
	for idy, row := range hikingMap {
		for idx, cell := range row {
			if cell == value {
				cells = append(cells, Point{idx, idy})
			}
		}
	}
	return cells
}

func solution(filename string, fillFunc func(grid, Point, grid)) {
	hikingMap := readData(filename)
	nines := cellsWithValue(hikingMap, 9)
	valueMap := make(grid, len(hikingMap))
	for i := 0; i < len(valueMap); i++ {
		valueMap[i] = make([]int, len(hikingMap[0]))
	}
	for _, nine := range nines {
		fillFunc(hikingMap, nine, valueMap)
	}

	sum := 0
	for _, zero := range cellsWithValue(hikingMap, 0) {
		sum += valueMap[zero.y][zero.x]
	}
	fmt.Println(sum)
}

func part2(filename string) {
	solution(filename, fillGrid2)
}

func part1(filename string) {
	solution(filename, fillGrid)
}

func main() {
	part1("input.txt")
	part2("input.txt")
}
