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
		k        int
		a        []int
		expected int
	}{
		{
			name:     "Пример 1",
			n:        1,
			k:        1,
			a:        []int{1},
			expected: 1,
		},
		{
			name:     "Пример 2",
			n:        5,
			k:        3,
			a:        []int{1, 5, 1},
			expected: 8,
		},
		{
			name:     "Пример 3",
			n:        3,
			k:        4,
			a:        []int{1, 2, 3, 1},
			expected: 1,
		},
		{
			name:     "Пример 4 (из дополнений)",
			n:        5,
			k:        2,
			a:        []int{1, 5},
			expected: 8,
		},
		{
			name:     "Пример 5 (из дополнений)",
			n:        3,
			k:        3,
			a:        []int{1, 2, 3},
			expected: 1,
		},
		{
			name:     "Пример 6 (из дополнений)",
			n:        11,
			k:        2,
			a:        []int{1, 7},
			expected: 270,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.k, tt.a)
			if result != tt.expected {
				t.Errorf("solve(%d, %d, %v) = %d, ожидалось %d",
					tt.n, tt.k, tt.a, result, tt.expected)
			}
		})
	}
}

func TestSolveConstraints(t *testing.T) {
	tests := []struct {
		name string
		n    int
		k    int
		a    []int
	}{
		{
			name: "Минимальные значения",
			n:    1,
			k:    1,
			a:    []int{1},
		},
		{
			name: "Средние значения",
			n:    100,
			k:    10,
			a:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "Большие значения n",
			n:    1000,
			k:    50,
			a:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solve(tt.n, tt.k, tt.a)

			// Проверяем, что результат неотрицательный и меньше mod
			if result < 0 || result >= mod {
				t.Errorf("solve(%d, %d, %v) вернул недопустимое значение %d (должно быть в [0, %d))",
					tt.n, tt.k, tt.a, result, mod)
			}

			// Проверяем ограничения из условия задачи
			if tt.n < 1 || tt.n > 1000000000 {
				t.Errorf("n=%d не соответствует ограничению 1 <= n <= 10^9", tt.n)
			}
			if tt.k < 1 || tt.k > min(tt.n, 100000) {
				t.Errorf("k=%d не соответствует ограничению 1 <= k <= min(n, 10^5)", tt.k)
			}
			for i, ai := range tt.a {
				if ai < 1 || ai > min(tt.n, 1000000) {
					t.Errorf("a[%d]=%d не соответствует ограничению 1 <= a_i <= min(n, 10^6)", i, ai)
				}
			}
		})
	}
}

func TestSolvePerformance(t *testing.T) {
	tests := []struct {
		name string
		n    int
		k    int
		a    []int
	}{
		{
			name: "Минимальные значения",
			n:    1,
			k:    1,
			a:    []int{1},
		},
		{
			name: "Средние значения",
			n:    1000,
			k:    100,
			a:    generateArray(100, 1000),
		},
		{
			name: "Большие значения n",
			n:    10000,
			k:    500,
			a:    generateArray(500, 1000000),
		},
		{
			name: "Очень большие значения n (n > 10^6)",
			n:    2000000,
			k:    1000,
			a:    generateArray(1000, 1000000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем время выполнения
			start := time.Now()
			result := solve(tt.n, tt.k, tt.a)
			elapsed := time.Since(start)

			// Проверяем ограничение времени: 1 с
			maxTime := 1 * time.Second
			if elapsed > maxTime {
				t.Errorf("solve(%d, %d, ...) выполнилось за %v, что превышает ограничение %v",
					tt.n, tt.k, elapsed, maxTime)
			}

			// Проверяем, что результат валидный
			if result < 0 || result >= mod {
				t.Errorf("solve(%d, %d, ...) вернул недопустимое значение %d",
					tt.n, tt.k, result)
			}

			t.Logf("solve(%d, %d, ...) выполнилось за %v, результат: %d", tt.n, tt.k, elapsed, result)
		})
	}
}

