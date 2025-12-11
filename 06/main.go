package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем N и M
	line, _ := reader.ReadString('\n')
	parts := strings.Fields(strings.TrimSpace(line))
	N, _ := strconv.Atoi(parts[0])
	_, _ = strconv.Atoi(parts[1]) // M - количество команд (используется при чтении строки команд)

	// Читаем координаты меток
	markersX := make([]int, N)
	markersY := make([]int, N)
	for i := 0; i < N; i++ {
		line, _ = reader.ReadString('\n')
		parts = strings.Fields(strings.TrimSpace(line))
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		markersX[i] = x
		markersY[i] = y
	}

	// Читаем команды
	line, _ = reader.ReadString('\n')
	commands := strings.TrimSpace(line)

	results := solve(markersX, markersY, commands)
	for _, result := range results {
		writer.WriteString(fmt.Sprintf("%d\n", result))
	}
}

// solve вычисляет сумму манхэттенских расстояний после каждой команды
func solve(markersX, markersY []int, commands string) []int64 {
	N := len(markersX)

	// Сортируем метки по x и y для быстрого вычисления суммы расстояний
	sortedX := make([]int, N)
	copy(sortedX, markersX)
	sort.Ints(sortedX)

	sortedY := make([]int, N)
	copy(sortedY, markersY)
	sort.Ints(sortedY)

	// Предвычисляем префиксные суммы для x и y
	prefixSumX := make([]int64, N+1)
	prefixSumY := make([]int64, N+1)
	for i := 0; i < N; i++ {
		prefixSumX[i+1] = prefixSumX[i] + int64(sortedX[i])
		prefixSumY[i+1] = prefixSumY[i] + int64(sortedY[i])
	}

	// Текущая позиция Кодеруна
	cx, cy := 0, 0
	results := make([]int64, 0, len(commands))

	// Обрабатываем каждую команду
	for _, cmd := range commands {
		// Обновляем позицию
		switch cmd {
		case 'N':
			cy++
		case 'S':
			cy--
		case 'E':
			cx++
		case 'W':
			cx--
		}

		// Вычисляем сумму манхэттенских расстояний
		// sum(|cx - mx| + |cy - my|) = sum(|cx - mx|) + sum(|cy - my|)

		// Для x координат
		idxX := sort.Search(N, func(i int) bool { return sortedX[i] > cx })
		// sortedX[0..idxX-1] <= cx, sortedX[idxX..N-1] > cx
		sumX := int64(cx)*int64(idxX) - prefixSumX[idxX] + (prefixSumX[N] - prefixSumX[idxX]) - int64(cx)*int64(N-idxX)

		// Для y координат
		idxY := sort.Search(N, func(i int) bool { return sortedY[i] > cy })
		// sortedY[0..idxY-1] <= cy, sortedY[idxY..N-1] > cy
		sumY := int64(cy)*int64(idxY) - prefixSumY[idxY] + (prefixSumY[N] - prefixSumY[idxY]) - int64(cy)*int64(N-idxY)

		total := sumX + sumY
		results = append(results, total)
	}

	return results
}
