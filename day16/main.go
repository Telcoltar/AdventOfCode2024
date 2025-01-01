package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Grid[T any] [][]T

type Point struct {
	x, y int
}

type FieldScore map[Point]int

func (g Grid[T]) getValue(p *Point) T {
	return g[p.y][p.x]
}

func (g Grid[T]) setValue(p *Point, v T) {
	g[p.y][p.x] = v
}

func (g Grid[T]) height() int {
	return len(g)
}

func (g Grid[T]) width() int {
	return len(g[0])
}

func addPoints(p1 *Point, p2 *Point) Point {
	return Point{p1.x + p2.x, p1.y + p2.y}
}

// p1 - p2
func subPoint(p1 *Point, p2 *Point) Point {
	return Point{p1.x - p2.x, p1.y - p2.y}
}

func turnRight(p *Point) Point {
	return Point{-p.y, p.x}
}

func turnLeft(p *Point) Point {
	return Point{p.y, -p.x}
}

func getNegativePoint(p *Point) Point {
	return Point{-p.x, -p.y}
}

var (
	NORTH = Point{0, -1}
	EAST  = Point{1, 0}
	SOUTH = Point{0, 1}
	WEST  = Point{-1, 0}
	DIRS  = []Point{
		NORTH, EAST, SOUTH, WEST,
	}
)

func getNeighbours(p *Point, g Grid[rune]) []Point {
	neighbours := make([]Point, 0)
	for _, dir := range DIRS {
		possibleNeighbour := addPoints(p, &dir)
		if g.getValue(&possibleNeighbour) != '#' {
			neighbours = append(neighbours, possibleNeighbour)
		}
	}
	return neighbours
}

func readData(filename string) Grid[rune] {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	grid := make(Grid[rune], 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		grid = append(grid, []rune(strings.TrimSpace(line)))
	}
	return grid
}

func findSymbol(grid Grid[rune], symbol rune) Point {
	for y, row := range grid {
		for x, cell := range row {
			if cell == symbol {
				return Point{x, y}
			}
		}
	}
	return Point{-1, -1}
}

func createScoreGrid(width, height int) Grid[FieldScore] {
	grid := make(Grid[FieldScore], height)
	for i := 0; i < height; i++ {
		grid[i] = make([]FieldScore, width)
		for j := 0; j < width; j++ {
			grid[i][j] = map[Point]int{
				NORTH: math.MaxInt,
				EAST:  math.MaxInt,
				SOUTH: math.MaxInt,
				WEST:  math.MaxInt,
			}
		}
	}
	return grid
}

func minInt(nums ...int) int {
	minValue := nums[0]
	for _, num := range nums[1:] {
		if num < minValue {
			minValue = num
		}
	}
	return minValue
}

func printGrid[T any](grid Grid[T]) {
	sb := strings.Builder{}
	for _, row := range grid {
		for _, cell := range row {
			s := fmt.Sprintf("%v ", cell)
			if len(s) > 5 {
				s = s[:5]
			}
			sb.WriteString(fmt.Sprintf("%6s", s))
		}
		sb.WriteString("\n")
	}
	fmt.Printf(sb.String())
}

func calculateScores(grid Grid[rune]) Grid[FieldScore] {
	start := findSymbol(grid, 'S')
	scores := createScoreGrid(grid.width(), grid.height())
	scores.setValue(&start, FieldScore{
		NORTH: 1000,
		EAST:  0,
		SOUTH: 1000,
		WEST:  2000,
	})
	queue := []*Point{&start}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		neighbours := getNeighbours(current, grid)
		for _, neighbour := range neighbours {
			currentScore := scores.getValue(current)
			neighbourScore := scores.getValue(&neighbour)
			dir := subPoint(&neighbour, current)
			newScore := currentScore[dir] + 1
			isBetter := false
			if neighbourScore[dir] > newScore {
				neighbourScore[dir] = newScore
				isBetter = true
			}
			newScore += 1000
			left := turnLeft(&dir)
			right := turnRight(&dir)
			opposite := getNegativePoint(&dir)
			if neighbourScore[left] > newScore {
				neighbourScore[left] = newScore
				isBetter = true
			}
			if neighbourScore[right] > newScore {
				neighbourScore[right] = newScore
				isBetter = true
			}
			newScore += 1000
			if neighbourScore[opposite] > newScore {
				neighbourScore[opposite] = newScore
			}
			if isBetter {
				scores.setValue(&neighbour, neighbourScore)
				queue = append(queue, &neighbour)
			}
		}
	}
	return scores
}

func part1(filename string) {
	grid := readData(filename)
	end := findSymbol(grid, 'E')
	scores := calculateScores(grid)
	endScores := scores.getValue(&end)
	minScore := minInt(endScores[NORTH], endScores[EAST], endScores[SOUTH], endScores[WEST])
	fmt.Printf("Part 1: minimum score: %d\n", minScore)
}

type Status struct {
	location  Point
	direction Point
}

func part2(filename string) {
	grid := readData(filename)
	end := findSymbol(grid, 'E')
	scores := calculateScores(grid)
	endScores := scores.getValue(&end)
	minScore := minInt(endScores[NORTH], endScores[EAST], endScores[SOUTH], endScores[WEST])
	var endDir Point
	for v, k := range endScores {
		if k == minScore {
			endDir = getNegativePoint(&v)
			break
		}
	}
	queue := []*Status{
		{
			location:  end,
			direction: endDir,
		},
	}
	visited := make(map[Point]bool)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visited[current.location] = true

		neighbour := addPoints(&current.location, &current.direction)
		oppositeDir := getNegativePoint(&current.direction)
		currentScore := scores.getValue(&current.location)[oppositeDir]
		if !visited[neighbour] && currentScore-scores.getValue(&neighbour)[oppositeDir] == 1 {
			queue = append(queue, &Status{location: neighbour, direction: current.direction})
		}
		dirs := []Point{
			turnLeft(&current.direction),
			turnRight(&current.direction),
		}
		for _, dir := range dirs {
			neighbour = addPoints(&current.location, &dir)
			oppositeDir = getNegativePoint(&dir)
			if !visited[neighbour] && currentScore-scores.getValue(&neighbour)[oppositeDir] == 1001 {
				queue = append(queue, &Status{location: neighbour, direction: dir})
			}
		}
	}
	fmt.Printf("Part 2: best path tiles: %d\n", len(visited))
}

func main() {
	part1("input.txt")
	part2("input.txt")
}
