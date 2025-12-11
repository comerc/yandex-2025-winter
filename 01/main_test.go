package main

import (
	"runtime"
	"testing"
	"time"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name      string
		R         int
		B         int
		expectedW int
		expectedH int
	}{
		{
			name:      "Пример 1",
			R:         8,
			B:         1,
			expectedW: 3,
			expectedH: 3,
		},
		{
			name:      "Пример 2",
			R:         10,
			B:         2,
			expectedW: 4,
			expectedH: 3,
		},
		{
			name:      "Пример 3",
			R:         24,
			B:         24,
			expectedW: 8,
			expectedH: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			W, H := solve(tt.R, tt.B)
			if W != tt.expectedW || H != tt.expectedH {
				t.Errorf("solve(%d, %d) = (%d, %d), ожидалось (%d, %d)",
					tt.R, tt.B, W, H, tt.expectedW, tt.expectedH)
			}
			// Проверяем условие W >= H
			if W < H {
				t.Errorf("solve(%d, %d) вернул W=%d < H=%d, но должно быть W >= H",
					tt.R, tt.B, W, H)
			}
			// Проверяем правильность формулы для красных плиток
			calculatedR := 2*W + 2*H - 4
			if calculatedR != tt.R {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но R должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.R, calculatedR)
			}
			// Проверяем правильность формулы для синих плиток
			calculatedB := (W - 2) * (H - 2)
			if calculatedB != tt.B {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но B должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.B, calculatedB)
			}
		})
	}
}

func TestSolveConstraints(t *testing.T) {
	tests := []struct {
		name string
		R    int
		B    int
		W    int
		H    int
	}{
		{
			name: "Минимальные значения",
			R:    8,
			B:    1,
			W:    3,
			H:    3,
		},
		{
			name: "Средние значения (W=10, H=8)",
			R:    32,
			B:    48,
			W:    10,
			H:    8,
		},
		{
			name: "Средние значения (W=5, H=4)",
			R:    14,
			B:    6,
			W:    5,
			H:    4,
		},
		{
			name: "Большие значения (W=20, H=15)",
			R:    66,
			B:    234,
			W:    20,
			H:    15,
		},
		{
			name: "Большие значения (W=100, H=50)",
			R:    296,
			B:    4704,
			W:    100,
			H:    50,
		},
		{
			name: "R большое, B минимальное (W=3, H=3)",
			R:    8,
			B:    1,
			W:    3,
			H:    3,
		},
		{
			name: "Прямоугольник (W=50, H=30)",
			R:    156,
			B:    1344,
			W:    50,
			H:    30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			W, H := solve(tt.R, tt.B)

			// Проверяем, что решение найдено
			if W == 0 || H == 0 {
				t.Errorf("solve(%d, %d) не нашло решение", tt.R, tt.B)
				return
			}

			// Проверяем, что решение совпадает с ожидаемым
			if W != tt.W || H != tt.H {
				t.Errorf("solve(%d, %d) = (%d, %d), ожидалось (%d, %d)",
					tt.R, tt.B, W, H, tt.W, tt.H)
			}

			// Проверяем условие W >= H
			if W < H {
				t.Errorf("solve(%d, %d) вернул W=%d < H=%d, но должно быть W >= H",
					tt.R, tt.B, W, H)
			}

			// Проверяем, что W и H >= 2 (иначе не будет внутренних синих плиток)
			if W < 2 || H < 2 {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но оба должны быть >= 2",
					tt.R, tt.B, W, H)
			}

			// Проверяем правильность формулы для красных плиток
			calculatedR := 2*W + 2*H - 4
			if calculatedR != tt.R {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но R должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.R, calculatedR)
			}

			// Проверяем правильность формулы для синих плиток
			calculatedB := (W - 2) * (H - 2)
			if calculatedB != tt.B {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но B должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.B, calculatedB)
			}
		})
	}
}

func TestSolveEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		R    int
		B    int
		W    int
		H    int
	}{
		{
			name: "Квадрат минимального размера",
			R:    8,
			B:    1,
			W:    3,
			H:    3,
		},
		{
			name: "Прямоугольник 4x3",
			R:    10,
			B:    2,
			W:    4,
			H:    3,
		},
		{
			name: "Большой квадрат (W=50, H=50)",
			R:    196,
			B:    2304, // (50-2)*(50-2) = 48*48 = 2304
			W:    50,
			H:    50,
		},
		{
			name: "Узкий прямоугольник (W=100, H=2)",
			R:    200,
			B:    0, // (100-2)*(2-2) = 98*0 = 0
			W:    100,
			H:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Пропускаем тест для B=0, так как по условию B>=1
			if tt.B == 0 {
				return
			}

			W, H := solve(tt.R, tt.B)

			if W == 0 || H == 0 {
				t.Errorf("solve(%d, %d) не нашло решение", tt.R, tt.B)
				return
			}

			// Проверяем, что решение совпадает с ожидаемым
			if W != tt.W || H != tt.H {
				t.Errorf("solve(%d, %d) = (%d, %d), ожидалось (%d, %d)",
					tt.R, tt.B, W, H, tt.W, tt.H)
			}

			// Проверяем условие W >= H
			if W < H {
				t.Errorf("solve(%d, %d) вернул W=%d < H=%d, но должно быть W >= H",
					tt.R, tt.B, W, H)
			}

			// Проверяем правильность формул
			calculatedR := 2*W + 2*H - 4
			if calculatedR != tt.R {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но R должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.R, calculatedR)
			}

			calculatedB := (W - 2) * (H - 2)
			if calculatedB != tt.B {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но B должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.B, calculatedB)
			}
		})
	}
}

// TestSolveBoundaryValues проверяет граничные значения из ограничений задачи
// R от 8 до 10^9, B от 1 до 10^9
func TestSolveBoundaryValues(t *testing.T) {
	tests := []struct {
		name string
		R    int
		B    int
		W    int
		H    int
	}{
		{
			name: "Минимальное R (8), минимальное B (1)",
			R:    8,
			B:    1,
			W:    3,
			H:    3,
		},
		{
			name: "Большие значения (W=1000, H=500)",
			R:    2996,
			B:    497004, // (1000-2)*(500-2) = 998*498 = 497004
			W:    1000,
			H:    500,
		},
		{
			name: "Очень большие значения (W=10000, H=5000)",
			R:    29996,
			B:    49970004, // (10000-2)*(5000-2) = 9998*4998 = 49970004
			W:    10000,
			H:    5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			W, H := solve(tt.R, tt.B)

			// Проверяем, что решение найдено
			if W == 0 || H == 0 {
				t.Errorf("solve(%d, %d) не нашло решение", tt.R, tt.B)
				return
			}

			// Проверяем, что решение совпадает с ожидаемым
			if W != tt.W || H != tt.H {
				t.Errorf("solve(%d, %d) = (%d, %d), ожидалось (%d, %d)",
					tt.R, tt.B, W, H, tt.W, tt.H)
			}

			// Проверяем условие W >= H
			if W < H {
				t.Errorf("solve(%d, %d) вернул W=%d < H=%d, но должно быть W >= H",
					tt.R, tt.B, W, H)
			}

			// Проверяем, что W и H >= 2
			if W < 2 || H < 2 {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но оба должны быть >= 2",
					tt.R, tt.B, W, H)
			}

			// Проверяем правильность формулы для красных плиток
			calculatedR := 2*W + 2*H - 4
			if calculatedR != tt.R {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но R должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.R, calculatedR)
			}

			// Проверяем правильность формулы для синих плиток
			calculatedB := (W - 2) * (H - 2)
			if calculatedB != tt.B {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но B должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.B, calculatedB)
			}

			// Проверяем ограничения из условия задачи
			if tt.R < 8 || tt.R > 1000000000 {
				t.Errorf("R=%d не соответствует ограничению 8 <= R <= 10^9", tt.R)
			}
			if tt.B < 1 || tt.B > 1000000000 {
				t.Errorf("B=%d не соответствует ограничению 1 <= B <= 10^9", tt.B)
			}
		})
	}
}

