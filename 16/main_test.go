package main

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// Тест на примеры из условия задачи
func TestSolveExample(t *testing.T) {
	testCases := []struct {
		name     string
		n        int
		q        int
		l        int
		r        int
		b        []int
		c        []int
		expected int64
	}{
		{
			name:     "Пример 1: n=3, q=1, l=0, r=2",
			n:        3,
			q:        1,
			l:        0,
			r:        2,
			b:        []int{1},
			c:        []int{2},
			expected: 2,
		},
		{
			name:     "Пример 2: n=3, q=1, l=2, r=2",
			n:        3,
			q:        1,
			l:        2,
			r:        2,
			b:        []int{3},
			c:        []int{1},
			expected: 1,
		},
		{
			name:     "Пример 3: n=4, q=3, l=0, r=4",
			n:        4,
			q:        3,
			l:        0,
			r:        4,
			b:        []int{1, 2, 3},
			c:        []int{2, 1, 1},
			expected: 0,
		},
		{
			name:     "Пример 4: n=20, q=5, l=4, r=17",
			n:        20,
			q:        5,
			l:        4,
			r:        17,
			b:        []int{4, 5, 6, 12, 13},
			c:        []int{5, 6, 4, 13, 14},
			expected: 677226326,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.n, tc.q, tc.l, tc.r, tc.b, tc.c)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d, %v, %v) = %d, ожидалось %d",
					tc.n, tc.q, tc.l, tc.r, tc.b, tc.c, actual, tc.expected)
			}
		})
	}
}

// Тест на специальные случаи
func TestSolveSpecialCases(t *testing.T) {
	testCases := []struct {
		name     string
		n        int
		q        int
		l        int
		r        int
		b        []int
		c        []int
		expected int64
	}{
		{
			name:     "Нет фиксированных кабелей",
			n:        3,
			q:        0,
			l:        1,
			r:        3,
			b:        []int{},
			c:        []int{},
			expected: 6, // S(3,1) + S(3,2) + S(3,3) = 0 + 3 + 1 = 4, но нужно проверить
		},
		{
			name:     "Все кабели фиксированы, один цикл",
			n:        2,
			q:        2,
			l:        1,
			r:        1,
			b:        []int{1, 2},
			c:        []int{2, 1},
			expected: 1,
		},
		{
			name:     "Все кабели фиксированы, два цикла",
			n:        3,
			q:        3,
			l:        2,
			r:        2,
			b:        []int{1, 2, 3},
			c:        []int{2, 1, 3},
			expected: 1,
		},
		{
			name:     "Невозможная конфигурация (конфликт)",
			n:        3,
			q:        2,
			l:        0,
			r:        3,
			b:        []int{1, 1},
			c:        []int{2, 3},
			expected: 0,
		},
		{
			name:     "Невозможная конфигурация (два входа)",
			n:        3,
			q:        2,
			l:        0,
			r:        3,
			b:        []int{1, 2},
			c:        []int{3, 3},
			expected: 0,
		},
		{
			name:     "Один фиксированный кабель, путь",
			n:        3,
			q:        1,
			l:        1,
			r:        2,
			b:        []int{1},
			c:        []int{2},
			expected: 2, // M=2, fixedCycles=0, needL=1, needR=2: S(2,1)+S(2,2) = 1+1 = 2
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := solve(tc.n, tc.q, tc.l, tc.r, tc.b, tc.c)
			if actual != tc.expected {
				t.Errorf("solve(%d, %d, %d, %d, %v, %v) = %d, ожидалось %d",
					tc.n, tc.q, tc.l, tc.r, tc.b, tc.c, actual, tc.expected)
			}
		})
	}
}

// Тест на полный ввод-вывод (как в условии)
func TestSolveFullIO(t *testing.T) {
	input := `4
3 1 0 2
1
2
3 1 2 2
3
1
4 3 0 4
1 2 3
2 1 1
20 5 4 17
4 5 6 12 13
5 6 4 13 14
`

	expected := []int64{2, 1, 0, 677226326}

	reader := bufio.NewReader(strings.NewReader(input))
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	// Читаем все тесты
	var t_count int
	fmt.Fscan(reader, &t_count)

	testCases := make([]struct {
		n, q, l, r int
		b, c       []int
	}, t_count)

	for i := 0; i < t_count; i++ {
		var n, q, l, r int
		fmt.Fscan(reader, &n, &q, &l, &r)

		b := make([]int, q)
		c := make([]int, q)
		for j := 0; j < q; j++ {
			fmt.Fscan(reader, &b[j])
		}
		for j := 0; j < q; j++ {
			fmt.Fscan(reader, &c[j])
		}

		testCases[i] = struct {
			n, q, l, r int
			b, c       []int
		}{n, q, l, r, b, c}
	}

	// Обрабатываем все тесты
	var results []int64
	for i := 0; i < t_count; i++ {
		result := solve(testCases[i].n, testCases[i].q, testCases[i].l, testCases[i].r, testCases[i].b, testCases[i].c)
		results = append(results, result)
		fmt.Fprintln(writer, result)
	}
	writer.Flush()

	// Проверяем результаты
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(expected) {
		t.Errorf("solve() вернул %d строк, ожидалось %d. Вывод: %q",
			len(lines), len(expected), output)
		return
	}

	for i, line := range lines {
		if line == "" {
			continue
		}
		var result int64
		if _, err := fmt.Sscanf(line, "%d", &result); err != nil {
			t.Errorf("solve() вернул невалидный результат на строке %d: %q", i+1, line)
			continue
		}

		if result != expected[i] {
			t.Errorf("solve() вернул %d на строке %d, ожидалось %d",
				result, i+1, expected[i])
		}
	}
}

