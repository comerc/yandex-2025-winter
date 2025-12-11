package main

import (
	"runtime"
	"testing"
	"time"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int64
	}{
		{
			name:     "n=1 (базовый случай)",
			n:        1,
			expected: 1,
		},
		{
			name:     "Пример 1: n=2",
			n:        2,
			expected: 4,
		},
		{
			name:     "n=3",
			n:        3,
			expected: 10,
		},
		{
			name:     "Пример 2: n=4",
			n:        4,
			expected: 333333357,
		},
		{
			name:     "n=5",
			n:        5,
			expected: 666666714, // E[5] = 64/3 mod (10^9+7)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n)
			if result != tt.expected {
				t.Errorf("solve(%d) = %d, ожидалось %d", tt.n, result, tt.expected)
			}
		})
	}
}

func TestSolveSmallValues(t *testing.T) {
	// Проверяем, что все значения в допустимом диапазоне [0, mod)
	for n := 1; n <= 100; n++ {
		result := solve(n)
		if result < 0 || result >= mod {
			t.Errorf("solve(%d) = %d, выходит за пределы [0, %d)", n, result, mod)
		}
	}
}

func TestSolveConstraints(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "Минимальное значение n=1",
			n:    1,
		},
		{
			name: "Малое значение n=10",
			n:    10,
		},
		{
			name: "Среднее значение n=1000",
			n:    1000,
		},
		{
			name: "Большое значение n=10000",
			n:    10000,
		},
		{
			name: "Очень большое значение n=100000",
			n:    100000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n)

			// Проверяем, что результат неотрицательный и меньше mod
			if result < 0 || result >= mod {
				t.Errorf("solve(%d) вернул недопустимое значение %d (должно быть в [0, %d))",
					tt.n, result, mod)
			}

			// Проверяем ограничения из условия задачи
			if tt.n < 1 || tt.n > 2000000 {
				t.Errorf("n=%d не соответствует ограничению 1 <= n <= 2*10^6", tt.n)
			}
		})
	}
}

func TestSolvePerformance(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "Минимальное значение",
			n:    1,
		},
		{
			name: "Малое значение",
			n:    100,
		},
		{
			name: "Среднее значение",
			n:    10000,
		},
		{
			name: "Большое значение",
			n:    100000,
		},
		{
			name: "Очень большое значение",
			n:    1000000,
		},
		{
			name: "Максимальное значение",
			n:    2000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем время выполнения
			start := time.Now()
			result := solve(tt.n)
			elapsed := time.Since(start)

			// Проверяем ограничение времени: 1 с
			maxTime := 1 * time.Second
			if elapsed > maxTime {
				t.Errorf("solve(%d) выполнилось за %v, что превышает ограничение %v",
					tt.n, elapsed, maxTime)
			}

			// Проверяем, что результат валидный
			if result < 0 || result >= mod {
				t.Errorf("solve(%d) вернул недопустимое значение %d", tt.n, result)
			}

			t.Logf("solve(%d) выполнилось за %v, результат: %d", tt.n, elapsed, result)
		})
	}
}

func TestSolveMemoryUsage(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "Минимальное значение",
			n:    1,
		},
		{
			name: "Малое значение",
			n:    1000,
		},
		{
			name: "Среднее значение",
			n:    100000,
		},
		{
			name: "Большое значение",
			n:    1000000,
		},
		{
			name: "Максимальное значение",
			n:    2000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем память до выполнения
			var m1, m2 runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&m1)

			result := solve(tt.n)

			// Измеряем память после выполнения
			runtime.GC()
			runtime.ReadMemStats(&m2)

			// Вычисляем использованную память в байтах
			allocated := m2.TotalAlloc - m1.TotalAlloc
			maxMemory := 64 * 1024 * 1024 // 64 МБ в байтах

			// Проверяем, что результат валидный
			if result < 0 || result >= mod {
				t.Errorf("solve(%d) вернул недопустимое значение %d", tt.n, result)
			}

			// Проверяем ограничение памяти: 64 МБ
			memoryMB := float64(allocated) / (1024 * 1024)
			if allocated > uint64(maxMemory) {
				t.Errorf("solve(%d) использовало %.2f МБ памяти, что превышает ограничение 64 МБ",
					tt.n, memoryMB)
			}

			// Для больших значений выводим предупреждение, если память близка к лимиту
			if memoryMB > 50 {
				t.Logf("⚠️  solve(%d) использовало %.2f МБ памяти (близко к лимиту 64 МБ), результат: %d",
					tt.n, memoryMB, result)
			} else {
				t.Logf("solve(%d) использовало %.2f МБ памяти, результат: %d",
					tt.n, memoryMB, result)
			}
		})
	}
}

