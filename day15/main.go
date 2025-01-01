package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strconv"
	"strings"
)

func doubleLine(line []rune) []rune {
	doubled := make([]rune, 0, len(line)*2)
	for _, r := range line {
		switch r {
		case '#':
			doubled = append(doubled, '#', '#')
		case 'O':
			doubled = append(doubled, '[', ']')
		case '.':
			doubled = append(doubled, '.', '.')
		case '@':
			doubled = append(doubled, '@', '.')
		}
	}
	return doubled
}

func readData(filename string, double bool) ([][]rune, []rune) {
	grid := make([][]rune, 0)
	moves := make([]rune, 0)
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(fileContent), "\n")
	i := 0
	for {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			break
		}
		if double {
			grid = append(grid, doubleLine([]rune(line)))
		} else {
			grid = append(grid, []rune(line))
		}
		i++
	}
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		moves = append(moves, []rune(line)...)
		i++
	}
	return grid, moves
}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}

func printMoves(moves []rune) {
	for _, move := range moves {
		fmt.Printf("%c", move)
	}
	fmt.Println()
}

type Point struct {
	x, y int
}

func (p *Point) opposite() Point {
	return Point{
		-p.x,
		-p.y,
	}
}

func addPoints(p1, p2 Point) Point {
	return Point{p1.x + p2.x, p1.y + p2.y}
}

func getRobotPosition(grid [][]rune) Point {
	for y, row := range grid {
		for x, cell := range row {
			if cell == '@' {
				return Point{x: x, y: y}
			}
		}
	}
	return Point{}
}

func getDirForMove(move rune) Point {
	switch move {
	case '^':
		return Point{0, -1}
	case '>':
		return Point{1, 0}
	case 'v':
		return Point{0, 1}
	case '<':
		return Point{-1, 0}
	}
	return Point{}
}

func moveRobot(start Point, dir Point, grid [][]rune) bool {
	currentPos := addPoints(start, dir)
	currentGridValue := grid[currentPos.y][currentPos.x]
	for currentGridValue == 'O' {
		currentPos = addPoints(currentPos, dir)
		currentGridValue = grid[currentPos.y][currentPos.x]
	}
	if currentGridValue == '#' {
		return false
	}
	grid[currentPos.y][currentPos.x] = 'O'
	nextPos := addPoints(start, dir)
	grid[nextPos.y][nextPos.x] = '.'
	return true
}

func countGPSCoordinates(grid [][]rune) int {
	sum := 0
	for y, row := range grid {
		for x, cell := range row {
			if cell == 'O' || cell == '[' {
				sum += y*100 + x
			}
		}
	}
	return sum
}

func part1(filename string) {
	grid, moves := readData(filename, false)
	robot := getRobotPosition(grid)
	grid[robot.y][robot.x] = '.'
	for _, move := range moves {
		dir := getDirForMove(move)
		if moveRobot(robot, dir, grid) {
			grid[robot.y][robot.x] = '.'
			robot = addPoints(robot, dir)
			grid[robot.y][robot.x] = '@'
		}
	}
	printGrid(grid)
	fmt.Printf("%d\n", countGPSCoordinates(grid))
}

type model struct {
	moveIndex int
	grid      [][]rune
	moves     []rune
	field     map[Point]struct{}
	endPoints map[Point]struct{}
	robot     Point
	gps       int
}

