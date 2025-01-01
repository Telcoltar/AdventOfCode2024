package main

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func part1(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	line := string(content)
	sum := processSection(line)
	fmt.Println(sum)
}

func processSection(section string) int {
	mulMatcher := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	sum := 0
	for _, match := range mulMatcher.FindAllStringSubmatchIndex(section, -1) {
		firstSubMatch, err := strconv.Atoi(section[match[2]:match[3]])
		if err != nil {
			panic(err)
		}
		secondSubMatch, err := strconv.Atoi(section[match[4]:match[5]])
		if err != nil {
			panic(err)
		}
		multResult := firstSubMatch * secondSubMatch
		sum += multResult
	}
	return sum
}

type maskPoint struct {
	index int
	value int
}

func part2(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	line := string(content)
	doMatcher := regexp.MustCompile(`do\(\)`)
	masked := make([]maskPoint, 0)
	doMatches := doMatcher.FindAllStringIndex(line, -1)
	for _, match := range doMatches {
		masked = append(masked, maskPoint{match[0], 1})
	}
	dontMatcher := regexp.MustCompile(`don't\(\)`)
	for _, match := range dontMatcher.FindAllStringIndex(line, -1) {
		masked = append(masked, maskPoint{match[0], 0})
	}
	slices.SortFunc(masked, func(a, b maskPoint) int {
		return cmp.Compare(a.index, b.index)
	})
	cleanedMask := make([]maskPoint, 0)
	cleanedMask = append(cleanedMask, maskPoint{index: 0, value: 1})
	current := 1
	for _, mask := range masked {
		if mask.value != current {
			cleanedMask = append(cleanedMask, mask)
			current = mask.value
		}
	}
	if cleanedMask[len(cleanedMask)-1].value == 1 {
		cleanedMask = append(cleanedMask, maskPoint{index: len(line), value: 0})
	}
	sum := 0
	for i := 0; i < len(cleanedMask); i += 2 {
		start := cleanedMask[i].index
		stop := cleanedMask[i+1].index
		sum += processSection(line[start:stop])
	}
	fmt.Println(sum)
}

func main() {
	part1("input.txt")
	part2("input.txt")
}
