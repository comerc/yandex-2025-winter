package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем количество тестов
	line, _ := reader.ReadString('\n')
	t, _ := strconv.Atoi(strings.TrimSpace(line))

	// Используем scanner для чтения чисел
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	readInt := func() int {
		scanner.Scan()
		v, _ := strconv.Atoi(scanner.Text())
		return v
	}

	for test := 0; test < t; test++ {
		// Читаем n и m
		n := readInt()
		m := readInt()

		// Читаем массивы a, b, c
		a := make([]int, m)
		b := make([]int, m)
		c := make([]int, m)

		for i := 0; i < m; i++ {
			a[i] = readInt()
		}
		for i := 0; i < m; i++ {
			b[i] = readInt()
		}
		for i := 0; i < m; i++ {
			c[i] = readInt()
		}

		edges := make([]Edge, m)
		for i := 0; i < m; i++ {
			edges[i] = Edge{a: a[i] - 1, b: b[i] - 1, w: c[i]} // 0-indexed
		}

		result := solve(n, edges)

		// Выводим результат
		sb := strings.Builder{}
		for i, v := range result {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		writer.WriteString(sb.String())
	}
}

type Edge struct {
	a, b, w int
}

// solve находит минимальный вес w для каждого размера компоненты k
func solve(n int, edges []Edge) []int {
	// Сортируем рёбра по весу
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})

	// Инициализируем DSU
	parent := make([]int, n)
	rank := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}

	// result[k] = минимальный вес для размера k+1
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = -1
	}

	// k=1: любая вершина достижима сама из себя с w=0
	result[0] = 0

	// maxReached отслеживает максимальный размер, для которого уже найден ответ
	maxReached := 1

	// Обрабатываем рёбра в порядке возрастания веса
	for _, e := range edges {
		// Пропускаем петли (они не меняют связность)
		if e.a == e.b {
			continue
		}

		rootA := find(parent, e.a)
		rootB := find(parent, e.b)

		if rootA != rootB {
			// Объединяем компоненты
			newSize := size[rootA] + size[rootB]
			union(parent, rank, size, rootA, rootB)

			// Обновляем результат для всех новых размеров
			for k := maxReached + 1; k <= newSize; k++ {
				result[k-1] = e.w
			}
			if newSize > maxReached {
				maxReached = newSize
			}
		}

		// Если достигли максимального размера, можно остановиться
		if maxReached == n {
			break
		}
	}

	return result
}

func find(parent []int, x int) int {
	if parent[x] != x {
		parent[x] = find(parent, parent[x])
	}
	return parent[x]
}

func union(parent, rank, size []int, x, y int) {
	if rank[x] < rank[y] {
		parent[x] = y
		size[y] += size[x]
	} else if rank[x] > rank[y] {
		parent[y] = x
		size[x] += size[y]
	} else {
		parent[y] = x
		size[x] += size[y]
		rank[x]++
	}
}