func TestSolveMemoryUsage(t *testing.T) {
	tests := []struct {
		name string
		n    int
		k    int
		a    []int
	}{
		{
			name: "Минимальные значения",
			n:    1,
			k:    1,
			a:    []int{1},
		},
		{
			name: "Средние значения",
			n:    1000,
			k:    100,
			a:    generateArray(100, 1000),
		},
		{
			name: "Большие значения",
			n:    10000,
			k:    500,
			a:    generateArray(500, 1000000),
		},
		{
			name: "Очень большие значения n (n > 10^6)",
			n:    2000000,
			k:    1000,
			a:    generateArray(1000, 1000000),
		},
		{
			name: "Очень большие значения n (n = 10^7)",
			n:    10000000,
			k:    10000,
			a:    generateArray(10000, 1000000),
		},
		{
			name: "Максимальные значения n (n = 10^8)",
			n:    100000000,
			k:    100000,
			a:    generateArray(100000, 1000000),
		},
		{
			name: "Максимальные значения n (n = 10^9) с минимальным k",
			n:    1000000000,
			k:    1,
			a:    []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем память до выполнения
			var m1, m2 runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&m1)

			result := solve(tt.n, tt.k, tt.a)

			// Измеряем память после выполнения
			runtime.GC()
			runtime.ReadMemStats(&m2)

			// Вычисляем использованную память в байтах
			allocated := m2.TotalAlloc - m1.TotalAlloc
			maxMemory := 128 * 1024 * 1024 // 128 МБ в байтах

			// Проверяем, что результат валидный
			if result < 0 || result >= mod {
				t.Errorf("solve(%d, %d, ...) вернул недопустимое значение %d",
					tt.n, tt.k, result)
			}

			// Проверяем ограничение памяти: 128 МБ
			memoryMB := float64(allocated) / (1024 * 1024)
			if allocated > uint64(maxMemory) {
				t.Errorf("solve(%d, %d, ...) использовало %.2f МБ памяти, что превышает ограничение 128 МБ",
					tt.n, tt.k, memoryMB)
			}

			// Для больших значений выводим предупреждение, если память близка к лимиту
			if memoryMB > 100 {
				t.Logf("⚠️  solve(%d, %d, ...) использовало %.2f МБ памяти (близко к лимиту 128 МБ), результат: %d",
					tt.n, tt.k, memoryMB, result)
			} else {
				t.Logf("solve(%d, %d, ...) использовало %.2f МБ памяти, результат: %d",
					tt.n, tt.k, memoryMB, result)
			}
		})
	}
}

// BenchmarkSolve проверяет производительность решения для различных размеров входных данных
func BenchmarkSolve(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int
		k    int
		a    []int
	}{
		{"Минимальные", 1, 1, []int{1}},
		{"Средние", 100, 10, generateArray(10, 100)},
		{"Большие", 1000, 100, generateArray(100, 1000)},
		{"Очень большие", 10000, 500, generateArray(500, 1000000)},
		{"n > 10^6", 2000000, 1000, generateArray(1000, 1000000)},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := solve(bm.n, bm.k, bm.a)
				_ = result
			}
		})
	}
}

// TestProcessLargePrimes проверяет корректность функции processLargePrimes
func TestProcessLargePrimes(t *testing.T) {
	tests := []struct {
		name     string
		low      int
		high     int
		expected []int
	}{
		{
			name:     "low > high (пустой диапазон)",
			low:      10,
			high:     5,
			expected: []int{},
		},
		{
			name:     "low == high (простое число)",
			low:      7,
			high:     7,
			expected: []int{7},
		},
		{
			name:     "low == high (составное число)",
			low:      8,
			high:     8,
			expected: []int{},
		},
		{
			name:     "Небольшой диапазон",
			low:      10,
			high:     20,
			expected: []int{11, 13, 17, 19},
		},
		{
			name:     "Диапазон с одним простым",
			low:      14,
			high:     16,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := []int{}
			processLargePrimes(tt.low, tt.high, func(p int) {
				result = append(result, p)
			})
			if len(result) != len(tt.expected) {
				t.Errorf("processLargePrimes(%d, %d) вернул %d простых чисел, ожидалось %d: %v vs %v",
					tt.low, tt.high, len(result), len(tt.expected), result, tt.expected)
				return
			}
			for i, prime := range result {
				if i >= len(tt.expected) || prime != tt.expected[i] {
					t.Errorf("processLargePrimes(%d, %d) = %v, ожидалось %v",
						tt.low, tt.high, result, tt.expected)
					return
				}
			}
		})
	}
}

// TestIntSqrt проверяет корректность функции intSqrt
func TestIntSqrt(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{
			name:     "n = 0",
			n:        0,
			expected: 0,
		},
		{
			name:     "n = 1",
			n:        1,
			expected: 1,
		},
		{
			name:     "n = 2",
			n:        2,
			expected: 1,
		},
		{
			name:     "n = 4 (точный квадрат)",
			n:        4,
			expected: 2,
		},
		{
			name:     "n = 9 (точный квадрат)",
			n:        9,
			expected: 3,
		},
		{
			name:     "n = 10",
			n:        10,
			expected: 3,
		},
		{
			name:     "n = 100 (точный квадрат)",
			n:        100,
			expected: 10,
		},
		{
			name:     "n = 1000000",
			n:        1000000,
			expected: 1000,
		},
		{
			name:     "n = 1000000000",
			n:        1000000000,
			expected: 31622,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intSqrt(tt.n)
			if result != tt.expected {
				t.Errorf("intSqrt(%d) = %d, ожидалось %d", tt.n, result, tt.expected)
			}
			// Проверяем, что результат корректен: result^2 <= n < (result+1)^2
			if result*result > tt.n {
				t.Errorf("intSqrt(%d) = %d, но %d^2 = %d > %d", tt.n, result, result, result*result, tt.n)
			}
			if tt.n > 0 && (result+1)*(result+1) <= tt.n {
				t.Errorf("intSqrt(%d) = %d, но (%d+1)^2 = %d <= %d", tt.n, result, result, (result+1)*(result+1), tt.n)
			}
		})
	}
}

// generateArray генерирует массив из k элементов со значениями от 1 до max
func generateArray(k, max int) []int {
	arr := make([]int, k)
	for i := 0; i < k; i++ {
		arr[i] = (i % max) + 1
	}
	return arr
}
