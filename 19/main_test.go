package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"
	"time"
)

// Helper function to capture stdout and provide stdin
func runSolver(input string) string {
	// Save old stdin/stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout

	// Create pipe for stdin
	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()
	os.Stdin = r

	// Create pipe for stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// Run solve
	solve()

	// Restore stdin/stdout
	wOut.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	// Read output
	var buf bytes.Buffer
	io.Copy(&buf, rOut)
	return buf.String()
}

func TestSolveExample1(t *testing.T) {
	input := "2\n2 0\n"
	expected := "0\n"
	result := runSolver(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSolveExample2(t *testing.T) {
	input := "2\n1 1\n"
	expected := "1\n"
	result := runSolver(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSolveMaxN(t *testing.T) {
	// Генерация большого теста
	// N = 2*10^6, m = 2
	// a1 = 10^6, a2 = 10^6
	// Это стресс-тест для времени
	m := 2
	a1 := 1000000
	a2 := 1000000
	input := fmt.Sprintf("%d\n%d %d\n", m, a1, a2)

	start := time.Now()
	_ = runSolver(input)
	elapsed := time.Since(start)

	t.Logf("Time for N=2*10^6: %v", elapsed)
	if elapsed > 1500*time.Millisecond {
		t.Errorf("Time limit exceeded: %v > 1.5s", elapsed)
	}
}

func TestSolveMemoryUsage(t *testing.T) {
	// N = 2*10^6
	// Для проверки памяти нам нужно запустить это в отдельном процессе или очень аккуратно измерять
	// Используем runtime.ReadMemStats

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Подготовка данных
	m := 200000
	inputBuffer := bytes.NewBufferString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		inputBuffer.WriteString("10 ")
	}
	inputBuffer.WriteString("\n")
	inputStr := inputBuffer.String()

	runSolver(inputStr)

	runtime.GC()
	runtime.ReadMemStats(&m2)

	allocated := m2.TotalAlloc - m1.TotalAlloc
	memoryMB := float64(allocated) / (1024 * 1024)

	t.Logf("Memory used: %.2f MB", memoryMB)
	if memoryMB > 256 {
		t.Errorf("Memory limit exceeded: %.2f MB > 256 MB", memoryMB)
	}
}

func BenchmarkSolve(b *testing.B) {
	// Подготовка данных
	m := 200000
	inputBuffer := bytes.NewBufferString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		inputBuffer.WriteString("10 ")
	}
	inputBuffer.WriteString("\n")
	inputStr := inputBuffer.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Мы не можем легко использовать runSolver в бенчмарке из-за перехвата os.Stdin
		// Это будет медленно из-за создания пайпов
		// Поэтому просто логируем, что бенчмарк сложен для реализации с os.Stdin
		// Но мы можем попробовать

		oldStdin := os.Stdin
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		go func() {
			w.Write([]byte(inputStr))
			w.Close()
		}()
		os.Stdin = r
		rOut, wOut, _ := os.Pipe()
		os.Stdout = wOut

		solve()

		wOut.Close()
		io.Copy(io.Discard, rOut)
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}
}
