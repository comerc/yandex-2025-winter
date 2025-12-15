package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const mod int64 = 1000000007
const INV2 int64 = 500000004 // 2^(-1) mod 10^9+7
const INV3 int64 = 333333336 // 3^(-1) mod 10^9+7

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем a, q, L, R
	line, _ := reader.ReadString('\n')
	parts := strings.Fields(strings.TrimSpace(line))
	a, _ := strconv.ParseInt(parts[0], 10, 64)
	q, _ := strconv.ParseInt(parts[1], 10, 64)
	L, _ := strconv.ParseInt(parts[2], 10, 64)
	R, _ := strconv.ParseInt(parts[3], 10, 64)

	var result int64
	checkLimits(1*time.Second, 256, func() {
		result = solve(a, q, L, R)
	})

	writer.WriteString(fmt.Sprintf("%d\n", result))
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

// solve находит количество четвёрок (n, m, k, s) таких, что (aq^n - aq^m) / (aq^k - aq^s) - целое число
func solve(a, q, L, R int64) int64 {
	N := R - L + 1
	if N <= 0 {
		return 0
	}

	nMod := N % mod

	// Специальные случаи
	if a == 0 {
		// Если a = 0, то числитель всегда 0, знаменатель тоже 0 - деление на 0 не определено
		return 0
	}

	// Случай q = 0 (из 14b/main.go)
	if q == 0 {
		zeros := int64(0)
		if L <= 0 && R >= 0 {
			zeros = 1
		}
		if zeros == 0 {
			return 0
		}
		// Знаменатель != 0 <=> k=0 xor s=0
		nonZeros := (N - zeros) % mod
		validDenoms := (2 * zeros) % mod
		validDenoms = (validDenoms * nonZeros) % mod
		allNums := (nMod * nMod) % mod
		return (validDenoms * allNums) % mod
	}

	// Случай |q| = 1 (из 14b/main.go)
	if q == 1 {
		return 0
	}
	if q == -1 {
		var evens, odds int64
		// Безопасная проверка четности без float
		lIsEven := (L%2 == 0)
		if N%2 == 0 {
			evens, odds = N/2, N/2
		} else {
			if lIsEven {
				evens, odds = N/2+1, N/2
			} else {
				evens, odds = N/2, N/2+1
			}
		}
		validDenoms := (2 * (evens % mod)) % mod
		validDenoms = (validDenoms * (odds % mod)) % mod
		allNums := (nMod * nMod) % mod
		return (validDenoms * allNums) % mod
	}

	// Общий случай: q ≠ 0, q ≠ 1, a ≠ 0
	return solveGeneral(a, q, L, R)
}

// solveQOne обрабатывает случай q = 1
func solveQOne(a, L, R int64) int64 {
	count := int64(0)

	for n := L; n <= R; n++ {
		for m := L; m <= R; m++ {
			// Вычисляем числитель
			var numer int64
			if n >= 0 {
				numer = a
			} else {
				numer = 0
			}
			if m >= 0 {
				numer -= a
			} else {
				numer -= 0
			}

			for k := L; k <= R; k++ {
				for s := L; s <= R; s++ {
					if k == s {
						continue
					}

					// Вычисляем знаменатель
					var denom int64
					if k >= 0 {
						denom = a
					} else {
						denom = 0
					}
					if s >= 0 {
						denom -= a
					} else {
						denom -= 0
					}

					if denom == 0 {
						continue
					}

					if checkDivisibilityInt(numer, denom) {
						count++
					}
				}
			}
		}
	}

	result := count % mod
	if result < 0 {
		result += mod
	}
	return result
}

// solveQZeroBruteForce обрабатывает случай q = 0 через brute force
func solveQZeroBruteForce(a, L, R int64) int64 {
	count := int64(0)

	for n := L; n <= R; n++ {
		for m := L; m <= R; m++ {
			// Вычисляем числитель: aq^n - aq^m
			// При q=0: q^n = 1 если n=0, иначе 0
			var qn, qm int64
			if n == 0 {
				qn = 1
			} else {
				qn = 0
			}
			if m == 0 {
				qm = 1
			} else {
				qm = 0
			}
			numer := a*qn - a*qm

			for k := L; k <= R; k++ {
				for s := L; s <= R; s++ {
					if k == s {
						continue
					}

					// Вычисляем знаменатель: aq^k - aq^s
					var qk, qs int64
					if k == 0 {
						qk = 1
					} else {
						qk = 0
					}
					if s == 0 {
						qs = 1
					} else {
						qs = 0
					}
					denom := a*qk - a*qs

					if denom == 0 {
						continue
					}

					if checkDivisibilityInt(numer, denom) {
						count++
					}
				}
			}
		}
	}

	result := count % mod
	if result < 0 {
		result += mod
	}
	return result
}

