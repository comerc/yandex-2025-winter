package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

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

func solve(n int) (string, string) {
	// Find M such that 10^M > n-1
	m := 0
	limit := 1
	// Using int for limit is safe because n <= 10^4.
	// 10^4 fits in int.
	for limit <= n-1 {
		limit *= 10
		m++
	}
	if m == 0 {
		m = 1
	}

	s := strings.Repeat("9", m)
	return s, s
}

func main() {
	var n int
	// Handle possible EOF or error
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	checkLimits(1*time.Second, 256, func() {
		a, d := solve(n)
		fmt.Println(a)
		fmt.Println(d)
	})
}
