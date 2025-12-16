package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"
)

const MOD = 998244353
const MAX_N = 3005

// Only store the prefix sums table to save memory (~72 MB)
var stirlingSum [MAX_N][MAX_N]int64

func init() {
	var prevStirling [MAX_N]int64
	var currStirling [MAX_N]int64

	prevStirling[0] = 1
	stirlingSum[0][0] = 1
	for j := 1; j < MAX_N; j++ {
		stirlingSum[0][j] = 1
	}

	for i := 1; i < MAX_N; i++ {
		currStirling[0] = 0 // S(n, 0) = 0 for n >= 1
		stirlingSum[i][0] = 0
		var currentSum int64 = 0

		for j := 1; j <= i; j++ {
			// S(n, k) = S(n-1, k-1) + (n-1)*S(n-1, k)
			val := prevStirling[j-1] + (int64(i-1)*prevStirling[j])%MOD
			currStirling[j] = val % MOD
			currentSum = (currentSum + currStirling[j]) % MOD
			stirlingSum[i][j] = currentSum
		}
		// Fill remaining sums with the total sum for this row
		for j := i + 1; j < MAX_N; j++ {
			stirlingSum[i][j] = currentSum
		}

		// Update prevStirling for next iteration
		for j := 0; j <= i; j++ {
			prevStirling[j] = currStirling[j]
		}
	}
}

// checkLimits проверяет ограничения времени и памяти (работает только если установлена переменная окружения CHECK_LIMITS)
func checkLimits(maxTime time.Duration, maxMemoryMB int, fn func()) {
	if os.Getenv("CHECK_LIMITS") == "" {
		fn()
		return
	}

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start := time.Now()
	fn()
	elapsed := time.Since(start)

	runtime.GC()
	runtime.ReadMemStats(&m2)

	allocated := m2.TotalAlloc - m1.TotalAlloc
	memoryMB := float64(allocated) / (1024 * 1024)

	maxMemoryBytes := uint64(maxMemoryMB) * 1024 * 1024

	timeOk := elapsed <= maxTime
	memoryOk := allocated <= maxMemoryBytes

	if !timeOk || !memoryOk {
		if elapsed > maxTime {
			fmt.Fprintf(os.Stderr, "⚠️ Превышено время: %v (лимит: %v)\n", elapsed, maxTime)
		}
		if memoryMB > float64(maxMemoryMB) {
			fmt.Fprintf(os.Stderr, "⚠️ Превышена память: %.2f МБ (лимит: %d МБ)\n", memoryMB, maxMemoryMB)
		}
	} else {
		fmt.Fprintf(os.Stderr, "✓ Время: %v, Память: %.2f МБ\n", elapsed, memoryMB)
	}
}

// solve находит количество способов завершить схему канатной дороги
func solve(n, q, l, r int, b, c []int) int64 {
	inDeg := make([]int, n+1)
	outDeg := make([]int, n+1)
	adj := make([]int, n+1)

	possible := true
	for i := 0; i < q; i++ {
		u, v := b[i], c[i]
		if outDeg[u] > 0 || inDeg[v] > 0 {
			possible = false
		}
		outDeg[u]++
		inDeg[v]++
		adj[u] = v
	}

	if !possible {
		return 0
	}

	// M = n - q (number of path components)
	// But we need to verify connectivity and count fixed cycles
	M := n - q

	// Count fixed cycles
	fixedCycles := 0
	visited := make([]bool, n+1)

	// First, traverse everything starting from nodes with inDeg == 0 (Path starts)
	for i := 1; i <= n; i++ {
		if inDeg[i] == 0 {
			curr := i
			for curr != 0 && !visited[curr] {
				visited[curr] = true
				if outDeg[curr] > 0 {
					curr = adj[curr]
				} else {
					curr = 0
				}
			}
		}
	}

	// Remaining unvisited nodes must be part of cycles
	for i := 1; i <= n; i++ {
		if !visited[i] {
			// Found a cycle
			fixedCycles++
			curr := i
			for !visited[curr] {
				visited[curr] = true
				curr = adj[curr]
			}
		}
	}

	// Range of cycles needed from path components
	needL := l - fixedCycles
	needR := r - fixedCycles

	if needL < 0 {
		needL = 0
	}
	if needR < 0 {
		return 0
	}
	if needL > M {
		return 0
	}
	if needR > M {
		needR = M
	}

	// Answer is sum of Stirling numbers [M][k] for k in [needL, needR]
	var sub int64 = 0
	if needL > 0 {
		sub = stirlingSum[M][needL-1]
	}
	ans := (stirlingSum[M][needR] - sub + MOD) % MOD
	return ans
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	// Читаем все тесты
	testCases := make([]struct {
		n, q, l, r int
		b, c       []int
	}, t)

	for i := 0; i < t; i++ {
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

	// Обрабатываем все тесты внутри checkLimits
	var results []int64
	checkLimits(1*time.Second, 128, func() {
		results = make([]int64, t)
		for i := 0; i < t; i++ {
			results[i] = solve(testCases[i].n, testCases[i].q, testCases[i].l, testCases[i].r, testCases[i].b, testCases[i].c)
		}
	})

	// Выводим результаты
	for i := 0; i < t; i++ {
		fmt.Fprintln(writer, results[i])
	}
}
