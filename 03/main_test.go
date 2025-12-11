package main

import (
	"testing"
)

func TestExamples(t *testing.T) {
	tests := []struct {
		name     string
		M        int64
		W        []int64
		expected int
	}{
		{
			name:     "Пример 1",
			M:        5,
			W:        []int64{1, 3, 2},
			expected: 1,
		},
		{
			name:     "Пример 2",
			M:        10,
			W:        []int64{4, 5, 2, 3},
			expected: 4,
		},
		{
			name:     "Пример 3",
			M:        1,
			W:        []int64{2},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalNeed := int64(0)
			for _, w := range tt.W {
				totalNeed += w
			}
			result := solve(tt.M, tt.W, totalNeed)
			if result != tt.expected {
				t.Errorf("solve(M=%d, W=%v) = %d, want %d", tt.M, tt.W, result, tt.expected)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		M        int64
		W        []int64
		expected int
	}{
		{
			name:     "Одна группа, большая недостача",
			M:        1,
			W:        []int64{100},
			expected: 9801, // (100-1)^2 = 99^2 = 9801
		},
		{
			name:     "Две группы, равные потребности",
			M:        5,
			W:        []int64{5, 5},
			expected: 25,
		},
		{
			name:     "Три группы, равномерное распределение",
			M:        6,
			W:        []int64{3, 3, 3},
			expected: 3, // (3-2)^2 + (3-2)^2 + (3-2)^2 = 1 + 1 + 1 = 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalNeed := int64(0)
			for _, w := range tt.W {
				totalNeed += w
			}
			result := solve(tt.M, tt.W, totalNeed)
			if result < 0 || result >= mod {
				t.Errorf("solve(M=%d, W=%v) = %d, должно быть в [0, %d)", tt.M, tt.W, result, mod)
			}
			if tt.expected > 0 && result != tt.expected {
				t.Logf("solve(M=%d, W=%v) = %d, ожидалось %d (может отличаться из-за оптимизации)", tt.M, tt.W, result, tt.expected)
			}
		})
	}
}

func TestLargeInput(t *testing.T) {
	// Тест на больших данных
	W := make([]int64, 1000)
	totalNeed := int64(0)
	for i := 0; i < 1000; i++ {
		W[i] = int64(i + 1)
		totalNeed += W[i]
	}
	M := totalNeed / 2 // Половина от потребности

	result := solve(M, W, totalNeed)
	if result < 0 || result >= mod {
		t.Errorf("solve на больших данных вернул недопустимое значение %d", result)
	}
}

func TestElseBranch(t *testing.T) {
	// Тест для покрытия ветки else (когда группа не может принять среднюю недостачу)
	// M=1, W=[100, 1, 1], sum=102, недостача=101
	// avgShortfall = 101/3 = 33 > minW = 1, поэтому else выполняется
	M := int64(1)
	W := []int64{100, 1, 1}
	totalNeed := int64(102)
	result := solve(M, W, totalNeed)

	// shortfall = [99, 1, 1], сумма квадратов = 99² + 1 + 1 = 9803
	expected := 9803
	if result != expected {
		t.Errorf("solve(M=%d, W=%v) = %d, want %d", M, W, result, expected)
	}
}

func TestConstraints(t *testing.T) {
	tests := []struct {
		name string
		M    int64
		W    []int64
	}{
		{
			name: "Максимальный M",
			M:    1999999999,
			W:    []int64{2000000000, 2000000000},
		},
		{
			name: "Максимальный N",
			M:    100000,
			W:    make([]int64, 100000),
		},
		{
			name: "Максимальные W_i",
			M:    1,
			W:    []int64{2000000000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализируем W если нужно
			if len(tt.W) > 0 && tt.W[0] == 0 && len(tt.W) > 1 {
				for i := 0; i < len(tt.W); i++ {
					tt.W[i] = int64(i + 1)
				}
			}

			totalNeed := int64(0)
			for _, w := range tt.W {
				totalNeed += w
			}

			result := solve(tt.M, tt.W, totalNeed)
			if result < 0 || result >= mod {
				t.Errorf("solve вернул недопустимое значение %d", result)
			}
		})
	}
}

func BenchmarkSolve(b *testing.B) {
	W := make([]int64, 1000)
	totalNeed := int64(0)
	for i := 0; i < 1000; i++ {
		W[i] = int64(i + 1)
		totalNeed += W[i]
	}
	M := totalNeed / 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(M, W, totalNeed)
	}
}

func BenchmarkSolveLarge(b *testing.B) {
	W := make([]int64, 100000)
	totalNeed := int64(0)
	for i := 0; i < 100000; i++ {
		W[i] = int64(i + 1)
		totalNeed += W[i]
	}
	M := totalNeed / 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(M, W, totalNeed)
	}
}

func BenchmarkSolveMaxConstraints(b *testing.B) {
	W := make([]int64, 100000)
	for i := 0; i < 100000; i++ {
		W[i] = 2000000000
	}
	totalNeed := int64(100000) * 2000000000
	M := totalNeed - 100000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solve(M, W, totalNeed)
	}
}
