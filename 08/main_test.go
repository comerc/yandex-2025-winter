package main

import (
	"runtime"
	"testing"
	"time"
)

func TestExamples(t *testing.T) {
	prefixSums := precompute()

	tests := []struct {
		name     string
		k, l, r  int
		expected int32
	}{
		{"1-интересные от 1 до 10", 1, 1, 10, 9},
		{"25 не является 2-интересным", 2, 25, 25, 0},
		{"30 является 3-интересным", 3, 30, 30, 1},
		{"3-интересные от 1 до 100", 3, 1, 100, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := query(prefixSums, tt.k, tt.l, tt.r)
			if result != tt.expected {
				t.Errorf("query(k=%d, l=%d, r=%d) = %d, expected %d",
					tt.k, tt.l, tt.r, result, tt.expected)
			}
		})
	}
}

func TestKInterestingNumbers(t *testing.T) {
	prefixSums := precompute()

	// Проверяем конкретные числа
	tests := []struct {
		name string
		n    int
		k    int
		want bool
	}{
		{"6=2*3 является 2-интересным", 6, 2, true},
		{"12=2*6 или 3*4 является 2-интересным", 12, 2, true},
		{"30=2*3*5 является 3-интересным", 30, 3, true},
		{"25=5*5 не является 2-интересным (множители не различны)", 25, 2, false},
		{"любое число >= 2 является 1-интересным", 7, 1, true},
		{"1 не является 1-интересным", 1, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Проверяем, является ли число k-интересным
			count := query(prefixSums, tt.k, tt.n, tt.n)
			got := count > 0
			if got != tt.want {
				t.Errorf("число %d k=%d: got %v, want %v", tt.n, tt.k, got, tt.want)
			}
		})
	}
}

func TestMaxK(t *testing.T) {
	prefixSums := precompute()

	// 2*3*4*5*6*7*8*9 = 362880 - это 8-интересное число
	count8 := query(prefixSums, 8, 362880, 362880)
	if count8 != 1 {
		t.Errorf("362880 должно быть 8-интересным, got count=%d", count8)
	}

	// Для k > 9 не должно быть k-интересных чисел в диапазоне до 700000
	count10 := query(prefixSums, 10, 1, 700000)
	if count10 != 0 {
		t.Errorf("Не должно быть 10-интересных чисел до 700000, got %d", count10)
	}
}

func TestEdgeCases(t *testing.T) {
	prefixSums := precompute()

	tests := []struct {
		name     string
		k, l, r  int
		expected int32
	}{
		{"k > maxK возвращает 0", 100, 1, 700000, 0},
		{"l > maxN возвращает 0", 1, 800000, 900000, 0},
		{"r > maxN ограничивается", 1, 699999, 800000, 2},
		{"l < 1 корректируется до 1 (l=0)", 1, 0, 10, 9},
		{"l < 1 корректируется до 1 (l=-1)", 1, -1, 10, 9},
		{"l < 1 корректируется до 1 (l=-5)", 1, -5, 10, 9},
		{"l = 0, r = 10 эквивалентно l = 1, r = 10", 1, 0, 10, 9},
		{"l < 1 для k=2", 2, 0, 20, query(prefixSums, 2, 1, 20)},
		{"l < 1 для k=3", 3, -10, 100, query(prefixSums, 3, 1, 100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := query(prefixSums, tt.k, tt.l, tt.r)
			if result != tt.expected {
				t.Errorf("query(k=%d, l=%d, r=%d) = %d, expected %d",
					tt.k, tt.l, tt.r, result, tt.expected)
			}
		})
	}
}

// TestLNegativeExplicit проверяет явно, что проверка l < 1 выполняется
func TestLNegativeExplicit(t *testing.T) {
	prefixSums := precompute()

	// Тест, который гарантированно проходит через проверку l < 1
	// k валидный (1 <= k <= maxK), l < 1, r валидный
	t.Run("l отрицательное, k=1", func(t *testing.T) {
		// При l=-1 проверка l > maxN будет false, так что дойдём до if l < 1
		result := query(prefixSums, 1, -1, 10)
		expected := query(prefixSums, 1, 1, 10) // должно быть эквивалентно l=1
		if result != expected {
			t.Errorf("query(k=1, l=-1, r=10) = %d, expected %d (same as l=1)", result, expected)
		}
	})

	t.Run("l=0, k=2", func(t *testing.T) {
		result := query(prefixSums, 2, 0, 20)
		expected := query(prefixSums, 2, 1, 20)
		if result != expected {
			t.Errorf("query(k=2, l=0, r=20) = %d, expected %d (same as l=1)", result, expected)
		}
	})

	t.Run("l отрицательное большое, k=3", func(t *testing.T) {
		result := query(prefixSums, 3, -100, 100)
		expected := query(prefixSums, 3, 1, 100)
		if result != expected {
			t.Errorf("query(k=3, l=-100, r=100) = %d, expected %d (same as l=1)", result, expected)
		}
	})
}

func TestPerformance(t *testing.T) {
	// Измеряем время предвычисления
	start := time.Now()
	prefixSums := precompute()
	precomputeTime := time.Since(start)

	t.Logf("Время предвычисления: %v", precomputeTime)

	if precomputeTime > 3*time.Second {
		t.Errorf("Предвычисление слишком долгое: %v > 3s", precomputeTime)
	}

	// Измеряем время запросов
	start = time.Now()
	for i := 0; i < 50000; i++ {
		query(prefixSums, (i%9)+1, 1, 700000)
	}
	queryTime := time.Since(start)

	t.Logf("Время 50000 запросов: %v", queryTime)

	if queryTime > 100*time.Millisecond {
		t.Errorf("Запросы слишком долгие: %v > 100ms", queryTime)
	}
}

func TestMemoryUsage(t *testing.T) {
	var m runtime.MemStats

	prefixSums := precompute()
	_ = prefixSums // используем переменную

	runtime.ReadMemStats(&m)

	allocatedMB := float64(m.Alloc) / (1024 * 1024)
	totalAllocMB := float64(m.TotalAlloc) / (1024 * 1024)
	t.Logf("Текущее выделение: %.2f MB", allocatedMB)
	t.Logf("Всего выделено: %.2f MB", totalAllocMB)

	// Лимит 256 МБ
	if allocatedMB > 256 {
		t.Errorf("Превышен лимит памяти: %.2f MB > 256 MB", allocatedMB)
	}
}

func BenchmarkPrecompute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		precompute()
	}
}

func BenchmarkQuery(b *testing.B) {
	prefixSums := precompute()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query(prefixSums, (i%9)+1, 1, 700000)
	}
}
