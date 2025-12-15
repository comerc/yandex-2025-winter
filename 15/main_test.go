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
func TestSolve(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "Пример 1: abc и def",
			input:    "1\nabc def\n",
			expected: []int{24},
		},
		{
			name:     "Пример 2: ab и cd",
			input:    "1\nab cd\n",
			expected: []int{6},
		},
		{
			name:     "Несколько тестов",
			input:    "3\nabc def\nab cd\nx y\n",
			expected: []int{24, 6, 1},
		},
		{
			name:     "Пустые строки",
			input:    "1\n \n",
			expected: []int{0},
		},
		{
			name:     "Одна буква",
			input:    "1\na b\n",
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			// Измеряем время выполнения
			start := time.Now()
			solve(reader, writer)
			writer.Flush()
			elapsed := time.Since(start)

			// Проверяем ограничение времени: 1 секунда
			maxTime := 1 * time.Second
			if elapsed > maxTime {
				t.Errorf("solve() выполнилось за %v, что превышает ограничение %v",
					elapsed, maxTime)
			}

			// Парсим результат
			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			if len(lines) != len(tt.expected) {
				t.Errorf("solve() вернул %d строк, ожидалось %d. Вывод: %q",
					len(lines), len(tt.expected), output)
				return
			}

			// Проверяем каждую строку результата
			for i, line := range lines {
				if line == "" {
					continue
				}
				var result int
				if _, err := fmt.Sscanf(line, "%d", &result); err != nil {
					t.Errorf("solve() вернул невалидный результат на строке %d: %q", i+1, line)
					continue
				}

				if result != tt.expected[i] {
					t.Errorf("solve() вернул %d на строке %d, ожидалось %d",
						result, i+1, tt.expected[i])
				}
			}

			t.Logf("solve() выполнилось за %v", elapsed)
		})
	}
}

// TestSolveTimeLimit проверяет ограничение времени на больших входных данных
func TestSolveTimeLimit(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Большие строки",
			input: generateLargeTest(100, 100),
		},
		{
			name:  "Очень большие строки",
			input: generateLargeTest(500, 500),
		},
		{
			name:  "Много тестов",
			input: generateManyTests(100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			start := time.Now()
			solve(reader, writer)
			writer.Flush()
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
		name  string
		input string
	}{
		{
			name:  "Средние строки",
			input: generateLargeTest(200, 200),
		},
		{
			name:  "Большие строки",
			input: generateLargeTest(500, 500),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Измеряем память до выполнения
			var m1, m2 runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&m1)

			reader := bufio.NewReader(strings.NewReader(tt.input))
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			solve(reader, writer)
			writer.Flush()

			// Измеряем память после выполнения
			runtime.GC()
			runtime.ReadMemStats(&m2)

			// Вычисляем использованную память в байтах
			allocated := m2.TotalAlloc - m1.TotalAlloc
			maxMemory := 256 * 1024 * 1024 // 256 МБ в байтах

			// Проверяем ограничение памяти: 256 МБ
			memoryMB := float64(allocated) / (1024 * 1024)
			if allocated > uint64(maxMemory) {
				t.Errorf("solve() использовало %.2f МБ памяти, что превышает ограничение 256 МБ",
					memoryMB)
			}

			t.Logf("solve() использовало %.2f МБ памяти", memoryMB)
		})
	}
}

// BenchmarkSolve проверяет производительность решения для различных размеров входных данных
func BenchmarkSolve(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{
			name:  "Маленькие строки (10x10)",
			input: generateLargeTest(10, 10),
		},
		{
			name:  "Средние строки (50x50)",
			input: generateLargeTest(50, 50),
		},
		{
			name:  "Большие строки (100x100)",
			input: generateLargeTest(100, 100),
		},
		{
			name:  "Очень большие строки (200x200)",
			input: generateLargeTest(200, 200),
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				reader := bufio.NewReader(strings.NewReader(bm.input))
				var buf bytes.Buffer
				writer := bufio.NewWriter(&buf)
				solve(reader, writer)
				writer.Flush()
			}
		})
	}
}

// Вспомогательные функции для генерации тестовых данных

func generateLargeTest(n, m int) string {
	var builder strings.Builder
	builder.WriteString("1\n")

	// Генерируем строку s длины n
	for i := 0; i < n; i++ {
		builder.WriteByte(byte('a' + (i % 26)))
	}
	builder.WriteString(" ")

	// Генерируем строку t длины m
	for i := 0; i < m; i++ {
		builder.WriteByte(byte('a' + (i % 26)))
	}
	builder.WriteString("\n")

	return builder.String()
}

func generateManyTests(count int) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d\n", count))

	for i := 0; i < count; i++ {
		// Генерируем строки разной длины
		n := 10 + (i % 20)
		m := 10 + (i % 20)

		// Генерируем строку s
		for j := 0; j < n; j++ {
			builder.WriteByte(byte('a' + (j % 26)))
		}
		builder.WriteString(" ")

		// Генерируем строку t
		for j := 0; j < m; j++ {
			builder.WriteByte(byte('a' + (j % 26)))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
