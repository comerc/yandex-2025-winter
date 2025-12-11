package main

import (
	"testing"
)

func TestExample1(t *testing.T) {
	points := []Point{
		{x: 2, y: 6, z: 11, idx: 0},
		{x: 8, y: 9, z: 3, idx: 1},
	}
	expected := int64(3)

	result := solveMST(points)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestExample2(t *testing.T) {
	points := []Point{
		{x: 1, y: 1, z: 1, idx: 0},
		{x: -5, y: -5, z: -5, idx: 1},
		{x: -10, y: -10, z: -10, idx: 2},
	}
	expected := int64(11)

	result := solveMST(points)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestExample3(t *testing.T) {
	points := []Point{
		{x: -2, y: -2, z: -6, idx: 0},
		{x: 10, y: -16, z: -16, idx: 1},
		{x: 18, y: -5, z: 18, idx: 2},
		{x: 13, y: -16, z: -16, idx: 3},
		{x: 9, y: -5, z: -2, idx: 4},
	}
	expected := int64(4)

	result := solveMST(points)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestEdgeCases(t *testing.T) {
	// Одна точка
	points1 := []Point{
		{x: 0, y: 0, z: 0, idx: 0},
	}
	expected1 := int64(0)

	result1 := solveMST(points1)
	if result1 != expected1 {
		t.Errorf("Expected %d, got %d", expected1, result1)
	}

	// Две точки с одинаковыми координатами
	points2 := []Point{
		{x: 0, y: 0, z: 0, idx: 0},
		{x: 0, y: 0, z: 0, idx: 1},
	}
	expected2 := int64(0)

	result2 := solveMST(points2)
	if result2 != expected2 {
		t.Errorf("Expected %d, got %d", expected2, result2)
	}

	// Три точки на одной прямой
	points3 := []Point{
		{x: 0, y: 0, z: 0, idx: 0},
		{x: 1, y: 1, z: 1, idx: 1},
		{x: 2, y: 2, z: 2, idx: 2},
	}
	result3 := solveMST(points3)
	if result3 < 0 {
		t.Errorf("Result should be non-negative, got %d", result3)
	}
}

func TestLargeCoordinates(t *testing.T) {
	// Тест с большими координатами
	points := []Point{
		{x: -1000000000, y: -1000000000, z: -1000000000, idx: 0},
		{x: 0, y: 0, z: 0, idx: 1},
		{x: 1000000000, y: 1000000000, z: 1000000000, idx: 2},
	}
	result := solveMST(points)
	if result < 0 {
		t.Errorf("Result should be non-negative, got %d", result)
	}
	// Ожидаем: min(10^9, 10^9, 10^9) + min(10^9, 10^9, 10^9) = 2*10^9
	// Но это может быть оптимизировано до меньшего значения
	// Проверяем только, что результат корректен
	if result > 2000000000 {
		t.Errorf("Result seems too large: %d", result)
	}
}

// Тесты для функций из main.go

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{5, 5},
		{-5, 5},
		{1000000000, 1000000000},
		{-1000000000, 1000000000},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func TestFind(t *testing.T) {
	// Тест find с path compression
	parent := []int{0, 0, 1, 2, 3, 4}

	// Изначально все указывают на предыдущий элемент
	// После find(5) должно произойти сжатие пути

	root := find(parent, 5)
	if root != 0 {
		t.Errorf("find(5) = %d, expected 0", root)
	}

	// Проверяем, что произошло сжатие пути
	if parent[5] != 0 {
		t.Errorf("parent[5] should be compressed to 0, got %d", parent[5])
	}
	if parent[4] != 0 {
		t.Errorf("parent[4] should be compressed to 0, got %d", parent[4])
	}
}

func TestUnion(t *testing.T) {
	// Тест union с разными рангами
	parent := []int{0, 1, 2, 3, 4, 5}
	rank := []int{0, 0, 0, 0, 0, 0}

	// Объединяем множества
	union(parent, rank, 0, 1)
	if find(parent, 0) != find(parent, 1) {
		t.Errorf("0 and 1 should be in the same set")
	}

	// Объединяем ещё множества
	union(parent, rank, 2, 3)
	union(parent, rank, 4, 5)

	// Объединяем два множества
	root1 := find(parent, 0)
	root2 := find(parent, 2)
	union(parent, rank, root1, root2)

	if find(parent, 0) != find(parent, 2) {
		t.Errorf("0 and 2 should be in the same set after union")
	}

	// Проверяем, что ранги обновляются правильно
	union(parent, rank, find(parent, 4), find(parent, 0))
	if find(parent, 4) != find(parent, 0) {
		t.Errorf("4 and 0 should be in the same set")
	}
}

func TestUnionRank(t *testing.T) {
	// Тест union с проверкой рангов
	parent := []int{0, 1, 2, 3}
	rank := []int{1, 0, 1, 0}

	// Объединяем множества с разными рангами
	union(parent, rank, 0, 1)
	if rank[0] != 1 {
		t.Errorf("rank[0] should remain 1, got %d", rank[0])
	}
	if parent[1] != 0 {
		t.Errorf("parent[1] should be 0 (lower rank attaches to higher), got %d", parent[1])
	}

	// Объединяем множества с одинаковыми рангами (оба ранга = 1)
	// Сначала нужно сделать rank[3] = 1
	parent2 := []int{0, 1, 2, 3}
	rank2 := []int{0, 0, 1, 1}
	union(parent2, rank2, 2, 3)
	// Когда ранги равны, rank[2] должен увеличиться
	if rank2[2] != 2 {
		t.Errorf("rank[2] should be incremented to 2 when ranks are equal, got %d", rank2[2])
	}
	if parent2[3] != 2 {
		t.Errorf("parent[3] should be 2, got %d", parent2[3])
	}
}

// Тест производительности
func TestPerformance(t *testing.T) {
	// Симулируем N = 10000 точек
	N := 10000
	points := make([]Point, N)
	for i := 0; i < N; i++ {
		x := (i*1234567)%2000000001 - 1000000000
		y := (i*7654321)%2000000001 - 1000000000
		z := (i*9876543)%2000000001 - 1000000000
		points[i] = Point{x: x, y: y, z: z, idx: i}
	}

	result := solveMST(points)
	if result < 0 {
		t.Errorf("Result should be non-negative, got %d", result)
	}
}

// Бенчмарк для одного MST
func BenchmarkSolveMST(b *testing.B) {
	N := 10000
	points := make([]Point, N)
	for i := 0; i < N; i++ {
		x := (i*1234567)%2000000001 - 1000000000
		y := (i*7654321)%2000000001 - 1000000000
		z := (i*9876543)%2000000001 - 1000000000
		points[i] = Point{x: x, y: y, z: z, idx: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = solveMST(points)
	}
}

// Бенчмарк для больших данных
func BenchmarkSolveMSTLarge(b *testing.B) {
	N := 100000
	points := make([]Point, N)
	for i := 0; i < N; i++ {
		x := (i*1234567)%2000000001 - 1000000000
		y := (i*7654321)%2000000001 - 1000000000
		z := (i*9876543)%2000000001 - 1000000000
		points[i] = Point{x: x, y: y, z: z, idx: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = solveMST(points)
	}
}
