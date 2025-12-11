package main

import (
	"testing"
)

var (
	testFact    []int
	testInvFact []int
)

func init() {
	maxN := 400000
	testFact = precomputeFactorials(maxN)
	testInvFact = precomputeInvFactorials(testFact, maxN)
}

func TestExamples(t *testing.T) {
	tests := []struct {
		name     string
		n, s     int
		expected int
	}{
		{"n=2, s=3", 2, 3, 18},
		{"n=3, s=2", 3, 2, 0},
		{"n=3, s=3", 3, 3, 24},
		{"n=1, s=100", 1, 100, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.s, testFact, testInvFact)
			if result != tt.expected {
				t.Errorf("solve(n=%d, s=%d) = %d, expected %d",
					tt.n, tt.s, result, tt.expected)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		n, s     int
		expected int
	}{
		{"n > s", 10, 3, 0},
		{"n = s", 3, 3, 24},
		{"n = 1, s = 1", 1, 1, 2},
		{"n = 1, s = 2", 1, 2, 4},
		{"n = 2, s = 2", 2, 2, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.s, testFact, testInvFact)
			if result != tt.expected {
				t.Errorf("solve(n=%d, s=%d) = %d, expected %d",
					tt.n, tt.s, result, tt.expected)
			}
		})
	}
}

func TestFormula(t *testing.T) {
	// Проверяем формулу: answer(n, s) = (n+1)! × C(s, n)
	tests := []struct {
		n, s int
	}{
		{1, 5},
		{2, 5},
		{3, 5},
		{4, 5},
		{5, 5},
		{1, 10},
		{5, 10},
		{10, 10},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := solve(tt.n, tt.s, testFact, testInvFact)
			// (n+1)! × C(s, n)
			expected := testFact[tt.n+1] * comb(tt.s, tt.n, testFact, testInvFact) % mod
			if result != expected {
				t.Errorf("solve(n=%d, s=%d) = %d, expected %d",
					tt.n, tt.s, result, expected)
			}
		})
	}
}

func TestLargeValues(t *testing.T) {
	// Проверяем на больших значениях (не должно падать)
	tests := []struct {
		n, s int
	}{
		{100000, 200000},
		{200000, 200000},
		{1, 200000},
		{199999, 200000},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := solve(tt.n, tt.s, testFact, testInvFact)
			if result < 0 || result >= mod {
				t.Errorf("solve(n=%d, s=%d) = %d, expected 0 <= result < mod",
					tt.n, tt.s, result)
			}
		})
	}
}

// Тест производительности
func TestPerformance(t *testing.T) {
	// Симулируем T = 50000 запросов
	T := 50000
	for i := 0; i < T; i++ {
		n := (i % 200000) + 1
		s := 200000
		solve(n, s, testFact, testInvFact)
	}
}

// Бенчмарк одного запроса
func BenchmarkSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		solve(100000, 200000, testFact, testInvFact)
	}
}

// Бенчмарк предвычисления факториалов
func BenchmarkPrecompute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fact := precomputeFactorials(400000)
		precomputeInvFactorials(fact, 400000)
	}
}

// Бенчмарк 50000 запросов
func BenchmarkMultipleQueries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 50000; j++ {
			n := (j % 200000) + 1
			s := 200000
			solve(n, s, testFact, testInvFact)
		}
	}
}
