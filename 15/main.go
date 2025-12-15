package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"
)

const MOD = 998244353

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	checkLimits(1*time.Second, 256, func() {
		solve(reader, writer)
	})
}

// checkLimits проверяет ограничения времени и памяти (работает только если установлена переменная окружения CHECK_LIMITS)
// Результаты выводятся в stderr, функция ничего не возвращает
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

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var T int
	fmt.Fscan(reader, &T)

	for tIdx := 0; tIdx < T; tIdx++ {
		var s, t string
		fmt.Fscan(reader, &s, &t)
		processTestCase(s, t, writer)
	}
}

func processTestCase(s, t string, writer *bufio.Writer) {
	n := len(s)
	m := len(t)

	// 1. Предподсчет количества разбиений для строки длины k.
	// Если длина k, то есть k-1 мест для разрыва, итого 2^(k-1) способов.
	// Для k=0 (пустая строка) считаем 1 способ.
	partitions := make([]int, n+1)
	partitions[0] = 1
	if n > 0 {
		p := 1
		for i := 1; i <= n; i++ {
			partitions[i] = p
			p = (p * 2) % MOD
		}
	}

	// 2. Предподсчет LCP (Longest Common Prefix) для всех пар суффиксов s и t.
	// lcp[i][j] = длина общего префикса s[i:] и t[j:]
	lcp := make([][]int, n+1)
	for i := range lcp {
		lcp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if s[i] == t[j] {
				lcp[i][j] = 1 + lcp[i+1][j+1]
			} else {
				lcp[i][j] = 0
			}
		}
	}

	// 3. Динамическое программирование
	// dp[i][u] - кол-во способов разбить суффикс s[i:], чтобы результат совпадал с t[u...]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	// База: пустой суффикс s совпадает с "началом" любой подстроки t (пустой строкой)
	for u := 0; u <= m; u++ {
		dp[n][u] = 1
	}

	ans := 0

	// Перебираем длину оставшегося суффикса s (от меньшего к большему с точки зрения "потребления" s справа налево)
	// i идет от n до 1. Мы пытаемся откусить кусок s[k...i-1].
	for i := n; i >= 1; i-- {
		// Длина уже сформированной части результата
		currentResLen := n - i

		for u := 0; u <= m; u++ {
			if dp[i][u] == 0 {
				continue
			}

			// Позиция в t, с которой мы должны сравнивать следующий кусок
			posInT := u + currentResLen

			// Если мы вышли за пределы t, сравнение невозможно (строка t кончилась раньше)
			if posInT >= m {
				continue
			}

			// Перебираем, где отрезать следующий кусок от s (индекс k)
			// Кусок будет s[k...i-1]
			for k := i - 1; k >= 0; k-- {
				chunkLen := i - k

				// Используем LCP для быстрого сравнения s[k...] и t[posInT...]
				val := lcp[k][posInT]

				if val >= chunkLen {
					// Кусок полностью совпал с частью t
					// Если мы не вышли за границы t, обновляем ДП
					if posInT+chunkLen <= m {
						dp[k][u] = (dp[k][u] + dp[i][u]) % MOD
					}
				} else {
					// Куски различаются.
					// Индекс различия относительно начала куска: val
					idxS := k + val
					idxT := posInT + val

					// Проверяем, что различие произошло в пределах строки t
					if idxT < m {
						// Если символ в s меньше символа в t, то результат кодирования лексикографически меньше
						if s[idxS] < t[idxT] {
							// Мы нашли "меньший" вариант.
							// Длина совпадающего префикса = currentResLen + val.
							matchLen := currentResLen + val

							// Этот результат будет меньше любой подстроки t, начинающейся в u,
							// которая длиннее matchLen.
							// Максимальная длина подстроки от u: m - u.
							// Подходящие длины: matchLen+1, matchLen+2, ..., m-u.
							count := (m - u) - matchLen

							if count > 0 {
								// Добавляем к ответу:
								// (способы дойти до i) * (способы разбить остаток s[0...k-1]) * (кол-во подстрок t)
								ways := (dp[i][u] * partitions[k]) % MOD
								term := (ways * count) % MOD
								ans = (ans + term) % MOD
							}
						}
						// Если s[idxS] > t[idxT], то результат больше, ничего не делаем.
					}
				}
			}
		}
	}

	// Обработка случаев, когда вся строка s (переставленная) является строгим префиксом подстроки t.
	// Это соответствует состояниям dp[0][u].
	// Результат имеет длину n и совпадает с t[u ... u+n-1].
	// Он будет меньше любой подстроки t[u...], длина которой > n.
	for u := 0; u <= m; u++ {
		if dp[0][u] > 0 {
			count := (m - u) - n
			if count > 0 {
				term := (dp[0][u] * count) % MOD
				ans = (ans + term) % MOD
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
