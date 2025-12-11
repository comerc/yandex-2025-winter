package main

import (
	"testing"
)

func TestExamples(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		edges    []Edge
		expected []int
	}{
		{
			name:     "Пример 1: одна вершина с петлёй",
			n:        1,
			edges:    []Edge{{a: 0, b: 0, w: 1}},
			expected: []int{0},
		},
		{
			name: "Пример 2: 4 вершины, 3 ребра",
			n:    4,
			// a = [1, 2, 1], b = [2, 2, 2], c = [4, 2, 3]
			// Ребро 1: 1-2 вес 4
			// Ребро 2: 2-2 вес 2 (петля)
			// Ребро 3: 1-2 вес 3
			edges: []Edge{
				{a: 0, b: 1, w: 4}, // 1-2 вес 4
				{a: 1, b: 1, w: 2}, // петля 2-2 вес 2
				{a: 0, b: 1, w: 3}, // 1-2 вес 3
			},
			expected: []int{0, 3, -1, -1},
		},
		{
			name: "Пример 3: 5 вершин, 5 рёбер",
			n:    5,
			edges: []Edge{
				{a: 0, b: 1, w: 1}, // 1-2 вес 1
				{a: 1, b: 0, w: 2}, // 2-1 вес 2
				{a: 2, b: 3, w: 2}, // 3-4 вес 2
				{a: 3, b: 2, w: 1}, // 4-3 вес 1
				{a: 1, b: 2, w: 5}, // 2-3 вес 5
			},
			expected: []int{0, 1, 5, 5, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.edges)
			if len(result) != len(tt.expected) {
				t.Fatalf("solve() returned %d elements, want %d", len(result), len(tt.expected))
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("solve()[%d] = %d, want %d", i, result[i], v)
				}
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		edges    []Edge
		expected []int
	}{
		{
			name:     "Без рёбер",
			n:        3,
			edges:    []Edge{},
			expected: []int{0, -1, -1},
		},
		{
			name:     "Только петли",
			n:        2,
			edges:    []Edge{{a: 0, b: 0, w: 1}, {a: 1, b: 1, w: 2}},
			expected: []int{0, -1},
		},
		{
			name:     "Полный граф из 3 вершин",
			n:        3,
			edges:    []Edge{{a: 0, b: 1, w: 1}, {a: 1, b: 2, w: 2}, {a: 0, b: 2, w: 3}},
			expected: []int{0, 1, 2},
		},
		{
			name:     "Цепочка",
			n:        4,
			edges:    []Edge{{a: 0, b: 1, w: 1}, {a: 1, b: 2, w: 2}, {a: 2, b: 3, w: 3}},
			expected: []int{0, 1, 2, 3},
		},
		{
			name:     "Две компоненты",
			n:        4,
			edges:    []Edge{{a: 0, b: 1, w: 1}, {a: 2, b: 3, w: 2}},
			expected: []int{0, 1, -1, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.edges)
			if len(result) != len(tt.expected) {
				t.Fatalf("solve() returned %d elements, want %d", len(result), len(tt.expected))
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("solve()[%d] = %d, want %d", i, result[i], v)
				}
			}
		})
	}
}

func TestDSU(t *testing.T) {
	// Тест функций find и union
	parent := []int{0, 1, 2, 3, 4}
	rank := []int{0, 0, 0, 0, 0}
	size := []int{1, 1, 1, 1, 1}

	// Проверяем find
	for i := 0; i < 5; i++ {
		if find(parent, i) != i {
			t.Errorf("find(%d) = %d, want %d", i, find(parent, i), i)
		}
	}

	// Объединяем 0 и 1
	union(parent, rank, size, 0, 1)
	root01 := find(parent, 0)
	if find(parent, 1) != root01 {
		t.Errorf("После union(0,1): find(0)=%d, find(1)=%d, должны быть равны", find(parent, 0), find(parent, 1))
	}

	// Объединяем 2 и 3
	union(parent, rank, size, 2, 3)
	root23 := find(parent, 2)
	if find(parent, 3) != root23 {
		t.Errorf("После union(2,3): find(2)=%d, find(3)=%d, должны быть равны", find(parent, 2), find(parent, 3))
	}

	// Объединяем компоненты
	union(parent, rank, size, root01, root23)
	rootAll := find(parent, 0)
	for i := 0; i < 4; i++ {
		if find(parent, i) != rootAll {
			t.Errorf("После объединения всех: find(%d)=%d, want %d", i, find(parent, i), rootAll)
		}
	}
}

func TestLargeGraph(t *testing.T) {
	// Тест на большом графе
	n := 1000
	edges := make([]Edge, n-1)
	for i := 0; i < n-1; i++ {
		edges[i] = Edge{a: i, b: i + 1, w: i + 1}
	}

	result := solve(n, edges)

	// k=1: w=0
	if result[0] != 0 {
		t.Errorf("result[0] = %d, want 0", result[0])
	}

	// k=2: w=1 (первое ребро)
	if result[1] != 1 {
		t.Errorf("result[1] = %d, want 1", result[1])
	}

	// k=n: w=n-1 (последнее ребро)
	if result[n-1] != n-1 {
		t.Errorf("result[%d] = %d, want %d", n-1, result[n-1], n-1)
	}
}

func BenchmarkSolve(b *testing.B) {
	n := 1000
	edges := make([]Edge, n-1)
	for i := 0; i < n-1; i++ {
		edges[i] = Edge{a: i, b: i + 1, w: i + 1}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(n, edges)
	}
}

func BenchmarkSolveLarge(b *testing.B) {
	n := 10000
	edges := make([]Edge, n*2)
	for i := 0; i < n-1; i++ {
		edges[i] = Edge{a: i, b: i + 1, w: i + 1}
	}
	// Добавляем случайные рёбра
	for i := n - 1; i < n*2; i++ {
		edges[i] = Edge{a: i % n, b: (i * 7) % n, w: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(n, edges)
	}
}
