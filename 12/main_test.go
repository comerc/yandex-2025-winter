package main

import (
	"runtime"
	"testing"
	"time"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		good     map[int]bool
		expected int
	}{
		{
			name:     "Пример 1: n=3, хорошее число 15",
			n:        3,
			good:     map[int]bool{15: true},
			expected: 69,
		},
		{
			name:     "Пример 2: n=5, хорошие числа 27, 10, 5, 7",
			n:        5,
			good:     map[int]bool{27: true, 10: true, 5: true, 7: true},
			expected: 518,
		},
		{
			name:     "n=3, все суммы хорошие (0-27)",
			n:        3,
			good:     makeAllGood(0, 27),
			expected: 900, // 9 * 10 * 10 = 900 (первая цифра 1-9, остальные 0-9)
		},
		{
			name:     "n=3, только сумма 0 хорошая",
			n:        3,
			good:     map[int]bool{0: true},
			expected: 0, // только число 000, но первая цифра не может быть 0, так что 0
		},
		{
			name:     "n=4, хорошие числа 0-27",
			n:        4,
			good:     makeAllGood(0, 27),
			expected: 9000, // 9 * 10 * 10 * 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.good)
			if result != tt.expected {
				t.Errorf("solve(n=%d, good=%v) = %d, ожидалось %d",
					tt.n, tt.good, result, tt.expected)
			}
		})
	}
}

func TestSolveEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		good     map[int]bool
		expected int
	}{
		{
			name:     "n=2 (граничное значение, должно вернуть 0)",
			n:        2,
			good:     map[int]bool{15: true},
			expected: 0, // по условию n >= 3, для n < 3 возвращаем 0
		},
		{
			name:     "n=3, только сумма 27 хорошая",
			n:        3,
			good:     map[int]bool{27: true},
			expected: 1, // только число 999
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.good)
			if result != tt.expected {
				t.Errorf("solve(n=%d, good=%v) = %d, ожидалось %d",
					tt.n, tt.good, result, tt.expected)
			}
		})
	}
}

func TestSolveLargeN(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		good map[int]bool
	}{
		{
			name: "n=10^12, хорошее число 15",
			n:    1000000000000,
			good: map[int]bool{15: true},
		},
		{
			name: "n=10^9, хорошие числа 0-27",
			n:    1000000000,
			good: makeAllGood(0, 27),
		},
		{
			name: "n=1000, хорошие числа 10, 15, 20",
			n:    1000,
			good: map[int]bool{10: true, 15: true, 20: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			result := solve(tt.n, tt.good)
			elapsed := time.Since(start)

			// Проверяем ограничение времени: 2 секунды
			maxTime := 2 * time.Second
			if elapsed > maxTime {
				t.Errorf("solve(n=%d) выполнилось за %v, что превышает ограничение %v",
					tt.n, elapsed, maxTime)
			}

			// Проверяем, что результат в допустимом диапазоне [0, mod)
			if result < 0 || result >= mod {
				t.Errorf("solve(n=%d) вернул %d, что вне диапазона [0, %d)",
					tt.n, result, mod)
			}

			t.Logf("solve(n=%d) = %d, время выполнения: %v", tt.n, result, elapsed)
		})
	}
}

func TestSolveConstraints(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		good map[int]bool
	}{
		{
			name: "Минимальное n=3",
			n:    3,
			good: map[int]bool{15: true},
		},
		{
			name: "Максимальное n=10^12",
			n:    1000000000000,
			good: map[int]bool{15: true},
		},
		{
			name: "Минимальное количество хороших чисел (m=1)",
			n:    100,
			good: map[int]bool{15: true},
		},
		{
			name: "Максимальное количество хороших чисел (m=28)",
			n:    100,
			good: makeAllGood(0, 27),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			result := solve(tt.n, tt.good)
			elapsed := time.Since(start)

			// Проверяем ограничение времени: 2 секунды
			maxTime := 2 * time.Second
			if elapsed > maxTime {
				t.Errorf("solve(n=%d) выполнилось за %v, что превышает ограничение %v",
					tt.n, elapsed, maxTime)
			}

			// Проверяем, что результат в допустимом диапазоне
			if result < 0 || result >= mod {
				t.Errorf("solve(n=%d) вернул %d, что вне диапазона [0, %d)",
					tt.n, result, mod)
			}

			t.Logf("solve(n=%d, |good|=%d) = %d, время: %v", tt.n, len(tt.good), result, elapsed)
		})
	}
}

func TestSolveMemoryUsage(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		good map[int]bool
	}{
		{
			name: "n=3",
			n:    3,
			good: map[int]bool{15: true},
		},
		{
			name: "n=10^12",
			n:    1000000000000,
			good: map[int]bool{15: true},
		},
		{
			name: "n=10^9",
			n:    1000000000,
			good: makeAllGood(0, 27),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем память до выполнения
			var m1, m2 runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&m1)

			result := solve(tt.n, tt.good)

			// Измеряем память после выполнения
			runtime.GC()
			runtime.ReadMemStats(&m2)

			// Вычисляем использованную память в байтах
			allocated := m2.TotalAlloc - m1.TotalAlloc
			maxMemory := 256 * 1024 * 1024 // 256 МБ в байтах

			// Проверяем ограничение памяти: 256 МБ
			if allocated > uint64(maxMemory) {
				t.Errorf("solve(n=%d) использовало %d байт памяти, что превышает ограничение %d байт (256 МБ)",
					tt.n, allocated, maxMemory)
			}

			// Проверяем, что результат в допустимом диапазоне
			if result < 0 || result >= mod {
				t.Errorf("solve(n=%d) вернул %d, что вне диапазона [0, %d)",
					tt.n, result, mod)
			}

			t.Logf("solve(n=%d) = %d, использовано памяти: %d байт (%.2f МБ)",
				tt.n, result, allocated, float64(allocated)/(1024*1024))
		})
	}
}

func BenchmarkSolve(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int64
		good map[int]bool
	}{
		{"n=3", 3, map[int]bool{15: true}},
		{"n=100", 100, map[int]bool{15: true}},
		{"n=1000", 1000, map[int]bool{15: true}},
		{"n=10^6", 1000000, map[int]bool{15: true}},
		{"n=10^9", 1000000000, map[int]bool{15: true}},
		{"n=10^12", 1000000000000, map[int]bool{15: true}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				solve(bm.n, bm.good)
			}
		})
	}
}

func BenchmarkSolveAllGood(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int64
		good map[int]bool
	}{
		{"n=3, все хорошие", 3, makeAllGood(0, 27)},
		{"n=100, все хорошие", 100, makeAllGood(0, 27)},
		{"n=1000, все хорошие", 1000, makeAllGood(0, 27)},
		{"n=10^6, все хорошие", 1000000, makeAllGood(0, 27)},
		{"n=10^9, все хорошие", 1000000000, makeAllGood(0, 27)},
		{"n=10^12, все хорошие", 1000000000000, makeAllGood(0, 27)},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				solve(bm.n, bm.good)
			}
		})
	}
}

// makeAllGood создает множество всех хороших чисел от min до max включительно
func makeAllGood(min, max int) map[int]bool {
	good := make(map[int]bool)
	for i := min; i <= max; i++ {
		good[i] = true
	}
	return good
}
