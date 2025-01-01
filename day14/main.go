package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Robot struct {
	x, y   int
	vx, vy int
}

func (r *Robot) print() {
	fmt.Printf("(%d, %d) (%d, %d)\n", r.x, r.y, r.vx, r.vy)
}

func (r *Robot) move() {
	r.x += r.vx
	r.y += r.vy
}

func (r *Robot) moveInGrid(width, height int) {
	r.x = positiveModulo(r.x+r.vx, width)
	r.y = positiveModulo(r.y+r.vy, height)
}

func (r *Robot) moveInGridBack(width, height int) {
	r.x = positiveModulo(r.x-r.vx, width)
	r.y = positiveModulo(r.y-r.vy, height)
}

func parseCommaSeparatedPair(pair string) (int, int) {
	split := strings.Split(strings.TrimSpace(pair), ",")
	first, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}
	second, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	return first, second
}

func readData(filename string) []*Robot {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	robots := make([]*Robot, 0)
	for _, line := range strings.Split(strings.TrimSpace(string(fileContent)), "\n") {
		split := strings.Split(strings.TrimSpace(line), " ")
		r := Robot{}
		r.x, r.y = parseCommaSeparatedPair(strings.Split(split[0], "=")[1])
		r.vx, r.vy = parseCommaSeparatedPair(strings.Split(split[1], "=")[1])
		robots = append(robots, &r)
	}
	return robots
}

func positiveModulo(num, modulo int) int {
	res := num % modulo
	if res < 0 {
		res += modulo
	}
	return res
}

func part1(filename string, width, height int) {
	robots := readData(filename)
	for _, robot := range robots {
		robot.x = positiveModulo(robot.x+100*robot.vx, width)
		robot.y = positiveModulo(robot.y+100*robot.vy, height)
	}
	q1 := 0
	q2 := 0
	q3 := 0
	q4 := 0
	mX := (width - 1) / 2
	mY := (height - 1) / 2
	for _, robot := range robots {
		if robot.x < mX {
			if robot.y < mY {
				q1 += 1
			} else if robot.y > mY {
				q2 += 1
			}
		} else if robot.x > mX {
			if robot.y < mY {
				q3 += 1
			} else if robot.y > mY {
				q4 += 1
			}
		}
	}
	fmt.Printf("Q1: %d, Q2: %d, Q3: %d, Q4: %d\n", q1, q2, q3, q4)
	fmt.Printf("Q1*Q2*Q3*Q4: %d\n", q1*q2*q3*q4)
}

type model struct {
	iteration     int
	width, height int
	robots        []*Robot
}

func (m *model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) moveRobots() {
	m.iteration++
	for _, robot := range m.robots {
		robot.moveInGrid(m.width, m.height)
	}
}

func (m *model) moveRobotsBack() {
	m.iteration--
	for _, robot := range m.robots {
		robot.moveInGridBack(m.width, m.height)
	}
}

type Point struct {
	x, y int
}

func (p *Point) getNeighboursInGrid(width, height int) []Point {
	neighbours := make([]Point, 0)
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *model) searchTree(grid [][]bool) int {
	maxSize := 0
	for x := 1; x < 20; x++ {
		for y := 1; y < 20; y++ {
			cluster := 0
			for i := -2; i < 3; i++ {
				for j := -2; j < 3; j++ {
					if grid[x*5+i][y*5+j] {
						cluster++
					}
				}
			}
			maxSize = maxInt(maxSize, cluster)
		}
	}
	return maxSize
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n":
			m.moveRobots()
			return m, nil
		case "p":
			m.moveRobotsBack()
			return m, nil
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func createGrid(width, height int, robots []*Robot) [][]bool {
	g := make([][]bool, height)
	for i := range g {
		g[i] = make([]bool, width)
	}
	for _, robot := range robots {
		g[robot.y][robot.x] = true
	}
	return g
}

func (m *model) View() string {
	var sb strings.Builder
	sb.WriteString("Iterations: ")
	sb.WriteString(strconv.Itoa(m.iteration))
	sb.WriteString("\n")
	sb.WriteString("\n")
	grid := createGrid(m.width, m.height, m.robots)
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				sb.WriteRune('*')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func part2(filename string) {
	robots := readData(filename)
	m := model{
		iteration: 0,
		robots:    robots,
		width:     101,
		height:    103,
	}
	for m.iteration < 20000 {
		grid := createGrid(m.width, m.height, m.robots)
		if m.searchTree(grid) > 15 {
			break
		}
		m.moveRobots()
	}
	// fmt.Printf("Iterations: %d\n", m.iteration)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	// part1("input.txt", 101, 103)
	part2("input.txt")
}
