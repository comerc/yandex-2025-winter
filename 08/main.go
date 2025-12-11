package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const maxN = 700000

// Максимальное k, для которого существуют k-интересные числа <= maxN
// 2*3*4*5*6*7*8*9 = 362880 <= 700000
// 2*3*4*5*6*7*8*9*10 = 3628800 > 700000
// Значит maxK = 9
const maxK = 9

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Предвычисляем префиксные суммы k-интересных чисел
	prefixSums := precompute()

	// Читаем количество запросов
	line, _ := reader.ReadString('\n')
	q, _ := strconv.Atoi(strings.TrimSpace(line))

	// Обрабатываем запросы
	for i := 0; i < q; i++ {
		line, _ = reader.ReadString('\n')
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) < 3 {
			continue
		}
		k, _ := strconv.Atoi(parts[0])
		l, _ := strconv.Atoi(parts[1])
		r, _ := strconv.Atoi(parts[2])

		result := query(prefixSums, k, l, r)
		writer.WriteString(fmt.Sprintf("%d\n", result))
	}
}

// precompute генерирует все k-интересные числа и строит префиксные суммы
func precompute() [][]int32 {
	// prefixSums[k][n] = количество k-интересных чисел от 1 до n
	prefixSums := make([][]int32, maxK+1)
	for k := 1; k <= maxK; k++ {
		prefixSums[k] = make([]int32, maxN+1)
	}

	// Генерируем k-интересные числа для k >= 2 и сразу помечаем
	generate(1, 2, 0, maxN, prefixSums)

	// Все числа >= 2 являются 1-интересными
	for n := 2; n <= maxN; n++ {
		prefixSums[1][n] = 1
	}

	// Преобразуем в префиксные суммы
	for k := 1; k <= maxK; k++ {
		for n := 1; n <= maxN; n++ {
			prefixSums[k][n] += prefixSums[k][n-1]
		}
	}

	return prefixSums
}

// generate рекурсивно генерирует все произведения возрастающих последовательностей множителей
func generate(product, minFactor, depth, limit int, prefixSums [][]int32) {
	for factor := minFactor; product*factor <= limit; factor++ {
		newProduct := product * factor
		newDepth := depth + 1

		if newDepth >= 2 && newDepth <= maxK {
			prefixSums[newDepth][newProduct] = 1
		}

		generate(newProduct, factor+1, newDepth, limit, prefixSums)
	}
}

// query возвращает количество k-интересных чисел в диапазоне [l, r]
func query(prefixSums [][]int32, k, l, r int) int32 {
	if k > maxK || l > maxN {
		return 0
	}
	if r > maxN {
		r = maxN
	}
	if l < 1 {
		l = 1
	}
	return prefixSums[k][r] - prefixSums[k][l-1]
}
