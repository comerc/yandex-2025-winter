package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestSolveTestCase(t *testing.T) {
	tests := []struct {
		n, m int
		a    []int64
		b    []int64
		want int64
	}{
		{
			n:    2,
			m:    2,
			a:    []int64{3, 13},
			b:    []int64{5, 7},
			want: 3,
		},
		{
			n:    1,
			m:    1,
			a:    []int64{3},
			b:    []int64{2},
			want: 1,
		},
		{
			n:    4,
			m:    2,
			a:    []int64{3, 11, 13, 15},
			b:    []int64{5, 7},
			want: 6,
		},
		{
			n:    4,
			m:    3,
			a:    []int64{3, 11, 13, 15},
			b:    []int64{5, 6, 7},
			want: 4,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			got := solveTestCase(tt.n, tt.m, tt.a, tt.b)
			if got != tt.want {
				t.Errorf("solveTestCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSolveTimeLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping time limit test in short mode")
	}

	// Максимальный тест: n=400, m=100
	n := 400
	m := 100
	a := make([]int64, n)
	b := make([]int64, m)

	// Заполняем случайными/большими числами
	for i := 0; i < n; i++ {
		a[i] = int64(1e9) - int64(i)
	}
	for i := 0; i < m; i++ {
		b[i] = int64(1e9) - int64(i)
	}

	start := time.Now()
	solveTestCase(n, m, a, b)
	elapsed := time.Since(start)

	t.Logf("Elapsed time for max test: %v", elapsed)

	if elapsed > 2*time.Second {
		t.Errorf("Time limit exceeded: %v > 2s", elapsed)
	}
}

func TestSolveMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory usage test in short mode")
	}

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Максимальный тест: n=400, m=100
	n := 400
	m := 100
	a := make([]int64, n)
	b := make([]int64, m)

	for i := 0; i < n; i++ {
		a[i] = int64(1e9)
	}
	for i := 0; i < m; i++ {
		b[i] = int64(1e9)
	}

	solveTestCase(n, m, a, b)

	runtime.GC()
	runtime.ReadMemStats(&m2)

	allocated := m2.TotalAlloc - m1.TotalAlloc
	memoryMB := float64(allocated) / (1024 * 1024)

	t.Logf("Memory used: %.2f MB", memoryMB)

	if memoryMB > 1024 { // 1 GB
		t.Errorf("Memory limit exceeded: %.2f MB > 1024 MB", memoryMB)
	}
}

func BenchmarkSolve(b *testing.B) {
	n := 400
	m := 100
	arrA := make([]int64, n)
	arrB := make([]int64, m)

	for i := 0; i < n; i++ {
		arrA[i] = int64(1e9)
	}
	for i := 0; i < m; i++ {
		arrB[i] = int64(1e9)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solveTestCase(n, m, arrA, arrB)
	}
}
