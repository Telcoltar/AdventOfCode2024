package main

import (
	"fmt"
	"os"
	"strings"
)

func readInput(fname string) [][]rune {
	content, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	board := make([][]rune, 0)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		board = append(board, []rune(line))
	}
	return board
}

type Direction struct {
	x, y int
}

var xmas []rune = []rune("XMAS")

func checkDirection(board [][]rune, startX, startY int, direction Direction, word []rune) bool {
	currentX := startX
	currentY := startY
	currentChar := 0
	for currentX >= 0 && currentX < len(board[0]) && currentY >= 0 && currentY < len(board) && currentChar < len(word) {
		if board[currentY][currentX] == word[currentChar] {
			currentChar++
			currentX += direction.x
			currentY += direction.y
		} else {
			return false
		}
	}
	return currentChar == len(word)
}

func checkDirections(board [][]rune, startX int, startY int) int {
	directions := []Direction{
		{0, 1},
		{1, 0},
		{-1, 0},
		{0, -1},
		{1, 1},
		{-1, 1},
		{1, -1},
		{-1, -1},
	}
	sum := 0
	for _, direction := range directions {
		if checkDirection(board, startX, startY, direction, xmas) {
			sum += 1
		}
	}
	return sum
}

var mas = []rune("MAS")

type Point struct {
	x, y int
}

func checkMasDirections(board [][]rune, startX int, startY int, aLoc map[Point]int) {
	directions := []Direction{
		{1, 1},
		{-1, 1},
		{1, -1},
		{-1, -1},
	}
	for _, direction := range directions {
		if checkDirection(board, startX, startY, direction, mas) {
			aLoc[Point{startX + direction.x, startY + direction.y}] += 1
		}
	}
}

func part1(board [][]rune) {
	sum := 0
	for indY, row := range board {
		for indX, cell := range row {
			if cell == 'X' {
				sum += checkDirections(board, indX, indY)
			}
		}
	}
	fmt.Println(sum)
}

func part2(board [][]rune) {
	sum := 0
	aLoc := make(map[Point]int)
	for indY, row := range board {
		for indX, cell := range row {
			if cell == 'M' {
				checkMasDirections(board, indX, indY, aLoc)
			}
		}
	}
	for _, count := range aLoc {
		if count == 2 {
			sum += 1
		}
	}
	fmt.Println(sum)
}

func main() {
	part1(readInput("input.txt"))
	part2(readInput("input.txt"))
}