// solveGeneral обрабатывает общий случай q ≠ 0, q ≠ 1, a ≠ 0
// Полностью копируем логику из 14b/main.go
func solveGeneral(a, q, L, R int64) int64 {
	N := R - L + 1
	if N <= 0 {
		return 0
	}

	nMod := N % mod

	// 1. Вклад n = m (числитель 0). Знаменатель любой ненулевой (k != s).
	// Кол-во = N * (N * (N - 1))
	ans := (nMod * nMod) % mod
	ans = (ans * ((nMod - 1 + mod) % mod)) % mod

	// 2. Вклад n != m. Используем Block Summation (Sqrt Decomposition)
	// Формула: 2 * [ (kB)^2 - kB(2N+1) + N(N+1) ]
	// где k - множитель (A = k*B)

	// Предподсчет постоянных частей формулы
	termN_N1 := (nMod * ((nMod + 1) % mod)) % mod // N(N+1)
	term2N_1 := ((2 * nMod) + 1) % mod            // 2N+1

	limit := N - 1
	var l int64 = 1

	// Основной цикл оптимизирован (без вызовов функций внутри)
	for l <= limit {
		kMax := limit / l
		r := limit / kMax
		if r > limit {
			r = limit
		}

		kMaxMod := kMax % mod

		// Вычисляем суммы для k от 1 до kMax
		// s1 = sum(k) = k(k+1)/2
		s1 := (kMaxMod * (kMaxMod + 1)) % mod
		s1 = (s1 * INV2) % mod

		// s2 = sum(k^2) = k(k+1)(2k+1)/6
		// s2 = s1 * (2k+1) / 3
		term2k1 := ((2 * kMaxMod) + 1) % mod
		s2 := (s1 * term2k1) % mod
		s2 = (s2 * INV3) % mod

		s0 := kMaxMod

		// Вычисляем суммы для B на отрезке [l, r]
		lMod := l % mod
		rMod := r % mod

		// Сумма 1..r для B и B^2
		sumR1 := (rMod * (rMod + 1)) % mod
		sumR1 = (sumR1 * INV2) % mod
		term2r1 := ((2 * rMod) + 1) % mod
		sumR2 := (sumR1 * term2r1) % mod
		sumR2 = (sumR2 * INV3) % mod

		// Сумма 1..l-1 для B и B^2
		lm1 := (lMod - 1 + mod) % mod
		sumL1 := (lm1 * (lm1 + 1)) % mod
		sumL1 = (sumL1 * INV2) % mod
		term2l1 := ((2 * lm1) + 1) % mod
		sumL2 := (sumL1 * term2l1) % mod
		sumL2 = (sumL2 * INV3) % mod

		// Разность сумм (значения на отрезке)
		ss2 := (sumR2 - sumL2 + mod) % mod
		ss1 := (sumR1 - sumL1 + mod) % mod
		ss0 := (rMod - lMod + 1 + mod) % mod

		// Собираем итоговое выражение для блока
		// blockSum = 2 * [ s2*ss2 - s1*ss1*(2N+1) + s0*ss0*N(N+1) ]

		p1 := (s2 * ss2) % mod

		p2 := (term2N_1 * s1) % mod
		p2 = (p2 * ss1) % mod

		p3 := (termN_N1 * s0) % mod
		p3 = (p3 * ss0) % mod

		blockSum := (p1 - p2 + mod) % mod
		blockSum = (blockSum + p3) % mod
		blockSum = (blockSum * 2) % mod

		ans = (ans + blockSum) % mod

		l = r + 1
	}

	// --- Коррекция для q = -2 ---
	// При B=2, частное целое для ВСЕХ нечетных A.
	// Стандартный алгоритм учел только четные A (кратные B=2).
	// Нужно добавить вклад всех нечетных A.
	if q == -2 && N >= 3 {
		maxOdd := limit
		if maxOdd%2 == 0 {
			maxOdd--
		}

		if maxOdd >= 1 {
			// Кол-во нечетных чисел <= maxOdd
			cnt := ((maxOdd + 1) / 2) % mod

			// sumA1 = сумма нечетных A = cnt^2
			sumA1 := (cnt * cnt) % mod

			// sumA2 = сумма квадратов нечетных A = cnt(4*cnt^2 - 1)/3
			termSq := (4 * cnt * cnt) % mod
			termSq = (termSq - 1 + mod) % mod
			sumA2 := (cnt * termSq) % mod
			sumA2 = (sumA2 * INV3) % mod

			// Подставляем в общую формулу 2 * [ A^2 - A(2N+1) + N(N+1) ]
			val := sumA2
			subVal := (term2N_1 * sumA1) % mod
			val = (val - subVal + mod) % mod
			addVal := (termN_N1 * cnt) % mod
			val = (val + addVal) % mod

			totalOdd := (val * 2) % mod

			// КОРРЕКЦИЯ ДЛЯ A=1 (пересечение условий)
			// При A=1, B=2 стандартная формула дает завышение на 4 единицы
			// (из-за того, что A < B, "хвосты" диапазонов ведут себя иначе).
			totalOdd = (totalOdd - 4 + mod) % mod

			ans = (ans + totalOdd) % mod
		}
	}

	return ans
}

