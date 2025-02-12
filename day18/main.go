package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y}
}

func MustParse(s string) Point {
	s = strings.TrimSpace(s)
	p := Point{}
	split := strings.Split(s, ",")
	p.x, _ = strconv.Atoi(split[0])
	p.y, _ = strconv.Atoi(split[1])
	return p
}

func readData(filename string) []Point {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	output := make([]Point, 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		output = append(output, MustParse(line))
	}
	return output
}

type Grid[T any] [][]T

func NewGrid[T any](size int) Grid[T] {
	grid := make([][]T, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]T, size)
	}
	return grid
}

func (g Grid[T]) Fill(value T) {
	for y, row := range g {
		for x := range row {
			g[y][x] = value
		}
	}
}

func (g Grid[T]) At(p Point) T {
	return g[p.y][p.x]
}

func (g Grid[T]) Set(p Point, v T) {
	g[p.y][p.x] = v
}

func (g Grid[T]) Width() int {
	return len(g[0])
}

func (g Grid[T]) Height() int {
	return len(g[0])
}

func (g Grid[T]) String() string {
	sb := strings.Builder{}
	for _, row := range g {
		for _, cell := range row {
			sb.WriteString(" ")
			sb.WriteString(fmt.Sprint(cell))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g Grid[T]) IsInGrid(p Point) bool {
	if p.x < 0 || p.x >= g.Width() {
		return false
	}
	if p.y < 0 || p.y >= g.Height() {
		return false
	}
	return true
}

var DIRS = []Point{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

func validNeighbours(g Grid[bool], p Point) []Point {
	neighbours := make([]Point, 0)
	for _, dir := range DIRS {
		n := p.Add(dir)
		if g.IsInGrid(n) && g.At(n) {
			neighbours = append(neighbours, n)
		}
	}
	return neighbours
}

func findPath(grid Grid[bool], size int) int {
	queue := []Point{{0, 0}}
	scores := NewGrid[int](size)
	scores.Fill(math.MaxInt)
	scores.Set(Point{0, 0}, 0)
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		score := scores.At(p)
		for _, neighbour := range validNeighbours(grid, p) {
			if scores.At(neighbour) > score+1 {
				scores.Set(neighbour, score+1)
				queue = append(queue, neighbour)
			}
		}
	}
	return scores.At(Point{size - 1, size - 1})
}

func part1(filename string, size int, byteLimit int) {
	bytesCoordinates := readData(filename)
	grid := NewGrid[bool](size)
	grid.Fill(true)
	for _, b := range bytesCoordinates[:byteLimit] {
		grid.Set(b, false)
	}
	fmt.Println(findPath(grid, size))
}

func part2(filename string, size int, byteLimit int) {
	bytesCoordinates := readData(filename)
	grid := NewGrid[bool](size)
	grid.Fill(true)
	for _, b := range bytesCoordinates[:byteLimit] {
		grid.Set(b, false)
	}
	index := byteLimit
	for findPath(grid, size) < math.MaxInt {
		grid.Set(bytesCoordinates[index], false)
		index++
	}
	fmt.Println(bytesCoordinates[index-1])
}

func main() {
	part1("input.txt", 71, 1024)
	part2("input.txt", 71, 1024)
}