// TestSolvePerformance проверяет ограничения времени выполнения (300 мс) и памяти (512 МБ)
func TestSolvePerformance(t *testing.T) {
	tests := []struct {
		name string
		R    int
		B    int
		W    int
		H    int
	}{
		{
			name: "Минимальные значения",
			R:    8,
			B:    1,
			W:    3,
			H:    3,
		},
		{
			name: "Средние значения (W=100, H=50)",
			R:    296,
			B:    4704, // (100-2)*(50-2) = 98*48 = 4704
			W:    100,
			H:    50,
		},
		{
			name: "Большие значения (W=1000, H=500)",
			R:    2996,
			B:    497004, // (1000-2)*(500-2) = 998*498 = 497004
			W:    1000,
			H:    500,
		},
		{
			name: "Очень большие значения (W=10000, H=5000)",
			R:    29996,
			B:    49970004, // (10000-2)*(5000-2) = 9998*4998 = 49970004
			W:    10000,
			H:    5000,
		},
		{
			name: "Большой квадрат (W=31623, H=31623)",
			R:    126488,
			B:    999887641, // (31623-2)*(31623-2) = 31621*31621 = 999887641
			W:    31623,
			H:    31623,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем время выполнения
			start := time.Now()
			W, H := solve(tt.R, tt.B)
			elapsed := time.Since(start)

			// Проверяем ограничение времени: 300 мс
			maxTime := 300 * time.Millisecond
			if elapsed > maxTime {
				t.Errorf("solve(%d, %d) выполнилось за %v, что превышает ограничение %v",
					tt.R, tt.B, elapsed, maxTime)
			}

			// Проверяем, что решение найдено
			if W == 0 || H == 0 {
				t.Errorf("solve(%d, %d) не нашло решение", tt.R, tt.B)
				return
			}

			// Проверяем, что решение совпадает с ожидаемым
			if W != tt.W || H != tt.H {
				t.Errorf("solve(%d, %d) = (%d, %d), ожидалось (%d, %d)",
					tt.R, tt.B, W, H, tt.W, tt.H)
			}

			// Проверяем правильность решения
			calculatedR := 2*W + 2*H - 4
			if calculatedR != tt.R {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но R должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.R, calculatedR)
			}

			calculatedB := (W - 2) * (H - 2)
			if calculatedB != tt.B {
				t.Errorf("solve(%d, %d) вернул W=%d, H=%d, но B должно быть %d, а получилось %d",
					tt.R, tt.B, W, H, tt.B, calculatedB)
			}

			t.Logf("solve(%d, %d) выполнилось за %v", tt.R, tt.B, elapsed)
		})
	}
}

// TestSolveMemoryUsage проверяет использование памяти
func TestSolveMemoryUsage(t *testing.T) {
	tests := []struct {
		name string
		R    int
		B    int
		W    int
		H    int
	}{
		{
			name: "Минимальные значения",
			R:    8,
			B:    1,
			W:    3,
			H:    3,
		},
		{
			name: "Большие значения (W=1000, H=500)",
			R:    2996,
			B:    497004,
			W:    1000,
			H:    500,
		},
		{
			name: "Очень большие значения (W=10000, H=5000)",
			R:    29996,
			B:    49970004,
			W:    10000,
			H:    5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем память до выполнения
			var m1, m2 runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&m1)

			W, H := solve(tt.R, tt.B)

			// Измеряем память после выполнения
			runtime.GC()
			runtime.ReadMemStats(&m2)

			// Вычисляем использованную память в байтах
			allocated := m2.TotalAlloc - m1.TotalAlloc
			maxMemory := 512 * 1024 * 1024 // 512 МБ в байтах

			// Проверяем, что решение найдено
			if W == 0 || H == 0 {
				t.Errorf("solve(%d, %d) не нашло решение", tt.R, tt.B)
				return
			}

			// Проверяем правильность решения
			if W != tt.W || H != tt.H {
				t.Errorf("solve(%d, %d) = (%d, %d), ожидалось (%d, %d)",
					tt.R, tt.B, W, H, tt.W, tt.H)
			}

			// Проверяем ограничение памяти: 512 МБ
			// Примечание: это приблизительная проверка, так как Go управляет памятью автоматически
			if allocated > uint64(maxMemory) {
				t.Errorf("solve(%d, %d) использовало %d байт памяти, что превышает ограничение %d байт (512 МБ)",
					tt.R, tt.B, allocated, maxMemory)
			}

			t.Logf("solve(%d, %d) использовало примерно %d байт памяти (%.2f МБ)",
				tt.R, tt.B, allocated, float64(allocated)/(1024*1024))
		})
	}
}

// BenchmarkSolve проверяет производительность решения для различных размеров входных данных
func BenchmarkSolve(b *testing.B) {
	benchmarks := []struct {
		name string
		R    int
		B    int
	}{
		{"Минимальные", 8, 1},
		{"Средние", 296, 4704},
		{"Большие", 2996, 497004},
		{"Очень большие", 29996, 49970004},
		{"Максимальные", 126488, 999887641},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				W, H := solve(bm.R, bm.B)
				_ = W
				_ = H
			}
		})
	}
}
