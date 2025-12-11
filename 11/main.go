package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const mod int64 = 1000000007

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	writer := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer writer.Flush()

	// Читаем n
	line, _ := reader.ReadString('\n')
	n, _ := strconv.Atoi(strings.TrimSpace(line))

	result := solve(n)
	writer.WriteString(fmt.Sprintf("%d\n", result))
}

// solve вычисляет E(n) mod M для гиперкуба размера n
// Решает систему линейных уравнений для ожидаемого времени в гиперкубе
// используя модульную арифметику
func solve(n int) int64 {
	// Специальный случай
	if n == 1 {
		return 1
	}

	// Система уравнений:
	// E[0] = 0
	// E[k] = 1 + (k/n)*E[k-1] + ((n-k)/n)*E[k+1]  для 1 <= k <= n-1
	// E[n] = 1 + E[n-1]
	//
	// Умножим на n:
	// n*E[k] = n + k*E[k-1] + (n-k)*E[k+1]
	// k*E[k-1] - n*E[k] + (n-k)*E[k+1] = -n
	//
	// Трехдиагональная система:
	// a[k]*E[k-1] + b[k]*E[k] + c[k]*E[k+1] = d[k]
	// где a[k] = k, b[k] = -n, c[k] = n-k, d[k] = -n

	// Метод прогонки (Thomas algorithm) в модульной арифметике
	// Прямой ход: преобразуем к виду E[k] = alpha[k]*E[k+1] + beta[k]

	nMod := int64(n) % mod

	// alpha[k] и beta[k] для прогонки
	// E[k] = alpha[k]*E[k+1] + beta[k]
	alpha := make([]int64, n+1)
	beta := make([]int64, n+1)

	// Начальное условие: E[0] = 0
	// Из уравнения для k=1: 1*E[0] - n*E[1] + (n-1)*E[2] = -n
	// -n*E[1] + (n-1)*E[2] = -n
	// E[1] = ((n-1)*E[2] + n) / n = (n-1)/n * E[2] + 1

	// Для k=0: E[0] = 0, нет E[-1], так что alpha[0] = 0, beta[0] = 0
	alpha[0] = 0
	beta[0] = 0

	// Прямой ход: для k = 1..n-1
	// k*E[k-1] - n*E[k] + (n-k)*E[k+1] = -n
	// Подставляем E[k-1] = alpha[k-1]*E[k] + beta[k-1]:
	// k*(alpha[k-1]*E[k] + beta[k-1]) - n*E[k] + (n-k)*E[k+1] = -n
	// (k*alpha[k-1] - n)*E[k] + (n-k)*E[k+1] = -n - k*beta[k-1]
	// E[k] = (n-k)/(n - k*alpha[k-1]) * E[k+1] + (-n - k*beta[k-1])/(n - k*alpha[k-1])

	for k := 1; k < n; k++ {
		kMod := int64(k) % mod
		nMinusK := int64(n-k) % mod

		// denominator = n - k*alpha[k-1]
		denom := (nMod - kMod*alpha[k-1]%mod + mod) % mod
		denomInv := modInverse(denom, mod)

		// alpha[k] = (n-k) / denom
		alpha[k] = nMinusK * denomInv % mod

		// beta[k] = (n + k*beta[k-1]) / denom
		numerator := (nMod + kMod*beta[k-1]%mod) % mod
		beta[k] = numerator * denomInv % mod
	}

	// Граничное условие: E[n] = 1 + E[n-1]
	// Подставляем E[n-1] = alpha[n-1]*E[n] + beta[n-1]:
	// E[n] = 1 + alpha[n-1]*E[n] + beta[n-1]
	// E[n]*(1 - alpha[n-1]) = 1 + beta[n-1]
	// E[n] = (1 + beta[n-1]) / (1 - alpha[n-1])

	oneMinusAlpha := (1 - alpha[n-1] + mod) % mod
	oneMinusAlphaInv := modInverse(oneMinusAlpha, mod)
	En := (1 + beta[n-1]) % mod * oneMinusAlphaInv % mod

	return En
}

// modInverse вычисляет обратное число a^(-1) mod m используя малую теорему Ферма
// a^(-1) = a^(m-2) mod m (для простого m)
func modInverse(a, m int64) int64 {
	return modPow(a, m-2, m)
}

// modPow вычисляет base^exp mod m с использованием быстрого возведения в степень
func modPow(base, exp, m int64) int64 {
	result := int64(1)
	base = base % m
	for exp > 0 {
		if exp%2 == 1 {
			result = result * base % m
		}
		exp = exp >> 1
		base = base * base % m
	}
	return result
}
