package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Тест на примере из условия
func TestSolveExample(t *testing.T) {
	input := `2
2 2
1 2
2 3
1 2
8 6
1 3 10 0 2 4 6 7
3 5
4 9
5 8
1 6
1 0
8 10`
	expected := `1
1
2
7
8
7
5
5
5
3`

	runTest(t, input, expected)
}

// Тест на граничные значения
func TestSolveEdgeCases(t *testing.T) {
	// 1. Один элемент, много запросов
	input1 := `1
1 3
0
1 1
1 2
1 0`
	// Пояснение:
	// Изначально [0]. Mex({0 ^ x}) = ?
	// При x=0: {0}. mex=1. max=1 (при x=0).
	// Запрос 1: a[0]=1. [1]. x=0->{1}, mex=0. x=1->{0}, mex=1. max=1.
	// Запрос 2: a[0]=2. [2]. x=0->{2}, mex=0. x=1->{3}, mex=0. x=2->{0}, mex=1. max=1.
	// Запрос 3: a[0]=0. [0]. max=1.
	expected1 := `1
1
1
1`
	runTest(t, input1, expected1)
}

// Тест на производительность (Time Limit)
func TestSolveTimeLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping time limit test in short mode")
	}

	// Генерируем большой тест
	// T=1, N=100000, Q=100000
	N := 100000
	Q := 100000
	
	// Создаем временный файл для ввода
	tmpfile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	writer := bufio.NewWriter(tmpfile)
	fmt.Fprintf(writer, "1\n%d %d\n", N, Q)
	
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < N; i++ {
		fmt.Fprintf(writer, "%d ", rng.Intn(1<<30))
	}
	fmt.Fprintln(writer)
	
	for i := 0; i < Q; i++ {
		idx := rng.Intn(N) + 1
		val := rng.Intn(1<<30)
		fmt.Fprintf(writer, "%d %d\n", idx, val)
	}
	writer.Flush()
	tmpfile.Close()

	// Перенаправляем stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	os.Stdin = file

	// Запускаем решение с измерением времени
	start := time.Now()
	
	// Перехватываем stdout, чтобы не засорять вывод теста
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	os.Stdout, _ = os.Open(os.DevNull)

	solve()
	
	elapsed := time.Since(start)
	t.Logf("Time elapsed: %v", elapsed)

	if elapsed > 4*time.Second {
		t.Errorf("Time limit exceeded: %v > 4s", elapsed)
	}
}

// Вспомогательная функция для запуска тестов
func runTest(t *testing.T, input, expected string) {
	// Создаем pipe для перехвата stdout
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	// Создаем pipe для подмены stdin
	rIn, wIn, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = rIn

	// Пишем входные данные
	go func() {
		wIn.Write([]byte(input))
		wIn.Close()
	}()

	// Запускаем решение
	solve()

	// Восстанавливаем stdout и stdin
	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	// Читаем вывод
	var outputBuf strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		outputBuf.WriteString(scanner.Text() + "\n")
	}

	output := strings.TrimSpace(outputBuf.String())
	expected = strings.TrimSpace(expected)

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

// Бенчмарк
func BenchmarkSolve(b *testing.B) {
	// Отключаем вывод
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	for i := 0; i < b.N; i++ {
		// Генерируем небольшой тест в памяти
		input := "1\n100 100\n"
		rng := rand.New(rand.NewSource(int64(i)))
		for j := 0; j < 100; j++ {
			input += fmt.Sprintf("%d ", rng.Intn(1000))
		}
		input += "\n"
		for j := 0; j < 100; j++ {
			input += fmt.Sprintf("%d %d\n", rng.Intn(100)+1, rng.Intn(1000))
		}

		// Подменяем stdin
		rIn, wIn, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = rIn
		
		go func() {
			wIn.Write([]byte(input))
			wIn.Close()
		}()

		solve()
		
		os.Stdin = oldStdin
	}
}

func init() {
	// Чтобы тесты запускались из любой директории
	os.Chdir(filepath.Dir(os.Args[0]))
}