func (m *model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
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
			if m.moveIndex < len(m.moves) {
				dir := getDirForMove(m.moves[m.moveIndex])
				if dir.y == 0 {
					m.moveRobot2Y(dir)
					m.moveIndex++
				} else {
					m.moveRobot2X(dir)
					m.moveIndex++
				}
			}
			return m, nil
		case "s":
			for _, move := range m.moves[m.moveIndex:] {
				dir := getDirForMove(move)
				if dir.y == 0 {
					m.moveRobot2Y(dir)
					m.moveIndex++
				} else {
					m.moveRobot2X(dir)
					m.moveIndex++
				}
			}
			m.gps = countGPSCoordinates(m.grid)
			return m, nil
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *model) View() string {
	styleField := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	styleEnd := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	var sb strings.Builder
	sb.WriteString("Index: ")
	sb.WriteString(strconv.Itoa(m.moveIndex))
	var move rune
	if m.moveIndex < len(m.moves) {
		move = m.moves[m.moveIndex]
	} else {
		move = 'x'
	}
	sb.WriteString(fmt.Sprintf(", next move %c", move))
	sb.WriteString("\n")
	sb.WriteString("Field size: ")
	sb.WriteString(strconv.Itoa(len(m.field)))
	sb.WriteString("\n")
	sb.WriteString("EndPoints size: ")
	sb.WriteString(strconv.Itoa(len(m.endPoints)))
	sb.WriteString("\n")
	sb.WriteString("GPS: ")
	sb.WriteString(strconv.Itoa(m.gps))
	sb.WriteString("\n")
	sb.WriteString("\n")
	for y, row := range m.grid {
		for x, cell := range row {
			p := Point{x, y}
			if p == m.robot {
				sb.WriteString("@")
			} else {
				if _, ok := m.field[p]; ok {
					sb.WriteString(styleField.Render(fmt.Sprintf("%c", cell)))
				} else if _, endOk := m.endPoints[p]; endOk {
					sb.WriteString(styleEnd.Render(fmt.Sprintf("%c", cell)))
				} else {
					sb.WriteString(fmt.Sprintf("%c", cell))
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m *model) moveRobot2Y(dir Point) bool {
	currentPos := addPoints(m.robot, dir)
	currentGridValue := m.grid[currentPos.y][currentPos.x]
	for currentGridValue == '[' || currentGridValue == ']' {
		currentPos = addPoints(currentPos, dir)
		currentGridValue = m.grid[currentPos.y][currentPos.x]
	}
	if currentGridValue == '#' {
		return false
	}
	oppositeDir := dir.opposite()
	nextPos := addPoints(currentPos, oppositeDir)
	for nextPos.x != m.robot.x {
		m.grid[currentPos.y][currentPos.x] = m.grid[nextPos.y][nextPos.x]
		currentPos = addPoints(currentPos, oppositeDir)
		nextPos = addPoints(nextPos, oppositeDir)
	}
	m.grid[currentPos.y][currentPos.x] = '.'
	m.robot = currentPos
	return true
}

var RIGHT = Point{1, 0}
var LEFT = Point{-1, 0}

func (m *model) moveRobot2X(dir Point) bool {
	start := addPoints(m.robot, dir)
	if m.grid[start.y][start.x] == '.' {
		m.robot = start
		return true
	}
	if m.grid[start.y][start.x] == '#' {
		return false
	}
	queue := []Point{start}
	visited := make(map[Point]struct{})
	field := make(map[Point]struct{})
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if _, ok := visited[current]; ok {
			continue
		}
		visited[current] = struct{}{}
		if m.grid[current.y][current.x] == '[' {
			field[current] = struct{}{}
			queue = append(queue, addPoints(current, dir))
			queue = append(queue, addPoints(current, RIGHT))
		} else if m.grid[current.y][current.x] == ']' {
			field[current] = struct{}{}
			queue = append(queue, addPoints(current, dir))
			queue = append(queue, addPoints(current, LEFT))
		}
	}
	m.field = field
	endPoints := make(map[Point]struct{})
	for p := range field {
		nextPoint := addPoints(p, dir)
		if _, ok := field[nextPoint]; !ok {
			endPoints[nextPoint] = struct{}{}
		}
	}
	m.endPoints = endPoints
	for p := range endPoints {
		if m.grid[p.y][p.x] != '.' {
			return false
		}
	}
	oppositeDir := dir.opposite()
	for p := range endPoints {
		currentPos := p
		nextPost := addPoints(currentPos, oppositeDir)
		_, ok := m.field[nextPost]
		for ok {
			m.grid[currentPos.y][currentPos.x] = m.grid[nextPost.y][nextPost.x]
			currentPos = addPoints(currentPos, oppositeDir)
			nextPost = addPoints(nextPost, oppositeDir)
			_, ok = m.field[nextPost]
		}
		m.grid[currentPos.y][currentPos.x] = '.'
	}
	m.robot = start
	return true
}

func part2(filename string) {
	grid, moves := readData(filename, true)
	robot := getRobotPosition(grid)
	grid[robot.y][robot.x] = '.'
	m := &model{
		moveIndex: 0,
		grid:      grid,
		moves:     moves,
		robot:     robot,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	// part1("input.txt")
	part2("input.txt")
}
