package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Point struct {
	x, y int
}

func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y}
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

func cheatNeighbours(g Grid[bool], p Point) []Point {
	neighbours := make([]Point, 0)
	for _, dir := range DIRS {
		n := p.Add(dir)
		n2 := n.Add(dir)
		if g.IsInGrid(n2) && !g.At(n) && g.At(n2) {
			neighbours = append(neighbours, n2)
		}
	}
	return neighbours
}

func fillScores(grid Grid[bool], start Point) Grid[int] {
	queue := []Point{start}
	scores := NewGrid[int](grid.Height())
	scores.Fill(math.MaxInt)
	scores.Set(start, 0)
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
	return scores
}

func readData(filename string) (Grid[bool], Point, Point) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(fileContent), "\n")
	grid := NewGrid[bool](len(lines))
	grid.Fill(true)
	var start Point
	var end Point
	for y, line := range lines {
		for x, cell := range []rune(strings.TrimSpace(line)) {
			if cell == '#' {
				grid[y][x] = false
			} else if cell == 'S' {
				start.x = x
				start.y = y
			} else if cell == 'E' {
				end.x = x
				end.y = y
			}
		}
	}
	return grid, start, end
}

func printScores(grid Grid[bool], scores Grid[int]) {
	sb := strings.Builder{}
	for y, row := range scores {
		for x, cell := range row {
			if grid[y][x] {
				sb.WriteString(fmt.Sprint(cell))
			} else {
				sb.WriteString("#")
			}
			sb.WriteString("\t")
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}

func part1(filename string) {
	grid, start, _ := readData(filename)
	scores := fillScores(grid, start)
	cheats := make(map[int]int)
	for y, row := range grid {
		for x, cell := range row {
			if cell {
				neighbours := cheatNeighbours(grid, Point{x, y})
				score := scores.At(Point{x, y})
				for _, neighbour := range neighbours {
					cheatDiff := scores.At(neighbour) - score - 2
					if cheatDiff > 0 {
						cheats[cheatDiff]++
					}
				}
			}
		}
	}
	sum := 0
	for cheatTime, count := range cheats {
		if cheatTime >= 100 {
			sum += count
		}
	}
	fmt.Println(sum)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func distance(p1, p2 Point) int {
	return absInt(p2.x-p1.x) + absInt(p2.y-p1.y)
}

func part2(filename string) {
	grid, start, _ := readData(filename)
	scores := fillScores(grid, start)
	maxScore := 0
	for y, row := range scores {
		for x, cell := range row {
			if grid[y][x] {
				maxScore = maxInt(maxScore, cell)
			}
		}
	}
	steps := make([]Point, maxScore+1)
	for y, row := range scores {
		for x, cell := range row {
			if grid[y][x] {
				steps[cell] = Point{x, y}
			}
		}
	}
	cheats := make(map[int]int)
	cutoff := 100
	for i := 0; i < maxScore+1-cutoff; i++ {
		current := steps[i]
		for j := i + cutoff; j < maxScore+1; j++ {
			next := steps[j]
			if distance(next, current) <= 20 {
				cheats[j-i-distance(next, current)]++
			}
		}
	}
	sum := 0
	for cheatTime, count := range cheats {
		if cheatTime >= cutoff {
			sum += count
		}
	}
	fmt.Println(sum)
}

func main() {
	filename := "input.txt"
	part1(filename)
	part2(filename)
}
