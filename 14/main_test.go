package main

import (
	"testing"
)

// Тест на пример из условия задачи
func TestSolveExample(t *testing.T) {
	testCases := []struct {
		name     string
		a        int64
		q        int64
		L        int64
		R        int64
		expected int64
	}{
		{"Пример из условия: a=1, q=2, L=0, R=3", 1, 2, 0, 3, 104},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.a, tc.q, tc.L, tc.R)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d) = %d, ожидалось %d",
					tc.a, tc.q, tc.L, tc.R, actual, tc.expected)
			}
		})
	}
}

// Тест на специальные случаи
func TestSolveSpecialCases(t *testing.T) {
	testCases := []struct {
		name     string
		a        int64
		q        int64
		L        int64
		R        int64
		expected int64
	}{
		{"a=0", 0, 2, 0, 3, 0},
		{"q=0, L=0, R=3", 1, 0, 0, 3, 96},
		{"q=0, L=-1, R=2", 1, 0, -1, 2, 96},
		{"q=1", 1, 1, 0, 3, 0},
		{"q=-1, L=0, R=3", 1, -1, 0, 3, 128},
		{"q=-1, L=1, R=4", 1, -1, 1, 4, 128},
		{"N <= 0", 1, 2, 5, 3, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.a, tc.q, tc.L, tc.R)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d) = %d, ожидалось %d",
					tc.a, tc.q, tc.L, tc.R, actual, tc.expected)
			}
		})
	}
}

// Тест на найденные ошибки (правильные значения из 14b/main.go)
func TestFoundBugs(t *testing.T) {
	testCases := []struct {
		name     string
		a        int64
		q        int64
		L        int64
		R        int64
		expected int64
	}{
		// q=20, L=5, R=14 - правильный ответ 2232 для всех a ≠ 0
		{"q=20, L=5, R=14, a=7", 7, 20, 5, 14, 2232},
		{"q=20, L=5, R=14, a=10", 10, 20, 5, 14, 2232},
		{"q=20, L=5, R=14, a=1", 1, 20, 5, 14, 2232},
		{"q=20, L=5, R=14, a=-1", -1, 20, 5, 14, 2232},
		{"q=20, L=5, R=14, a=0", 0, 20, 5, 14, 0},

		// Другие случаи с правильными значениями из 14b/main.go
		{"a=15, q=47, L=9, R=15", 15, 47, 9, 15, 690},
		{"a=-11, q=39, L=6, R=17", -11, 39, 6, 17, 4048},
		{"a=13, q=33, L=5, R=16", 13, 33, 5, 16, 4048},
		{"a=-3, q=22, L=10, R=17", -3, 22, 10, 17, 1072},
		{"a=-7, q=38, L=5, R=14", -7, 38, 5, 14, 2232},
		{"a=11, q=49, L=9, R=20", 11, 49, 9, 20, 4048},
		{"a=-1, q=30, L=7, R=14", -1, 30, 7, 14, 1072},
		{"a=14, q=50, L=8, R=18", 14, 50, 8, 18, 3050},
		{"a=13, q=11, L=8, R=18", 13, 11, 8, 18, 3050},
		{"a=-2, q=10, L=9, R=20", -2, 10, 9, 20, 4048},
		{"a=12, q=20, L=8, R=14", 12, 20, 8, 14, 690},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.a, tc.q, tc.L, tc.R)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d) = %d, ожидалось %d",
					tc.a, tc.q, tc.L, tc.R, actual, tc.expected)
			}
		})
	}
}

// Тест на различные диапазоны и значения q
func TestSolveVariousRanges(t *testing.T) {
	testCases := []struct {
		name     string
		a        int64
		q        int64
		L        int64
		R        int64
		expected int64
	}{
		{"q=2, L=0, R=2", 1, 2, 0, 2, 38},
		{"q=2, L=0, R=4", 1, 2, 0, 4, 224},
		{"q=3, L=0, R=3", 1, 3, 0, 3, 104},
		{"q=4, L=0, R=3", 1, 4, 0, 3, 104},
		{"q=5, L=0, R=3", 1, 5, 0, 3, 104},
		{"q=10, L=0, R=3", 1, 10, 0, 3, 104},
		{"q=2, L=1, R=4", 1, 2, 1, 4, 104},
		{"q=2, L=5, R=8", 1, 2, 5, 8, 104},
		{"a=2, q=2, L=0, R=3", 2, 2, 0, 3, 104},
		{"a=-2, q=2, L=0, R=3", -2, 2, 0, 3, 104},
		{"a=5, q=2, L=0, R=3", 5, 2, 0, 3, 104},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.a, tc.q, tc.L, tc.R)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d) = %d, ожидалось %d",
					tc.a, tc.q, tc.L, tc.R, actual, tc.expected)
			}
		})
	}
}

// Тест на отрицательные индексы
func TestSolveNegativeIndices(t *testing.T) {
	testCases := []struct {
		name     string
		a        int64
		q        int64
		L        int64
		R        int64
		expected int64
	}{
		{"q=2, L=-2, R=2", 1, 2, -2, 2, 224},
		{"q=3, L=-2, R=2", 1, 3, -2, 2, 224},
		{"q=2, L=-1, R=1", 1, 2, -1, 1, 38},
		{"q=2, L=-3, R=0", 1, 2, -3, 0, 104},
		{"q=2, L=-2, R=0", 1, 2, -2, 0, 38},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.a, tc.q, tc.L, tc.R)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d) = %d, ожидалось %d",
					tc.a, tc.q, tc.L, tc.R, actual, tc.expected)
			}
		})
	}
}

// Тест на большие q
func TestSolveLargeQ(t *testing.T) {
	testCases := []struct {
		name     string
		a        int64
		q        int64
		L        int64
		R        int64
		expected int64
	}{
		{"q=50, L=0, R=5", 1, 50, 0, 5, 412},
		{"q=100, L=0, R=5", 1, 100, 0, 5, 412},
		{"q=200, L=0, R=5", 1, 200, 0, 5, 412},
		{"q=500, L=0, R=5", 1, 500, 0, 5, 412},
		{"q=1000, L=0, R=5", 1, 1000, 0, 5, 412},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.a, tc.q, tc.L, tc.R)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d) = %d, ожидалось %d",
					tc.a, tc.q, tc.L, tc.R, actual, tc.expected)
			}
		})
	}
}
