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

	// Предвычисляем факториалы и обратные факториалы
	maxN := 400001
	fact := precomputeFactorials(maxN)
	invFact := precomputeInvFactorials(fact, maxN)

	// Читаем количество тестов
	line, _ := reader.ReadString('\n')
	T, _ := strconv.Atoi(strings.TrimSpace(line))

	// Обрабатываем тесты
	for i := 0; i < T; i++ {
		line, _ = reader.ReadString('\n')
		parts := strings.Fields(strings.TrimSpace(line))
		n, _ := strconv.Atoi(parts[0])
		s, _ := strconv.Atoi(parts[1])

		result := solve(n, s, fact, invFact)
		writer.WriteString(fmt.Sprintf("%d\n", result))
	}
}

// solve считает количество замечательных массивов
// Формула: answer(n, s) = (n+1)! × C(s, n)
func solve(n, s int, fact, invFact []int) int {
	if n > s {
		return 0
	}

	// C(s, n) = s! / (n! * (s-n)!)
	c := comb(s, n, fact, invFact)

	// (n+1)! × C(s, n)
	return fact[n+1] * c % mod
}

// comb вычисляет C(n, k) по модулю mod
func comb(n, k int, fact, invFact []int) int {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

// precomputeFactorials предвычисляет факториалы до n
func precomputeFactorials(n int) []int {
	fact := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	return fact
}

// precomputeInvFactorials предвычисляет обратные факториалы
func precomputeInvFactorials(fact []int, n int) []int {
	invFact := make([]int, n+1)
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * i % mod
	}
	return invFact
}

// modPow вычисляет a^b mod mod
func modPow(a, b int) int {
	result := 1
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			result = result * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return result
}
