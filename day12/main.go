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
	grid := make([][]rune, 0)
	for _, line := range strings.Split(strings.TrimSpace(string(fileContent)), "\n") {
		grid = append(grid, []rune(strings.TrimSpace(line)))
	}
	return grid
}

type Point struct {
	x, y int
}

func (p *Point) Add(other Point) {
	p.x = p.x + other.x
	p.y = p.y + other.y
}

func (p *Point) Neighbour(dir Point) Point {
	return Point{p.x + dir.x, p.y + dir.y}
}

func (p *Point) Sub(other Point) {
	p.x = p.x - other.x
	p.y = p.y - other.y
}

func (p *Point) RotateRight() {
	p.x = -p.y
	p.y = p.x
}

func (p *Point) RotateLeft() {
	p.x = p.y
	p.y = -p.x
}

func (p *Point) Rotate180() {
	p.x = -p.x
	p.y = -p.y
}

func (p *Point) getNeighbours(height, width int) []Point {
	neighbours := make([]Point, 0, 4)
	if p.x-1 >= 0 {
		neighbours = append(neighbours, Point{p.x - 1, p.y})
	}
	if p.x+1 < width {
		neighbours = append(neighbours, Point{p.x + 1, p.y})
	}
	if p.y-1 >= 0 {
		neighbours = append(neighbours, Point{p.x, p.y - 1})
	}
	if p.y+1 < height {
		neighbours = append(neighbours, Point{p.x, p.y + 1})
	}
	return neighbours
}

func (p *Point) getAllNeighbours() []Point {
	return []Point{
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y - 1},
		{p.x, p.y + 1},
	}
}

func expandRegion(start Point, grid [][]rune, visitedPoints map[Point]struct{}) map[Point]struct{} {
	value := grid[start.y][start.x]
	queue := make([]Point, 0)
	queue = append(queue, start)
	region := make(map[Point]struct{})
	visitedPoints[start] = struct{}{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		region[current] = struct{}{}

		for _, neighbour := range current.getNeighbours(len(grid), len(grid[0])) {
			if _, ok := visitedPoints[neighbour]; grid[neighbour.y][neighbour.x] == value && !ok {
				visitedPoints[neighbour] = struct{}{}
				queue = append(queue, neighbour)
			}
		}

	}
	return region
}

func calculateFence(region map[Point]struct{}) int {
	fence := 0
	for key := range region {
		for _, neighbour := range key.getAllNeighbours() {
			if _, ok := region[neighbour]; !ok {
				fence++
			}
		}
	}
	return fence
}

func mapHas[M ~map[K]V, K comparable, V any](key K, searchMap M) bool {
	_, ok := searchMap[key]
	return ok
}

type Edge struct {
	p1, p2 Point
}

func calculateFence2(region map[Point]struct{}) int {
	fences := make(map[Edge]struct{})
	for key := range region {
		// up
		if !mapHas(key.Neighbour(Point{0, -1}), region) {
			fences[Edge{key, key.Neighbour(Point{1, 0})}] = struct{}{}
		}
		// down
		if !mapHas(key.Neighbour(Point{0, 1}), region) {
			fences[Edge{key.Neighbour(Point{1, 1}), key.Neighbour(Point{0, 1})}] = struct{}{}
		}
		// left
		if !mapHas(key.Neighbour(Point{-1, 0}), region) {
			fences[Edge{key.Neighbour(Point{0, 1}), key}] = struct{}{}
		}
		// right
		if !mapHas(key.Neighbour(Point{1, 0}), region) {
			fences[Edge{key.Neighbour(Point{1, 0}), key.Neighbour(Point{1, 1})}] = struct{}{}
		}
	}
	fenceSegments := make([]Edge, 0)
	visited := make(map[Edge]struct{})
	for fence := range fences {
		if mapHas(fence, visited) {
			continue
		}
		dir := Point{fence.p2.x - fence.p1.x, fence.p2.y - fence.p1.y}

		currentFence := fence
		for mapHas(currentFence, fences) {
			visited[currentFence] = struct{}{}
			currentFence.p1.Add(dir)
			currentFence.p2.Add(dir)
		}
		endPoint := currentFence.p2

		dir.Rotate180()
		currentFence = fence
		for mapHas(currentFence, fences) {
			visited[currentFence] = struct{}{}
			currentFence.p1.Add(dir)
			currentFence.p2.Add(dir)
		}
		startPoint := currentFence.p1

		fenceSegments = append(fenceSegments, Edge{startPoint, endPoint})
	}

	return len(fenceSegments)
}

func solution(filename string, fenceFunc func(map[Point]struct{}) int) {
	grid := readData(filename)
	visitedPoints := make(map[Point]struct{})
	sum := 0
	count := 0
	for y, row := range grid {
		for x := range row {
			if _, ok := visitedPoints[Point{x, y}]; !ok {
				count += 1
				region := expandRegion(Point{x, y}, grid, visitedPoints)
				fence := fenceFunc(region)
				size := len(region)
				sum += fence * size
				// fmt.Printf("Point: %d, %d, size: %d, fence: %d, prod: %d\n", y, x, size, fence, fence*size)
			}
		}
	}
	fmt.Printf("Total regions: %d, sum: %d\n", count, sum)
}

func part1(filename string) {
	solution(filename, calculateFence)
}

func part2(filename string) {
	solution(filename, calculateFence2)
}

func main() {
	part1("input.txt")
	part2("input.txt")
}
