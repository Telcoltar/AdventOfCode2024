package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

func readData(filepath string) ([][2]string, error) {
	fileContent, readFileErr := os.ReadFile(filepath)
	if readFileErr != nil {
		return nil, readFileErr
	}
	connections := make([][2]string, 0)
	for _, line := range strings.Split(string(fileContent), "\n") {
		split := strings.Split(strings.TrimSpace(line), "-")
		connections = append(connections, [2]string{split[0], split[1]})
		connections = append(connections, [2]string{split[1], split[0]})
	}
	return connections, nil
}

func buildConnectionMap(connections [][2]string) map[string][]string {
	connectionMap := make(map[string][]string)
	for _, connection := range connections {
		if _, exists := connectionMap[connection[0]]; !exists {
			connectionMap[connection[0]] = make([]string, 0)
		}
		connectionMap[connection[0]] = append(connectionMap[connection[0]], connection[1])
	}
	return connectionMap
}

func getPairs(items []string) [][2]string {
	pairs := make([][2]string, 0)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			pairs = append(pairs, [2]string{items[i], items[j]})
		}
	}
	return pairs
}

func solutionPart1(connections [][2]string) int {
	connectionsMap := buildConnectionMap(connections)
	connectionSet := make(map[string]bool)
	for _, connection := range connections {
		connectionSet[strings.Join(connection[:], "-")] = true
	}

	triplesSet := make(map[string]bool)

	for node, connectedNodes := range connectionsMap {
		// for pairs in connectedNodes
		for _, pair := range getPairs(connectedNodes) {
			if _, exists := connectionSet[strings.Join(pair[:], "-")]; exists {
				nodes := []string{node, pair[0], pair[1]}
				slices.Sort(nodes)
				triplesSet[strings.Join(nodes, "-")] = true
			}
		}

	}
	countOfTriplesWithT := 0
	for triple := range triplesSet {
		if triple[0] == 't' || triple[3] == 't' || triple[6] == 't' {
			countOfTriplesWithT++
		}
	}
	return countOfTriplesWithT
}

func countConnectedPairs(connectedNodes []string, connectionSet map[string]bool) int {
	count := 0
	for _, pair := range getPairs(connectedNodes) {
		if _, exists := connectionSet[strings.Join(pair[:], "-")]; exists {
			count++
		}
	}
	return count
}

func filterList(list []string, filter []string) []string {
	filterSet := make(map[string]bool)
	for _, item := range filter {
		filterSet[item] = true
	}
	filteredList := make([]string, 0)
	for _, item := range list {
		if _, exists := filterSet[item]; exists {
			filteredList = append(filteredList, item)
		}
	}
	return filteredList
}

func solutionPart2(connections [][2]string) int {
	connectionsMap := buildConnectionMap(connections)
	connectionSet := make(map[string]bool)
	for _, connection := range connections {
		connectionSet[strings.Join(connection[:], "-")] = true
	}

	currentMaxLen := 0
	var currentMaxLenNodes string

	for node, connectedNodes := range connectionsMap {
		filteredNodes := connectedNodes
		for _, innerNode := range connectedNodes {
			filteredNodes = filterList(filteredNodes, connectionsMap[innerNode])
			filteredNodes = append(filteredNodes, innerNode)
		}
		filteredNodes = append(filteredNodes, node)
		if len(filteredNodes) > currentMaxLen {
			currentMaxLen = len(filteredNodes)
			slices.Sort(filteredNodes)
			currentMaxLenNodes = strings.Join(filteredNodes, ",")
		}
	}
	println(currentMaxLen)
	println(currentMaxLenNodes)
	return 0
}

func main() {
	connections, readDataErr := readData("input.txt")
	if readDataErr != nil {
		panic(readDataErr)
	}
	startTime := time.Now()
	println(solutionPart1(connections))
	fmt.Printf("Time: %s\n", time.Since(startTime))
	startTime = time.Now()
	println(solutionPart2(connections))
	fmt.Printf("Time: %s\n", time.Since(startTime))
}