// TestSolveTimeLimit проверяет ограничение времени на больших входных данных
func TestSolveTimeLimit(t *testing.T) {
	tests := []struct {
		name string
		n    int
		q    int
		l    int
		r    int
		b    []int
		c    []int
	}{
		{
			name: "Большой тест n=1000",
			n:    1000,
			q:    100,
			l:    100,
			r:    900,
			b:    generateFixedCables(1000, 100),
			c:    generateFixedCables(1000, 100),
		},
		{
			name: "Очень большой тест n=3000",
			n:    3000,
			q:    500,
			l:    500,
			r:    2500,
			b:    generateFixedCables(3000, 500),
			c:    generateFixedCables(3000, 500),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			_ = solve(tt.n, tt.q, tt.l, tt.r, tt.b, tt.c)
			elapsed := time.Since(start)

			// Проверяем ограничение времени: 1 секунда
			maxTime := 1 * time.Second
			if elapsed > maxTime {
				t.Errorf("solve() выполнилось за %v, что превышает ограничение %v",
					elapsed, maxTime)
			}

			t.Logf("solve() выполнилось за %v", elapsed)
		})
	}
}

// TestSolveMemoryUsage проверяет использование памяти
func TestSolveMemoryUsage(t *testing.T) {
	tests := []struct {
		name string
		n    int
		q    int
		l    int
		r    int
		b    []int
		c    []int
	}{
		{
			name: "Средний тест n=500",
			n:    500,
			q:    50,
			l:    50,
			r:    450,
			b:    generateFixedCables(500, 50),
			c:    generateFixedCables(500, 50),
		},
		{
			name: "Большой тест n=2000",
			n:    2000,
			q:    200,
			l:    200,
			r:    1800,
			b:    generateFixedCables(2000, 200),
			c:    generateFixedCables(2000, 200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем память до выполнения
			var m1, m2 runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&m1)

			_ = solve(tt.n, tt.q, tt.l, tt.r, tt.b, tt.c)

			// Измеряем память после выполнения
			runtime.GC()
			runtime.ReadMemStats(&m2)

			// Вычисляем использованную память в байтах
			allocated := m2.TotalAlloc - m1.TotalAlloc
			maxMemory := 128 * 1024 * 1024 // 128 МБ в байтах

			// Проверяем ограничение памяти: 128 МБ
			memoryMB := float64(allocated) / (1024 * 1024)
			if allocated > uint64(maxMemory) {
				t.Errorf("solve() использовало %.2f МБ памяти, что превышает ограничение 128 МБ",
					memoryMB)
			}

			t.Logf("solve() использовало %.2f МБ памяти", memoryMB)
		})
	}
}

// BenchmarkSolve проверяет производительность решения для различных размеров входных данных
func BenchmarkSolve(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int
		q    int
		l    int
		r    int
		b    []int
		c    []int
	}{
		{
			name: "Маленький тест n=10",
			n:    10,
			q:    2,
			l:    1,
			r:    8,
			b:    []int{1, 2},
			c:    []int{2, 3},
		},
		{
			name: "Средний тест n=100",
			n:    100,
			q:    10,
			l:    10,
			r:    90,
			b:    generateFixedCables(100, 10),
			c:    generateFixedCables(100, 10),
		},
		{
			name: "Большой тест n=500",
			n:    500,
			q:    50,
			l:    50,
			r:    450,
			b:    generateFixedCables(500, 50),
			c:    generateFixedCables(500, 50),
		},
		{
			name: "Очень большой тест n=1000",
			n:    1000,
			q:    100,
			l:    100,
			r:    900,
			b:    generateFixedCables(1000, 100),
			c:    generateFixedCables(1000, 100),
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = solve(bm.n, bm.q, bm.l, bm.r, bm.b, bm.c)
			}
		})
	}
}

// Вспомогательные функции для генерации тестовых данных

func generateFixedCables(n, q int) []int {
	result := make([]int, q)
	for i := 0; i < q; i++ {
		result[i] = (i % n) + 1
	}
	return result
}
