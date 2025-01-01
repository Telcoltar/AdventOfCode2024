package main

import (
	"fmt"
	"os"
	"strings"
)

func readData(filename string) [][]rune {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	data := make([][]rune, 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		data = append(data, []rune(strings.TrimSpace(line)))
	}
	return padGrid(data)
}

func padGrid(grid [][]rune) [][]rune {
	paddedGrid := make([][]rune, len(grid)+2)
	width := len(grid[0]) + 2
	paddedGrid[0] = make([]rune, width)
	for i := 0; i < width; i++ {
		paddedGrid[0][i] = 'b'
	}
	for i, row := range grid {
		paddedGrid[i+1] = append(paddedGrid[i+1], 'b')
		paddedGrid[i+1] = append(paddedGrid[i+1], row...)
		paddedGrid[i+1] = append(paddedGrid[i+1], 'b')
	}
	paddedGrid[len(grid)+1] = make([]rune, width)
	for i := 0; i < width; i++ {
		paddedGrid[len(grid)+1][i] = 'b'
	}
	return paddedGrid
}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}

func findGuard(grid [][]rune) Point {
	for idy, row := range grid {
		for idx, cell := range row {
			if cell == '^' {
				return Point{idx, idy}
			}
		}
	}
	return Point{-1, -1}
}

type Direction struct {
	x, y int
}

func (d *Direction) turnRight() {
	tmp := d.x
	d.x = -d.y
	d.y = tmp
}

type Point = Direction

func (p *Point) Add(d Direction) {
	p.x += d.x
	p.y += d.y
}

func (p *Point) Sub(d Direction) {
	p.x -= d.x
	p.y -= d.y
}

func nextCell(grid [][]rune, point Point, dir Direction) rune {
	return grid[point.y+dir.y][point.x+dir.x]
}

func takeStep(grid [][]rune, start Point, direction Direction,
	visitedCells map[Point]struct{}) Point {
	currentCell := start
	for nextCell(grid, currentCell, direction) == '.' {
		currentCell.Add(direction)
		visitedCells[currentCell] = struct{}{}
	}
	return currentCell
}

func part1(filename string) {
	grid := readData(filename)
	dir := Direction{0, -1}
	currentCell := findGuard(grid)
	grid[currentCell.y][currentCell.x] = '.'
	visitedCells := make(map[Point]struct{})
	visitedCells[currentCell] = struct{}{}
	for {
		currentCell = takeStep(grid, currentCell, dir, visitedCells)
		if nextCell(grid, dir, currentCell) == 'b' {
			break
		}
		dir.turnRight()
	}
	fmt.Printf("%d\n", len(visitedCells))
}

type PointDir struct {
	p   Point
	dir Direction
}

func checkForLoops(grid [][]rune, start Point) bool {
	turningPoints := make(map[PointDir]struct{})
	currentDir := Direction{0, -1}
	// turningPoints[PointDir{start, currentDir}] = struct{}{}
	currentCell := start
	for {
		for nextCell(grid, currentCell, currentDir) == '.' {
			currentCell.Add(currentDir)
		}
		if nextCell(grid, currentCell, currentDir) == 'b' {
			return false
		}
		if _, ok := turningPoints[PointDir{currentCell, currentDir}]; ok {
			return true
		}
		turningPoints[PointDir{currentCell, currentDir}] = struct{}{}
		currentDir.turnRight()
	}
}

func part2(filename string) {
	grid := readData(filename)
	start := findGuard(grid)
	grid[start.y][start.x] = '.'
	possibilities := 0
	for i := 1; i < len(grid)-1; i++ {
		for j := 1; j < len(grid[i])-1; j++ {
			if i == start.y && j == start.x {
				continue
			}
			if grid[i][j] != '.' {
				continue
			}
			grid[i][j] = '#'
			if checkForLoops(grid, start) {
				// fmt.Printf("x: %d, y %d\n", j, i)
				possibilities++
			}
			grid[i][j] = '.'
		}
	}
	fmt.Printf("%d\n", possibilities)
}

func main() {
	// part1("input.txt")
	part2("input.txt")
}
