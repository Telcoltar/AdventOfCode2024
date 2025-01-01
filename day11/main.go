package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readData(filename string) []int {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	split := strings.Split(strings.TrimSpace(string(fileContent)), " ")
	result := make([]int, len(split))
	for i, s := range split {
		parsed, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			panic(err)
		}
		result[i] = parsed
	}
	return result
}

type Node struct {
	number, level int
}

type Cache struct {
	memo map[Node]int
	hits int
}

func rec(node Node, cache *Cache) int {
	if node.level == 0 {
		return 1
	}
	if length, ok := cache.memo[node]; ok {
		cache.hits++
		return length
	}
	if node.number == 0 {
		value := rec(Node{1, node.level - 1}, cache)
		cache.memo[node] = value
		return value
	}
	strNum := strconv.Itoa(node.number)
	if len(strNum)%2 == 0 {
		middle := len(strNum) / 2
		left, err := strconv.Atoi(strNum[:middle])
		if err != nil {
			panic(err)
		}
		right, err := strconv.Atoi(strNum[middle:])
		if err != nil {
			panic(err)
		}
		value := rec(Node{left, node.level - 1}, cache) +
			rec(Node{right, node.level - 1}, cache)
		cache.memo[node] = value
		return value
	}
	value := rec(Node{node.number * 2024, node.level - 1}, cache)
	cache.memo[node] = value
	return value
}

func solution(filename string, iterations int) {
	data := readData(filename)
	sum := 0
	cache := Cache{
		memo: make(map[Node]int),
		hits: 0,
	}
	for _, d := range data {
		sum += rec(Node{d, iterations}, &cache)
	}
	fmt.Printf("Cache hits: %d\n", cache.hits)
	fmt.Printf("Length: %d\n", sum)
}

func part1(filename string) {
	solution(filename, 25)
}

func part2(filename string) {
	solution(filename, 75)
}

func main() {
	start := time.Now()
	part1("input.txt")
	part2("input.txt")
	fmt.Printf("Time: %s\n", time.Since(start))
}
