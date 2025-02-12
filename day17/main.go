package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Computer struct {
	A int
	B int
	C int

	InstructionPointer int
	Instructions       []int
	Output             []int
}

func getComboOpValue(computer *Computer, operand int) int {
	switch {
	case operand < 4:
		return operand
	case operand == 4:
		return computer.A
	case operand == 5:
		return computer.B
	case operand == 6:
		return computer.C
	case operand == 7:
		panic("invalid operand 7")
	}
	return -1
}

func formatIntegerSlice(in []int) string {
	stringSlice := make([]string, len(in))
	for i, v := range in {
		stringSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(stringSlice, ",")
}

func (c *Computer) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("A: %d, B: %d, C: %d\n", c.A, c.B, c.C))
	sb.WriteString(fmt.Sprintf("InstructionsPointer: %d, \n", c.InstructionPointer))
	sb.WriteString(fmt.Sprintf("Instructions: %v\n", c.Instructions))
	sb.WriteString(fmt.Sprintf("Output: %s\n", formatIntegerSlice(c.Output)))
	return sb.String()
}

func (c *Computer) Step() bool {
	if c.InstructionPointer >= len(c.Instructions) {
		return false
	}
	instruction := c.Instructions[c.InstructionPointer]
	operand := c.Instructions[c.InstructionPointer+1]
	c.InstructionPointer += 2
	switch instruction {
	case 0:
		operand = getComboOpValue(c, operand)
		c.A = c.A / IntPow(2, operand)
	case 1:
		c.B = c.B ^ operand
	case 2:
		operand = getComboOpValue(c, operand)
		c.B = operand % 8
	case 3:
		if c.A != 0 {
			c.InstructionPointer = operand
		}
	case 4:
		c.B = c.B ^ c.C
	case 5:
		c.Output = append(c.Output, getComboOpValue(c, operand)%8)
	case 6:
		operand = getComboOpValue(c, operand)
		c.B = c.A / IntPow(2, operand)
	case 7:
		operand = getComboOpValue(c, operand)
		c.C = c.A / IntPow(2, operand)
	}
	return true
}

func (c *Computer) Run() {
	for c.Step() {
	}
}

func (c *Computer) Reset() {
	c.InstructionPointer = 0
	c.Output = []int{}
}

func IntPow(a, b int) int {
	if b == 0 {
		return 1
	}
	if b == 1 {
		return a
	}

	output := a
	for i := 2; i <= b; i++ {
		output *= a
	}
	return output
}

func MustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func ParseCommaSeperatedString(s string) []int {
	s = strings.TrimSpace(s)
	split := strings.Split(s, ",")
	result := make([]int, len(split))
	for j, i := range split {
		result[j] = MustParseInt(i)
	}
	return result
}

func readData(filename string) Computer {
	com := Computer{}
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(fileContent), "\n")
	registerRegex := regexp.MustCompile("Register [ABC]: (\\d+)")
	aMatch := registerRegex.FindStringSubmatch(strings.TrimSpace(lines[0]))
	bMatch := registerRegex.FindStringSubmatch(strings.TrimSpace(lines[1]))
	cMatch := registerRegex.FindStringSubmatch(strings.TrimSpace(lines[2]))

	programMatcher := regexp.MustCompile("Program: ([0-9,]+)")
	programMatch := programMatcher.FindStringSubmatch(strings.TrimSpace(lines[4]))

	com.A = MustParseInt(aMatch[1])
	com.B = MustParseInt(bMatch[1])
	com.C = MustParseInt(cMatch[1])

	com.Instructions = ParseCommaSeperatedString(programMatch[1])

	return com
}

func part1(filename string) {
	c := readData(filename)
	c.Run()
	fmt.Println(c.String())
}

func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func part2() {
	table := make([][]int, 8)
	for i := 0; i < 8; i++ {
		table[i] = make([]int, 0)
	}
	for i := 0; i < 1024; i++ {
		table[combinedFunction(i)] = append(table[combinedFunction(i)], i)
	}
	program := []int{2, 4, 1, 6, 7, 5, 4, 6, 1, 4, 5, 5, 0, 3, 3, 0}
	d := []int{table[0][0]}
	for i := len(program) - 2; i >= 0; i-- {
		d = findElementsWithRemainders(table[program[i]], d)
	}
	fmt.Println(d[0])
	com := Computer{
		A:            d[0],
		Instructions: program,
	}
	com.Run()
	fmt.Println(com.String())
}

func divEightSlice(a []int) []int {
	result := make([]int, len(a))
	for i := range a {
		result[i] = a[i] / 8
	}
	return result
}

func addSlice(a []int, add int) []int {
	result := make([]int, len(a))
	for i := range a {
		result[i] = a[i] + add
	}
	return result
}

func findElementsWithRemainder(a []int, b []int, x int) []int {
	index := -1
	for i, el := range b {
		if el == x {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	start := index
	for b[index] == x {
		index += 1
	}
	return a[start:index]
}

func findElementsWithRemainders(a []int, rem []int) []int {
	b := divEightSlice(a)
	var elements []int
	for _, el := range rem {
		elements = append(elements, addSlice(findElementsWithRemainder(a, b, el%128), (el/128)*1024)...)
	}
	return elements
}

func combinedFunction(a int) int {
	return ((a % 8) ^ 2 ^ (a >> ((a % 8) ^ 6))) % 8
}

func main() {
	part1("input.txt")
	part2()
}
