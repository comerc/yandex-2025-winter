package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<16)
	writer := bufio.NewWriterSize(os.Stdout, 1<<16)
	defer writer.Flush()

	// Читаем 10 чисел
	nums := make([]int, 10)
	for i := 0; i < 10; i++ {
		line, _ := reader.ReadString('\n')
		nums[i], _ = strconv.Atoi(strings.TrimSpace(line))
	}

	result := solve(nums)
	writer.WriteString(fmt.Sprintf("%d\n", result))
}

// solve находит сумму подмножества, ближайшую к 100
// При равном расстоянии выбирает большую сумму
func solve(nums []int) int {
	const target = 100
	bestSum := 0
	bestDist := target // расстояние от 0 до 100

	// Перебираем все 2^10 = 1024 подмножества
	for mask := 0; mask < (1 << 10); mask++ {
		sum := 0
		for i := 0; i < 10; i++ {
			if mask&(1<<i) != 0 {
				sum += nums[i]
			}
		}

		dist := abs(sum - target)

		// Выбираем лучшую сумму:
		// - меньшее расстояние до 100
		// - при равном расстоянии — большую сумму
		if dist < bestDist || (dist == bestDist && sum > bestSum) {
			bestSum = sum
			bestDist = dist
		}
	}

	return bestSum
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
