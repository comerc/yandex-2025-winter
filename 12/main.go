package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const mod = 998244353

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем n и m
	line, _ := reader.ReadString('\n')
	parts := strings.Fields(strings.TrimSpace(line))
	n, _ := strconv.ParseInt(parts[0], 10, 64)
	m, _ := strconv.Atoi(parts[1])

	// Читаем хорошие числа
	line, _ = reader.ReadString('\n')
	parts = strings.Fields(strings.TrimSpace(line))
	good := make(map[int]bool, m)
	for i := 0; i < m; i++ {
		val, _ := strconv.Atoi(parts[i])
		good[val] = true
	}

	result := solve(int64(n), good)
	writer.WriteString(fmt.Sprintf("%d\n", result))
}

// solve находит количество чудесных чисел длины n
// Чудесное число: без лидирующих нулей, сумма любых трех последовательных цифр - хорошее число
func solve(n int64, good map[int]bool) int {
	if n < 3 {
		return 0
	}

	// Состояние: последние две цифры (d1, d2) -> индекс = d1*10 + d2
	// Матрица переходов: M[i][j] = 1, если можно перейти от состояния i к состоянию j
	// i = d1*10 + d2, j = d2*10 + d3, переход возможен если d1+d2+d3 - хорошее число

	// Строим матрицу переходов 100x100
	M := make([][]int, 100)
	for i := 0; i < 100; i++ {
		M[i] = make([]int, 100)
	}

	for d1 := 0; d1 < 10; d1++ {
		for d2 := 0; d2 < 10; d2++ {
			from := d1*10 + d2
			for d3 := 0; d3 < 10; d3++ {
				sum := d1 + d2 + d3
				if good[sum] {
					to := d2*10 + d3
					M[from][to] = 1
				}
			}
		}
	}

	// Начальный вектор: для чисел длины 2 (d1, d2), где d1 != 0
	start := make([]int, 100)
	for d1 := 1; d1 < 10; d1++ {
		for d2 := 0; d2 < 10; d2++ {
			idx := d1*10 + d2
			start[idx] = 1
		}
	}

	// Если n == 2, возвращаем количество начальных состояний
	if n == 2 {
		result := 0
		for i := 0; i < 100; i++ {
			result = (result + start[i]) % mod
		}
		return result
	}

	// Возводим матрицу в степень (n-2), так как у нас уже есть первые 2 цифры
	power := n - 2
	M_power := matrixPower(M, power)

	// Умножаем начальный вектор на матрицу
	result := 0
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			result = (result + start[i]*M_power[i][j]%mod) % mod
		}
	}

	return result
}

// matrixPower возводит матрицу в степень используя быстрое возведение
func matrixPower(M [][]int, power int64) [][]int {
	n := len(M)

	// Инициализируем единичную матрицу
	result := make([][]int, n)
	for i := 0; i < n; i++ {
		result[i] = make([]int, n)
		result[i][i] = 1
	}

	base := make([][]int, n)
	for i := 0; i < n; i++ {
		base[i] = make([]int, n)
		copy(base[i], M[i])
	}

	// Быстрое возведение в степень
	for power > 0 {
		if power&1 == 1 {
			result = matrixMultiply(result, base)
		}
		base = matrixMultiply(base, base)
		power >>= 1
	}

	return result
}

// matrixMultiply умножает две матрицы по модулю
func matrixMultiply(A, B [][]int) [][]int {
	n := len(A)
	m := len(B[0])
	k := len(B)

	result := make([][]int, n)
	for i := 0; i < n; i++ {
		result[i] = make([]int, m)
		for j := 0; j < m; j++ {
			sum := 0
			for t := 0; t < k; t++ {
				sum = (sum + A[i][t]*B[t][j]%mod) % mod
			}
			result[i][j] = sum
		}
	}

	return result
}
