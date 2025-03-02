package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-set/v3"
)

type Operator string

const AND Operator = "AND"
const OR Operator = "OR"
const XOR Operator = "XOR"

type Op struct {
	Left, Right, Out string
	Operator         Operator
}

func readInput(filepath string) (map[string]bool, []Op, error) {
	fileContent, readFileErr := os.ReadFile(filepath)
	if readFileErr != nil {
		return nil, nil, readFileErr
	}
	parts := strings.Split(string(fileContent), "\n\n")
	inputs := make(map[string]bool)
	for _, input := range strings.Split(parts[0], "\n") {
		inputParts := strings.Split(input, ": ")
		inputs[strings.TrimSpace(inputParts[0])] = strings.TrimSpace(inputParts[1]) == "1"
	}

	ops := make([]Op, 0)
	for _, op := range strings.Split(parts[1], "\n") {
		opParts := strings.Split(op, " ")
		ops = append(ops, Op{
			Left:     strings.TrimSpace(opParts[0]),
			Right:    strings.TrimSpace(opParts[2]),
			Out:      strings.TrimSpace(opParts[4]),
			Operator: Operator(opParts[1]),
		})
	}
	return inputs, ops, nil
}

func getValue(output string, outputMap map[string]Op, cache map[string]bool, depth int) bool {
	if depth > 1000 {
		panic("to deep")
	}
	if value, ok := cache[output]; ok {
		return value
	}
	op, ok := outputMap[output]
	if !ok {
		panic("got output value without op")
	}
	left := getValue(op.Left, outputMap, cache, depth+1)
	right := getValue(op.Right, outputMap, cache, depth+1)
	switch op.Operator {
	case AND:
		cache[output] = left && right
		return left && right
	case OR:
		cache[output] = left || right
		return left || right
	case XOR:
		cache[output] = left != right
		return left != right
	}
	panic("should not land here")
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func getNumFromBits(cache map[string]bool, indicator string) (int64, error) {
	bits := make(map[string]bool)
	for label, value := range cache {
		if !strings.Contains(label, indicator) {
			continue
		}
		bits[label] = value
	}
	keys := slices.Collect(maps.Keys(bits))
	slices.Sort(keys)
	slices.Reverse(keys)
	bitString := strings.Builder{}
	for _, key := range keys {
		bitString.WriteString(fmt.Sprintf("%d", btoi(cache[key])))
	}
	return strconv.ParseInt(bitString.String(), 2, 64)
}

func solutionPart1(filepath string) int {
	inputs, ops, err := readInput(filepath)
	if err != nil {
		panic(err)
	}
	outputMap := make(map[string]Op)
	for _, op := range ops {
		outputMap[op.Out] = op
	}
	cache := inputs

	for output := range outputMap {
		if !strings.Contains(output, "z") {
			continue
		}

		getValue(output, outputMap, cache, 0)
	}

	resultNum, boolToNumErr := getNumFromBits(cache, "z")
	if boolToNumErr != nil {
		panic(boolToNumErr)
	}
	return int(resultNum)
}

func getParticipatingWires(outputMap map[string]Op, cache map[string]*set.Set[string], startWire string) *set.Set[string] {
	if wires, ok := cache[startWire]; ok {
		return wires
	}
	if strings.Contains(startWire, "x") || strings.Contains(startWire, "y") {
		return set.From([]string{startWire})
	}

	op, ok := outputMap[startWire]
	if !ok {
		panic("got output value without op")
	}
	left := getParticipatingWires(outputMap, cache, op.Left)
	right := getParticipatingWires(outputMap, cache, op.Right)

	wires := set.New[string](1 + left.Size() + right.Size())
	wires.Insert(startWire)
	wires.InsertSet(right)
	wires.InsertSet(left)
	cache[startWire] = wires
	return wires
}

func getDefectPositions(a string, b string) []int {
	defects := make([]int, 0)
	for idx, el := range a {
		if el != []rune(b)[idx] {
			defects = append(defects, idx)
		}
	}
	return defects
}

func calculateOutput(outputMap map[string]Op, inputs map[string]bool) (string, string) {
	cache := inputs

	for output := range outputMap {
		if !strings.Contains(output, "z") {
			continue
		}

		getValue(output, outputMap, cache, 0)
	}
	output, _ := getNumFromBits(cache, "z")
	x, _ := getNumFromBits(cache, "x")
	y, _ := getNumFromBits(cache, "y")
	return strconv.FormatInt(output, 2), strconv.FormatInt(x+y, 2)
}

func getValueWithBothPerm(inputMap map[string]string, a, op, b string) string {
	pos1, ok1 := inputMap[fmt.Sprintf("%s:%s:%s", a, op, b)]
	pos2, ok2 := inputMap[fmt.Sprintf("%s:%s:%s", b, op, a)]
	if (ok1 && ok2) && (pos1 != pos2) {
		panic("permutations have different values")
	}
	if !ok1 && !ok2 {
		panic(fmt.Sprintf("no values for %s:%s:%s", a, op, b))
	}
	if ok1 {
		return pos1
	}
	return pos2
}

func setValueWithBothPerm(inputMap map[string]string, a, op, b string, setValue string) {
	if _, ok1 := inputMap[fmt.Sprintf("%s:%s:%s", a, op, b)]; ok1 {
		inputMap[fmt.Sprintf("%s:%s:%s", a, op, b)] = setValue
	}
	inputMap[fmt.Sprintf("%s:%s:%s", b, op, a)] = setValue
}

func solutionPart2(filepath string) string {
	inputs, ops, err := readInput(filepath)
	if err != nil {
		panic(err)
	}
	outputMap := make(map[string]Op)
	for _, op := range ops {
		outputMap[op.Out] = op
	}
	inputMap := make(map[string]string)
	for _, op := range ops {
		inputMap[fmt.Sprintf("%s:%s:%s", op.Left, op.Operator, op.Right)] = op.Out
	}

	calculateOutput(outputMap, inputs)

	keys := slices.Collect(maps.Keys(outputMap))
	zs := make([]string, 0)
	for _, key := range keys {
		if strings.Contains(key, "z") {
			zs = append(zs, key)
		}
	}
	bits := len(zs)

	outputAdder := make([]string, bits)
	carryOut := make([]string, bits)
	carryOut[0] = "hjp"

	for i := 1; i < bits-1; i++ {
		fmt.Printf("Start: %d\n", i)
		currentZ := outputMap[fmt.Sprintf("z%02d", i)]

		xor := getValueWithBothPerm(inputMap, fmt.Sprintf("x%02d", i), "XOR", fmt.Sprintf("y%02d", i))
		if currentZ.Operator != XOR {
			fmt.Println("Z operator not XOR")
			z := getValueWithBothPerm(inputMap, xor, "XOR", carryOut[i-1])
			realZ := outputMap[z]

			outputMap[fmt.Sprintf("z%02d", i)] = realZ
			outputMap[z] = currentZ

			setValueWithBothPerm(inputMap, xor, "XOR", carryOut[i-1], fmt.Sprintf("z%02d", i))
			setValueWithBothPerm(inputMap, currentZ.Left, string(currentZ.Operator), currentZ.Right, realZ.Out)

			currentZ = outputMap[fmt.Sprintf("z%02d", i)]
		}

		if xor != currentZ.Left && xor != currentZ.Right {
			var targetXor string
			if currentZ.Left == carryOut[i-1] {
				targetXor = currentZ.Right
			} else if currentZ.Right == carryOut[i-1] {
				targetXor = currentZ.Left
			} else {
				panic("currentZ should have one part from carryOut")
			}

			xorOp, targetXorOp := outputMap[xor], outputMap[targetXor]
			outputMap[xor], outputMap[targetXor] = targetXorOp, xorOp

			setValueWithBothPerm(inputMap, xorOp.Left, "XOR", xorOp.Right, targetXor)
			setValueWithBothPerm(inputMap, targetXorOp.Left, string(targetXorOp.Operator), targetXorOp.Right, xor)

			xor = getValueWithBothPerm(inputMap, fmt.Sprintf("x%02d", i), "XOR", fmt.Sprintf("y%02d", i))
		}
		z := getValueWithBothPerm(inputMap, xor, "XOR", carryOut[i-1])
		outputAdder[i] = z

		fmt.Printf("XOR: %s\n", xor)
		fmt.Printf("Z: %s\n", z)

		and := getValueWithBothPerm(inputMap, fmt.Sprintf("x%02d", i), "AND", fmt.Sprintf("y%02d", i))
		c2 := getValueWithBothPerm(inputMap, carryOut[i-1], "AND", xor)
		c := getValueWithBothPerm(inputMap, and, "OR", c2)

		carryOut[i] = c

		fmt.Printf("AND: %s\n", and)
		fmt.Printf("C2: %s\n", c2)
		fmt.Printf("C: %s\n", c)
	}

	originalInputMap := make(map[string]string)
	for _, op := range ops {
		originalInputMap[fmt.Sprintf("%s:%s:%s", op.Left, op.Operator, op.Right)] = op.Out
	}

	swappedWires := set.New[string](0)
	for key := range originalInputMap {
		if originalInputMap[key] != inputMap[key] {
			swappedWires.Insert(originalInputMap[key])
			swappedWires.Insert(inputMap[key])
		}
	}

	swappedWiresSlice := swappedWires.Slice()
	slices.Sort(swappedWiresSlice)
	return strings.Join(swappedWiresSlice, ",")
}

func main() {
	fmt.Println(solutionPart1("input.txt"))
	startTime := time.Now()
	fmt.Println(solutionPart2("input.txt"))
	fmt.Println(time.Since(startTime))
}
