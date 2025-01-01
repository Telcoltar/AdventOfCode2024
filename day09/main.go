package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func readData(filename string) []int {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	numbers := make([]int, 0)
	for _, r := range []rune(string(fileContent)) {
		parsed, err := strconv.Atoi(string(r))
		if err != nil {
			panic(err)
		}
		numbers = append(numbers, parsed)
	}
	return numbers
}

func part1(filename string) {
	data := readData(filename)
	// len
	diskSize := 0
	for _, number := range data {
		diskSize += number
	}
	disk := make([]int, diskSize)
	diskIndex := 0
	emptyCount := 0
	fullCount := 0

	for i, number := range data {
		// file
		if i%2 == 0 {
			fullCount += number
			index := i / 2
			for j := diskIndex; j < number+diskIndex; j++ {
				disk[j] = index
			}
		} else {
			emptyCount += number
			for j := diskIndex; j < number+diskIndex; j++ {
				disk[j] = -1
			}
		}
		diskIndex += number
	}
	emptyPoint := data[0]
	for i := len(disk) - 1; i >= 0; i-- {
		number := disk[i]
		if number == -1 {
			continue
		}
		for disk[emptyPoint] != -1 {
			emptyPoint++
		}
		if emptyPoint >= i {
			break
		}
		disk[emptyPoint] = number
		disk[i] = -1
		emptyPoint++
	}
	checkSum := 0
	for i, number := range disk {
		if number == -1 {
			break
		}
		checkSum += number * i
	}
	fmt.Printf("Part 1: %d\n", checkSum)
}

type Block struct {
	start, size, number int
}

func part2(filename string) {
	data := readData(filename)
	emptyBlocks := make([]*Block, 0)
	fullBlocks := make([]*Block, 0)
	index := 0
	for i, number := range data {
		if i%2 == 0 {
			fullBlocks = append(fullBlocks, &Block{index, number, i / 2})
		} else {
			emptyBlocks = append(emptyBlocks, &Block{index, number, -1})
		}
		index += number
	}
	newBlocks := make([]Block, 0)
	slices.Reverse(fullBlocks)
	for _, block := range fullBlocks {
		for _, emptyBlock := range emptyBlocks {
			if emptyBlock.start > block.start {
				break
			}
			if emptyBlock.size >= block.size {
				newBlocks = append(newBlocks,
					Block{emptyBlock.start, block.size, block.number},
				)
				emptyBlock.size = emptyBlock.size - block.size
				emptyBlock.start = emptyBlock.start + block.size
				block.number = -1
				break
			}
		}
	}
	for _, block := range fullBlocks {
		if block.number != -1 {
			newBlocks = append(newBlocks, *block)
		}
	}
	slices.SortFunc(newBlocks, func(a, b Block) int {
		return cmp.Compare(a.start, b.start)
	})
	checksum := 0
	for _, block := range newBlocks {
		for i := block.start; i < block.size+block.start; i++ {
			checksum += block.number * i
		}
	}
	fmt.Printf("Part 2: %d\n", checksum)
}

func main() {
	//part1("input.txt")
	part2("input.txt")
}
