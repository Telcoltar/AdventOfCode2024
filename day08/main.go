package main

import (
	"fmt"
	"os"
	"strings"
)

type Grid = [][]rune

func readData(filename string) Grid {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var grid Grid
	for _, line := range strings.Split(string(fileContents), "\n") {
		grid = append(grid, []rune(line))
	}
	return grid
}

type Point struct {
	x, y int
}

func (p *Point) Sub(other Point) Point {
	return Point{p.x - other.x, p.y - other.y}
}

func (p *Point) Add(other Point) Point {
	return Point{p.x + other.x, p.y + other.y}
}

func (p *Point) checkBound(width, height int) bool {
	if p.x < 0 || p.x >= width || p.y < 0 || p.y >= height {
		return false
	}
	return true
}

func getAntennas(grid Grid) map[rune][]Point {
	var antennas map[rune][]Point = make(map[rune][]Point, 0)
	for idy, row := range grid {
		for idx, ch := range row {
			if ch != '.' {
				antennas[ch] = append(antennas[ch], Point{idx, idy})
			}
		}
	}
	return antennas
}

func calcAntinodes(antennas []Point, width, height int) []Point {
	antinodes := make([]Point, 0)
	for i := 0; i < len(antennas); i++ {
		for j := i + 1; j < len(antennas); j++ {
			first := antennas[i]
			second := antennas[j]
			diff := first.Sub(second)
			firstAntinode := first.Add(diff)
			if firstAntinode.checkBound(width, height) {
				antinodes = append(antinodes, firstAntinode)
			}
			secondAntinode := second.Sub(diff)
			if secondAntinode.checkBound(width, height) {
				antinodes = append(antinodes, secondAntinode)
			}
		}
	}
	return antinodes
}

func calcAntinodesHarmnics(antennas []Point, width, height int) []Point {
	antinodes := make([]Point, 0)
	for i := 0; i < len(antennas); i++ {
		for j := i + 1; j < len(antennas); j++ {
			first := antennas[i]
			second := antennas[j]
			diff := first.Sub(second)

			firstAntinode := first
			for firstAntinode.checkBound(width, height) {
				antinodes = append(antinodes, firstAntinode)
				firstAntinode = firstAntinode.Add(diff)
			}

			secondAntinode := second
			for secondAntinode.checkBound(width, height) {
				antinodes = append(antinodes, secondAntinode)
				secondAntinode = secondAntinode.Sub(diff)
			}
		}
	}
	return antinodes
}

func part1(filename string) {
	grid := readData(filename)
	antennas := getAntennas(grid)
	antinodeLocation := make(map[Point]struct{}, 0)
	for _, locations := range antennas {
		for _, loc := range calcAntinodes(locations, len(grid[0]), len(grid)) {
			antinodeLocation[loc] = struct{}{}
		}
	}
	fmt.Println(len(antinodeLocation))
}

func part2(filename string) {
	grid := readData(filename)
	antennas := getAntennas(grid)
	antinodeLocation := make(map[Point]struct{}, 0)
	for _, locations := range antennas {
		for _, loc := range calcAntinodesHarmnics(locations, len(grid[0]), len(grid)) {
			antinodeLocation[loc] = struct{}{}
		}
	}
	fmt.Println(len(antinodeLocation))
}

func main() {
	part1("input.txt")
	part2("input.txt")
}
