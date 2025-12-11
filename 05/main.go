package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
	idx     int
}

type Edge struct {
	from, to int
	weight   int
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем N
	line, _ := reader.ReadString('\n')
	N, _ := strconv.Atoi(strings.TrimSpace(line))

	// Читаем точки
	points := make([]Point, N)
	for i := 0; i < N; i++ {
		line, _ = reader.ReadString('\n')
		parts := strings.Fields(strings.TrimSpace(line))
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		points[i] = Point{x: x, y: y, z: z, idx: i}
	}

	result := solveMST(points)
	writer.WriteString(fmt.Sprintf("%d\n", result))
}

func find(parent []int, x int) int {
	if parent[x] != x {
		parent[x] = find(parent, parent[x])
	}
	return parent[x]
}

func union(parent, rank []int, x, y int) {
	if rank[x] < rank[y] {
		parent[x] = y
	} else if rank[x] > rank[y] {
		parent[y] = x
	} else {
		parent[y] = x
		rank[x]++
	}
}

// solveMST находит минимальное остовное дерево для заданных точек
func solveMST(points []Point) int64 {
	N := len(points)

	// Строим рёбра: для каждой координаты сортируем точки и добавляем рёбра между соседями
	edges := make([]Edge, 0, 3*N) // Максимум 3*(N-1) рёбер

	// Сортируем по x и добавляем рёбра между соседями
	sortedByX := make([]Point, N)
	copy(sortedByX, points)
	sort.Slice(sortedByX, func(i, j int) bool {
		return sortedByX[i].x < sortedByX[j].x
	})
	for i := 0; i < N-1; i++ {
		cost := min(
			abs(sortedByX[i].x-sortedByX[i+1].x),
			min(
				abs(sortedByX[i].y-sortedByX[i+1].y),
				abs(sortedByX[i].z-sortedByX[i+1].z),
			),
		)
		edges = append(edges, Edge{
			from:   sortedByX[i].idx,
			to:     sortedByX[i+1].idx,
			weight: cost,
		})
	}

	// Сортируем по y и добавляем рёбра между соседями
	sortedByY := make([]Point, N)
	copy(sortedByY, points)
	sort.Slice(sortedByY, func(i, j int) bool {
		return sortedByY[i].y < sortedByY[j].y
	})
	for i := 0; i < N-1; i++ {
		cost := min(
			abs(sortedByY[i].x-sortedByY[i+1].x),
			min(
				abs(sortedByY[i].y-sortedByY[i+1].y),
				abs(sortedByY[i].z-sortedByY[i+1].z),
			),
		)
		edges = append(edges, Edge{
			from:   sortedByY[i].idx,
			to:     sortedByY[i+1].idx,
			weight: cost,
		})
	}

	// Сортируем по z и добавляем рёбра между соседями
	sortedByZ := make([]Point, N)
	copy(sortedByZ, points)
	sort.Slice(sortedByZ, func(i, j int) bool {
		return sortedByZ[i].z < sortedByZ[j].z
	})
	for i := 0; i < N-1; i++ {
		cost := min(
			abs(sortedByZ[i].x-sortedByZ[i+1].x),
			min(
				abs(sortedByZ[i].y-sortedByZ[i+1].y),
				abs(sortedByZ[i].z-sortedByZ[i+1].z),
			),
		)
		edges = append(edges, Edge{
			from:   sortedByZ[i].idx,
			to:     sortedByZ[i+1].idx,
			weight: cost,
		})
	}

	// Сортируем рёбра по весу
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].weight < edges[j].weight
	})

	// Алгоритм Крускала с DSU
	parent := make([]int, N)
	rank := make([]int, N)
	for i := 0; i < N; i++ {
		parent[i] = i
		rank[i] = 0
	}

	var totalCost int64
	edgesUsed := 0

	for _, edge := range edges {
		if edgesUsed == N-1 {
			break
		}
		fromRoot := find(parent, edge.from)
		toRoot := find(parent, edge.to)
		if fromRoot != toRoot {
			union(parent, rank, fromRoot, toRoot)
			totalCost += int64(edge.weight)
			edgesUsed++
		}
	}

	return totalCost
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