// BenchmarkSolve проверяет производительность решения для различных размеров входных данных
func BenchmarkSolve(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int
	}{
		{"n=1", 1},
		{"n=10", 10},
		{"n=100", 100},
		{"n=1000", 1000},
		{"n=10000", 10000},
		{"n=100000", 100000},
		{"n=1000000", 1000000},
		{"n=2000000", 2000000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := solve(bm.n)
				_ = result
			}
		})
	}
}

// TestModPow проверяет корректность функции modPow
func TestModPow(t *testing.T) {
	tests := []struct {
		name     string
		base     int64
		exp      int64
		m        int64
		expected int64
	}{
		{
			name:     "2^10 mod 1000",
			base:     2,
			exp:      10,
			m:        1000,
			expected: 24, // 1024 mod 1000 = 24
		},
		{
			name:     "3^0 mod 7",
			base:     3,
			exp:      0,
			m:        7,
			expected: 1,
		},
		{
			name:     "5^3 mod 13",
			base:     5,
			exp:      3,
			m:        13,
			expected: 8, // 125 mod 13 = 8
		},
		{
			name:     "2^(mod-2) mod mod (обратный элемент)",
			base:     2,
			exp:      mod - 2,
			m:        mod,
			expected: 500000004, // 2^(-1) mod (10^9+7)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modPow(tt.base, tt.exp, tt.m)
			if result != tt.expected {
				t.Errorf("modPow(%d, %d, %d) = %d, ожидалось %d",
					tt.base, tt.exp, tt.m, result, tt.expected)
			}
		})
	}
}

// TestModInverse проверяет корректность функции modInverse
func TestModInverse(t *testing.T) {
	tests := []struct {
		name     string
		a        int64
		m        int64
		expected int64
	}{
		{
			name:     "2^(-1) mod (10^9+7)",
			a:        2,
			m:        mod,
			expected: 500000004,
		},
		{
			name:     "3^(-1) mod (10^9+7)",
			a:        3,
			m:        mod,
			expected: 333333336,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modInverse(tt.a, tt.m)
			if result != tt.expected {
				t.Errorf("modInverse(%d, %d) = %d, ожидалось %d",
					tt.a, tt.m, result, tt.expected)
			}
			// Проверяем, что a * result ≡ 1 (mod m)
			product := (tt.a * result) % tt.m
			if product != 1 {
				t.Errorf("Проверка: %d * %d mod %d = %d, ожидалось 1",
					tt.a, result, tt.m, product)
			}
		})
	}
}

// TestSolveConsistency проверяет, что результаты согласованы для последовательных n
func TestSolveConsistency(t *testing.T) {
	// E[n] должно возрастать с ростом n (математическое ожидание увеличивается)
	// Но в модульной арифметике это не обязательно верно для результата
	// Просто проверяем, что результаты стабильны
	for n := 1; n <= 20; n++ {
		result1 := solve(n)
		result2 := solve(n)
		if result1 != result2 {
			t.Errorf("solve(%d) возвращает нестабильные результаты: %d != %d",
				n, result1, result2)
		}
	}
}
