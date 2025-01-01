package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type FloatPoint struct {
	x, y float64
}

type Puzzle struct {
	a, b, prize Point
}

type Combination struct {
	a, b, token int
}

func readData(filename string, part2 bool) []Puzzle {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	puzzles := make([]Puzzle, 0)
	for _, puzzle := range strings.Split(string(fileContent), "\n\n") {
		parts := make([]Point, 3)
		for i, line := range strings.Split(puzzle, "\n") {
			split := strings.Split(line, ": ")
			secondSplit := strings.Split(split[1], ", ")
			parts[i].x, err = strconv.Atoi(secondSplit[0][2:])
			if err != nil {
				panic(err)
			}
			parts[i].y, err = strconv.Atoi(secondSplit[1][2:])
			if err != nil {
				panic(err)
			}
		}
		if part2 {
			parts[2].x += 10000000000000
			parts[2].y += 10000000000000
		}
		puzzles = append(puzzles, Puzzle{parts[0], parts[1], parts[2]})
	}
	return puzzles
}

func isDivisible(a, b int) (int, bool) {
	if a%b == 0 {
		return a / b, true
	}
	return 0, false
}

func solution(puzzles []Puzzle) {
	sum := 0
	for _, puzzle := range puzzles {
		// check discriminant
		// fmt.Printf("%#v\n", puzzle)
		d := puzzle.a.x*puzzle.b.y - puzzle.b.x*puzzle.a.y
		if d == 0 {
			if puzzle.a.x*puzzle.prize.y-puzzle.a.y*puzzle.prize.x != 0 {
				continue
			} else {
				a := 0
				b, bDiv := isDivisible(puzzle.prize.x, puzzle.b.x)
				for b >= 0 {
					if bDiv {
						// fmt.Printf("a: %d, b: %d\n", a, b)
						sum += a*3 + b
						break
					}
					a++
					b, bDiv = isDivisible(puzzle.prize.x-a*puzzle.a.x, puzzle.b.x)
				}
			}
		}
		a := puzzle.b.y*puzzle.prize.x - puzzle.b.x*puzzle.prize.y
		b := puzzle.a.x*puzzle.prize.y - puzzle.a.y*puzzle.prize.x
		// fmt.Printf("d: %d, a: %d, b: %d\n", d, a, b)
		aRed, aDiv := isDivisible(a, d)
		if !aDiv {
			// fmt.Printf("a is not divisible by d\n")
			continue
		}
		bRed, bDiv := isDivisible(b, d)
		if !bDiv {
			// fmt.Printf("b is not divisible by d\n")
			continue
		}
		// fmt.Printf("aRed: %d, bRed: %d\n", aRed, bRed)
		sum += 3*aRed + bRed
	}
	fmt.Printf("Sum: %d\n", sum)
}

func part1(filename string) {
	puzzles := readData(filename, false)
	solution(puzzles)
}

func part2(filename string) {
	puzzles := readData(filename, true)
	solution(puzzles)
}

func main() {
	part1("input.txt")
	part2("input.txt")
}