// countValidKSForPair считает количество валидных (k, s) для конкретной пары (n, m)
// Использует оптимизацию: перебираем только делители d для |k-s|
// numer содержит (q^n - q^m), проверяем делимость на (q^k - q^s)
// Избегаем переполнения, не умножая на a, так как делимость эквивалентна при a ≠ 0
func countValidKSForPair(q, L, R int64, d int64, numer int64) int64 {
	N := R - L + 1
	if N < 2 {
		return 0
	}

	count := int64(0)

	// Находим все делители d
	divisors := []int64{}
	for e := int64(1); e <= d; e++ {
		if d%e == 0 {
			divisors = append(divisors, e)
		}
	}

	// Для каждого делителя e перебираем пары (k, s) с |k-s| = e
	for _, e := range divisors {
		if e >= N {
			continue
		}

		// Перебираем пары (k, s) с |k-s| = e
		for k := L; k <= R; k++ {
			s1 := k + e
			s2 := k - e

			// Проверяем пару (k, k+e)
			if s1 >= L && s1 <= R {
				qk := computePowerSafe(q, k)
				qs := computePowerSafe(q, s1)
				// Избегаем переполнения: проверяем делимость (q^n - q^m) на (q^k - q^s)
				denom := qk - qs

				if denom != 0 && checkDivisibilityInt(numer, denom) {
					count++
				}
			}

			// Проверяем пару (k, k-e), если она отличается
			if s2 >= L && s2 <= R && s2 != s1 {
				qk := computePowerSafe(q, k)
				qs := computePowerSafe(q, s2)
				// Избегаем переполнения: проверяем делимость (q^n - q^m) на (q^k - q^s)
				denom := qk - qs

				if denom != 0 && checkDivisibilityInt(numer, denom) {
					count++
				}
			}
		}
	}

	return count
}

// checkDivisibilityInt проверяет, делится ли numer на denom в целых числах
func checkDivisibilityInt(numer, denom int64) bool {
	if denom == 0 {
		return false
	}
	absNumer := numer
	if absNumer < 0 {
		absNumer = -absNumer
	}
	absDenom := denom
	if absDenom < 0 {
		absDenom = -absDenom
	}
	return absNumer%absDenom == 0
}

// checkDivisibilityGeometric проверяет делимость через прямую проверку без переполнения
// Проверяет делимость (q^n - q^m) на (q^k - q^s)
func checkDivisibilityGeometric(q, n, m, k, s int64) bool {
	// Вычисляем q^n, q^m, q^k, q^s
	qn := computePowerSafe(q, n)
	qm := computePowerSafe(q, m)
	qk := computePowerSafe(q, k)
	qs := computePowerSafe(q, s)

	// Вычисляем числитель и знаменатель
	numer := qn - qm
	denom := qk - qs

	if denom == 0 {
		return false
	}

	// Проверяем делимость напрямую
	return checkDivisibilityInt(numer, denom)
}

// checkDivisibilityIntDirect проверяет делимость напрямую с обработкой переполнения
func checkDivisibilityIntDirect(q, diffNM, diffKS, powerDiff int64) bool {
	// Вычисляем числитель: q^diffNM - 1
	qDiffNM := computePowerSafe(q, diffNM)
	numer := qDiffNM - 1

	// Вычисляем знаменатель: q^powerDiff * (q^diffKS - 1)
	qPowerDiff := computePowerSafe(q, powerDiff)
	qDiffKS := computePowerSafe(q, diffKS)
	denom := qPowerDiff * (qDiffKS - 1)

	if denom == 0 {
		return false
	}

	return checkDivisibilityInt(numer, denom)
}

// computePowerSafe вычисляет q^exp безопасно
func computePowerSafe(q, exp int64) int64 {
	if exp < 0 {
		return 0
	}
	if exp == 0 {
		return 1
	}
	if exp > 100 {
		return 0
	}
	result := int64(1)
	for i := int64(0); i < exp; i++ {
		prev := result
		result *= q
		if result < 0 || (q != 0 && result/q != prev) {
			return 0
		}
	}
	return result
}

// addMod вычисляет (a + b) mod mod
func addMod(a, b int64) int64 {
	result := (a + b) % mod
	if result < 0 {
		result += mod
	}
	return result
}

// mulMod вычисляет (a * b) mod mod
func mulMod(a, b int64) int64 {
	a %= mod
	b %= mod
	if a < 0 {
		a += mod
	}
	if b < 0 {
		b += mod
	}
	result := (a * b) % mod
	if result < 0 {
		result += mod
	}
	return result
}
