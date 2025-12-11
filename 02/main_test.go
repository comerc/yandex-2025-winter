package main

import (
	"testing"
)

func TestExamples(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "Пример 1",
			nums:     []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			expected: 100,
		},
		{
			name:     "Пример 2 (числа Фибоначчи)",
			nums:     []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89},
			expected: 100,
		},
		{
			name:     "Пример 3 (все по 40)",
			nums:     []int{40, 40, 40, 40, 40, 40, 40, 40, 40, 40},
			expected: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.nums)
			if result != tt.expected {
				t.Errorf("solve(%v) = %d, want %d", tt.nums, result, tt.expected)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "Все единицы",
			nums:     []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expected: 10, // сумма всех = 10, ближе к 100 чем 0
		},
		{
			name:     "Все сотни",
			nums:     []int{100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
			expected: 100, // одна сотня = 100
		},
		{
			name:     "Точно 100",
			nums:     []int{100, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expected: 100,
		},
		{
			name:     "98 vs 102 - выбираем 102",
			nums:     []int{51, 51, 1, 1, 1, 1, 1, 1, 1, 1},
			expected: 102, // 51 + 51 = 102
		},
		{
			name:     "Пустое подмножество vs полное",
			nums:     []int{50, 50, 50, 50, 50, 50, 50, 50, 50, 50},
			expected: 100, // 2 * 50 = 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.nums)
			if result != tt.expected {
				t.Errorf("solve(%v) = %d, want %d", tt.nums, result, tt.expected)
			}
		})
	}
}

func TestTieBreaker(t *testing.T) {
	// Проверяем, что при равном расстоянии выбирается большая сумма
	// Пример из условия: для сумм 98 и 102 нужно вывести 102
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "99 vs 101 - выбираем 101",
			nums:     []int{99, 2, 3, 3, 3, 3, 3, 3, 3, 3},
			expected: 101, // 99 + 2 = 101, 99 - тоже вариант, но 101 > 99
		},
		{
			name:     "Пример 3 из условия - 80 vs 120",
			nums:     []int{40, 40, 40, 40, 40, 40, 40, 40, 40, 40},
			expected: 120, // 3 × 40 = 120 ближе к 100 чем 2 × 40 = 80, и 120 > 80
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.nums)
			if result != tt.expected {
				t.Errorf("solve(%v) = %d, want %d", tt.nums, result, tt.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
		{100, 100},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}

func BenchmarkSolve(b *testing.B) {
	nums := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(nums)
	}
}

func BenchmarkSolveWorstCase(b *testing.B) {
	nums := []int{100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(nums)
	}
}
